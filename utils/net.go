package utils

import (
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GoListenAndServe starts the server and does meaningful error handling on exit.
func GoListenAndServe(server *http.Server) {
	go func() {
		switch err := server.ListenAndServe(); err {
		case nil, http.ErrServerClosed:
			// do nothing
		default:
			log.WithError(err).Panic("server error")
		}
	}()
}

// GoGrpcServe starts the server and does meaningful error handling on exit.
func GoGrpcServe(server *grpc.Server, lis net.Listener) {
	go func() {
		switch err := server.Serve(lis); err {
		case nil, grpc.ErrServerStopped:
			// do nothing
		default:
			log.WithError(err).Panic("server error")
		}
	}()
}
