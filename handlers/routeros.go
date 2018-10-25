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
