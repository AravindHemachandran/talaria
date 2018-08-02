package main

import (
	"fmt"

	"github.com/Comcast/webpa-common/logging/logginghttp"
	"github.com/Comcast/webpa-common/xhttp"
	"github.com/Comcast/webpa-common/xhttp/gate"
	"github.com/Comcast/webpa-common/xhttp/xcontext"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	ControlKey = "control"
)

func StartControlServer(logger log.Logger, v *viper.Viper) error {
	if !v.IsSet(ControlKey) {
		return nil
	}

	var options xhttp.ServerOptions
	if err := v.UnmarshalKey(ControlKey, &options); err != nil {
		return err
	}

	options.Logger = logger

	var (
		g          = gate.New(gate.Open)
		r          = mux.NewRouter()
		apiHandler = r.PathPrefix(fmt.Sprintf("%s/%s", baseURI, version)).Subrouter()
	)

	apiHandler.Handle("/device/gate", &gate.Lever{Gate: g, Parameter: "open"}).
		Methods("POST", "PUT", "PATCH")

	apiHandler.Handle("/device/gate", &gate.Status{Gate: g}).
		Methods("GET")

	server := xhttp.NewServer(options)
	server.Handler = xcontext.Populate(0, logginghttp.SetLogger(logger))(r)
	return xhttp.NewStarter(options.StartOptions(), server)()
}
