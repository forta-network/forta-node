package containers

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/sirupsen/logrus"
)

// ListenToDockerEvents creates new.
func ListenToDockerEvents(
	ctx context.Context, dockerClient clients.DockerClient, msgClient clients.MessageClient,
	startFrom time.Time,
) {
	handler := &eventHandler{
		dockerClient: dockerClient,
		msgClient:    msgClient,
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Info("stopping docker events listener")
			return
		default:
		}

		events, errs := dockerClient.Events(ctx, startFrom)

		var restartListening bool
		for {
			select {
			case event := <-events:
				handler.HandleEvent(ctx, &event)

			case err := <-errs:
				logrus.WithError(err).Error("error while listening to docker events")
				// set the start time and restart listening
				startFrom = time.Now().Add(-1 * time.Second)
				restartListening = true

			case <-ctx.Done():
				logrus.Info("stopping docker events listener")
				return
			}
			if restartListening {
				break
			}
		}
	}
}

type eventHandler struct {
	dockerClient clients.DockerClient
	msgClient    clients.MessageClient
}

func (es *eventHandler) HandleEvent(ctx context.Context, event *events.Message) {
	var metric *protocol.AgentMetric
	ts := time.Unix(0, event.TimeNano)
	switch event.Type {
	case "image":
		if event.Action != "pull" {
			return
		}
		imageRef := getEventAttribute(event, "name")
		metric = metrics.CreateEventMetric(ts, "system", metricNameFrom(event), imageRef)

	case "container", "network":
		if !isOneOf(event.Action, "create", "destroy", "connect", "disconnect") {
			return
		}
		botID, ok := getBotID(event)
		if !ok {
			botID = "system"
		}
		containerName := getEventAttribute(event, "name")
		metric = metrics.CreateEventMetric(ts, botID, metricNameFrom(event), containerName)

	default:
		return
	}

	metrics.SendAgentMetrics(es.msgClient, []*protocol.AgentMetric{metric})
	return
}

func isOneOf(input string, values ...string) bool {
	for _, value := range values {
		if input == value {
			return true
		}
	}
	return false
}

func metricNameFrom(event *events.Message) string {
	return strings.Join([]string{"docker", event.Type, event.Action}, ".")
}

func getBotID(event *events.Message) (string, bool) {
	val := getEventAttribute(event, docker.LabelFortaBotID)
	return val, len(val) > 0
}

func getEventAttribute(event *events.Message, attr string) string {
	if event.Actor.Attributes == nil {
		return ""
	}
	return event.Actor.Attributes[attr]
}
