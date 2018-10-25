package handlers

import (
	"context"
	"net/http"

	"github.com/ecadlabs/rosgw/command"
	"github.com/ecadlabs/rosgw/config"
	"github.com/ecadlabs/rosgw/conn"
	"github.com/ecadlabs/rosgw/response"
	"github.com/ecadlabs/rosgw/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type DeviceDB interface {
	GetDevice(context.Context, string) (*config.Device, error)
}

type Routeros struct {
	cmd    *command.RemoteCommand
	logger *logrus.Logger
	db     DeviceDB
}

func NewRouterosHandler(db DeviceDB, maxconn int, logger *logrus.Logger) *Routeros {
	return &Routeros{
		logger: logger,
		db:     db,
		cmd: &command.RemoteCommand{
			Logger: logger,
			Pool:   conn.NewConnPool(maxconn),
		},
	}
}

func (r *Routeros) GetInterfaces(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	dev, err := r.db.GetDevice(req.Context(), id)
	if err != nil {
		utils.JSONErrorResponse(w, err)
		return
	}

	ctx := req.Context()
	if to := dev.GetTimeout(); to != 0 {
		ctx, _ = context.WithTimeout(ctx, to)
	}

	res, err := r.cmd.Run(ctx, dev.Address(), dev.Config(), "interface print detail")
	if err != nil {
		r.logger.Error(err)
		utils.JSONErrorResponse(w, err)
		return
	}
	defer res.Close()

	var ifaces response.Interfaces
	if err := ifaces.ParseResponse(res); err != nil {
		r.logger.Error(err)
		utils.JSONErrorResponse(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, ifaces)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (r *Routeros) MonitorTraffic(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	dev, err := r.db.GetDevice(req.Context(), id)
	if err != nil {
		utils.JSONErrorResponse(w, err)
		return
	}

	iface := mux.Vars(req)["iface"]

	res, err := r.cmd.Run(req.Context(), dev.Address(), dev.Config(), "interface monitor-traffic "+iface)
	if err != nil {
		r.logger.Error(err)
		utils.JSONErrorResponse(w, err)
		return
	}

	defer func() {
		res.Close()
		r.logger.Info("Session closed")
	}()

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		r.logger.Error(err)
		return
	}

	defer conn.Close()

	ch := make(chan *response.TrafficMonitorResponse, 100)

	mon := response.TrafficMonitor{
		C:        ch,
		NonBlock: true,
	}

	// Run parser
	go func() {
		if err := mon.ParseResponse(res); err != nil {
			r.logger.Error(err)
		}

		close(ch)
	}()

	// WS handler loop
	for pkt := range ch {
		if err := conn.WriteJSON(pkt); err != nil {
			r.logger.Error(err)
			return
		}
	}
}
