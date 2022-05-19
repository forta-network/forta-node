package network

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/forta-network/forta-core-go/utils"
	"github.com/sirupsen/logrus"
)

type botAdminUnixSockServer struct {
	server   *http.Server
	listener net.Listener
	botAdmin *botAdmin
}

// NewBotAdminServer starts a new server.
func NewBotAdminServer(containerName string) (*botAdminUnixSockServer, error) {
	listener, err := net.Listen("unix", sockPath(containerName))
	if err != nil {
		return nil, err
	}

	admin := &botAdmin{}

	server := &http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var ruleCmds [][]string
			if err := json.NewDecoder(r.Body).Decode(&ruleCmds); err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				logrus.WithError(err).Error("failed to decode request")
				return
			}
			if err := admin.IPTables(ruleCmds); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				logrus.WithError(err).Error("failed to execute iptables rules")
				return
			}
		}),
	}

	return &botAdminUnixSockServer{
		server:   server,
		listener: listener,
		botAdmin: admin,
	}, nil
}

func (ba *botAdminUnixSockServer) Start() error {
	utils.GoListenAndServe(ba.server)
	return nil
}

func (ba *botAdminUnixSockServer) Stop() error {
	return nil
}

func (ba *botAdminUnixSockServer) Name() string {
	return "bot-admin"
}
