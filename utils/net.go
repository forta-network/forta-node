package utils

import (
	"net"
	"net/http"
	"strings"

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

// ConvertToDockerHostURL converts localhost HTTP URLs to docker host URLs that are
// dialable from within a scanner.
func ConvertToDockerHostURL(rawurl string) string {
	rawurl = strings.ReplaceAll(rawurl, "http://127.0.0.1", "http://host.docker.internal")
	rawurl = strings.ReplaceAll(rawurl, "http://localhost", "http://host.docker.internal")
	return rawurl
}
