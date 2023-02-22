package jwt_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/forta-network/forta-node/config"
)

type CreateJWTMessage struct {
	Claims map[string]interface{} `json:"claims"`
}
type CreateJWTResponse struct {
	Token string `json:"token"`
}

const envPrefix = config.EnvFortaBotID + "="

const (
	errBadCreateMessage  = "bad create jwt message body"
	errFailedToCreateJWT = "can't find bot id from request source"
)

// agentIDReverseLookup reverse lookup from ip to agent id.
func (j *JWTProvider) agentIDReverseLookup(ctx context.Context, ipAddr string) (string, error) {
	container, err := j.findContainerByIP(ctx, ipAddr)
	if err != nil {
		return "", err
	}

	botID, err := j.extractBotIDFromContainer(ctx, container)
	if err != nil {
		return "", err
	}

	return botID, nil
}

func (j *JWTProvider) extractBotIDFromContainer(ctx context.Context, container types.Container) (string, error) {
	// container struct doesn't have the "env" information, inspection required.
	c, err := j.dockerClient.InspectContainer(ctx, container.ID)
	if err != nil {
		return "", err
	}

	// find the env variable with bot id
	for _, s := range c.Config.Env {
		if env := strings.SplitAfter(s, envPrefix); len(env) == 2 {
			return env[1], nil
		}
	}

	return "", fmt.Errorf("can't extract bot id from container")
}

func (j *JWTProvider) findContainerByIP(ctx context.Context, ipAddr string) (types.Container, error) {
	containers, err := j.dockerClient.GetContainers(ctx)
	if err != nil {
		return types.Container{}, err
	}

	// find the container that has the same ip
	for _, container := range containers {
		for _, network := range container.NetworkSettings.Networks {
			if network.IPAddress == ipAddr {
				return container, nil
			}
		}
	}
	return types.Container{}, fmt.Errorf("can't find container %s", ipAddr)
}

// createJWTHandler returns a scanner jwt token with claims [hash] = hash(uri,payload) and [bot] = "bot id"
func (j *JWTProvider) createJWTHandler(w http.ResponseWriter, req *http.Request) {
	var msg CreateJWTMessage
	if req.Body != http.NoBody {
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, errBadCreateMessage)
			return
		}
	}

	jwt, err := j.doCreateJWT(req.Context(), req.RemoteAddr, msg.Claims)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "can not create jwt: %v", err)
		return
	}

	resp, err := json.Marshal(CreateJWTResponse{Token: jwt})

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", resp)
}

func (j *JWTProvider) doCreateJWT(ctx context.Context, remoteAddr string, claims map[string]interface{}) (string, error) {
	ipAddr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", fmt.Errorf("can't extract ip from request %s", remoteAddr)
	}

	agentID, err := j.agentIDReverseLookup(ctx, ipAddr)
	if err != nil {
		return "", fmt.Errorf("can't find bot id from request source %s, err: %v", ipAddr, err)
	}

	jwt, err := CreateBotJWT(j.cfg.Key, agentID, claims)
	if err != nil {
		return "", fmt.Errorf("%s: %v", errFailedToCreateJWT, err)
	}

	return jwt, nil
}