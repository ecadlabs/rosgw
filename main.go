package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecadlabs/rosgw/config"
	"github.com/ecadlabs/rosgw/errors"
	"github.com/ecadlabs/rosgw/handlers"
	"github.com/ecadlabs/rosgw/middleware"
	"github.com/ecadlabs/rosgw/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfgFile := flag.String("c", "", "Config file.")
	flag.Parse()

	if *cfgFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	conf, err := config.Load(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewRouterosHandler(conf, conf.MaxConn, log.StandardLogger())

	m := mux.NewRouter()

	m.Use((&middleware.Logging{}).Handler)
	m.Use((&middleware.Recover{}).Handler)

	m.Methods("GET").Path("/devices/{id}/interfaces").HandlerFunc(handler.GetInterfaces)

	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.JSONErrorResponse(w, errors.ErrResourceNotFound)
	})

	log.Printf("HTTP service listening on %s", conf.Address)

	httpServer := &http.Server{
		Addr:    conf.Address,
		Handler: m,
	}

	errChan := make(chan error, 10)
	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	defer httpServer.Shutdown(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}

		case s := <-signalChan:
			log.Printf("Captured %v. Exiting...\n", s)
			return
		}
	}
}
