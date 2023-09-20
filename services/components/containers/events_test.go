package containers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/golang/mock/gomock"
)

var (
	testEventBotID         = "0x12345"
	testEventContainerName = "test-container-name"
	testImageName          = "test-image-name"
	testEventTime          = time.Now()
)

func TestListenToDockerEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)
	msgClient := mock_clients.NewMockMessageClient(ctrl)

	eventsSince := time.Now()
	ctx := context.Background()
	eventCh := make(chan events.Message)
	errCh := make(chan error)

	doneCh := make(chan struct{})
	dockerClient.EXPECT().Events(ctx, eventsSince).Return(eventCh, errCh).Do(func(v1, v2 interface{}) {
		close(doneCh)
	})
	resettedSince := gomock.Not(eventsSince)
	dockerClient.EXPECT().Events(ctx, resettedSince).Return(eventCh, errCh)
	go ListenToDockerEvents(ctx, dockerClient, msgClient, eventsSince)
	errCh <- errors.New("test error") // causes the second dockerClient.Events() call
	<-doneCh

	doneCh = make(chan struct{})
	msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any()).Do(func(v1, v2 interface{}) {
		close(doneCh)
	})
	eventCh <- events.Message{
		Type:   "container",
		Action: "create",
		Actor: events.Actor{
			Attributes: map[string]string{
				docker.LabelFortaBotID: testEventBotID,
			},
		},
	}
	<-doneCh
}

func TestHandleEvent(t *testing.T) {
	testCases := []struct {
		incomingEvent  events.Message
		producedMetric *protocol.AgentMetric
	}{
		{
			incomingEvent: events.Message{
				Type:   "image",
				Action: "foo", // causes skip
				Actor: events.Actor{
					Attributes: map[string]string{
						docker.LabelFortaBotID: testEventBotID,
						"name":                 testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: nil,
		},
		{
			incomingEvent: events.Message{
				Type:   "foo", // causes skip
				Action: "bar",
				Actor: events.Actor{
					Attributes: map[string]string{
						docker.LabelFortaBotID: testEventBotID,
						"name":                 testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: nil,
		},
		{
			incomingEvent: events.Message{
				Type:   "image",
				Action: "pull",
				Actor: events.Actor{
					Attributes: map[string]string{
						"name": testImageName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   "system",
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.image.pull",
				Value:     1,
				Details:   testImageName,
			},
		},
		{
			incomingEvent: events.Message{
				Type:   "container",
				Action: "create",
				Actor: events.Actor{
					Attributes: map[string]string{
						docker.LabelFortaBotID: testEventBotID,
						"name":                 testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   testEventBotID,
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.container.create",
				Value:     1,
				Details:   testEventContainerName,
			},
		},
		{
			incomingEvent: events.Message{
				Type:   "container",
				Action: "destroy",
				Actor: events.Actor{
					Attributes: map[string]string{
						docker.LabelFortaBotID: testEventBotID,
						"name":                 testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   testEventBotID,
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.container.destroy",
				Value:     1,
				Details:   testEventContainerName,
			},
		},
		{
			incomingEvent: events.Message{
				Type:   "network",
				Action: "create",
				Actor: events.Actor{
					Attributes: map[string]string{
						"name": testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   "system",
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.network.create",
				Value:     1,
				Details:   testEventContainerName,
			},
		},
		{
			incomingEvent: events.Message{
				Type:   "network",
				Action: "destroy",
				Actor: events.Actor{
					Attributes: map[string]string{
						"name": testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   "system",
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.network.destroy",
				Value:     1,
				Details:   testEventContainerName,
			},
		},
		{
			incomingEvent: events.Message{
				Type:   "network",
				Action: "connect",
				Actor: events.Actor{
					Attributes: map[string]string{
						"name": testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   "system", // disregards bot id in the label
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.network.connect",
				Value:     1,
				Details:   testEventContainerName,
			},
		},
		{
			incomingEvent: events.Message{
				Type:   "network",
				Action: "disconnect",
				Actor: events.Actor{
					Attributes: map[string]string{
						"name": testEventContainerName,
					},
				},
				Time:     testEventTime.Unix(),
				TimeNano: testEventTime.UnixNano(),
			},
			producedMetric: &protocol.AgentMetric{
				AgentId:   "system", // disregards bot id in the label
				Timestamp: testEventTime.Format(time.RFC3339),
				Name:      "docker.network.disconnect",
				Value:     1,
				Details:   testEventContainerName,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.incomingEvent.Type, testCase.incomingEvent.Action), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dockerClient := mock_clients.NewMockDockerClient(ctrl)
			msgClient := mock_clients.NewMockMessageClient(ctrl)
			handler := &eventHandler{
				dockerClient: dockerClient,
				msgClient:    msgClient,
			}
			ctx := context.Background()

			if testCase.producedMetric != nil {
				msgClient.EXPECT().PublishProto(
					messaging.SubjectMetricAgent, newMetrics(testCase.producedMetric),
				)
			}
			handler.HandleEvent(ctx, &testCase.incomingEvent)
		})
	}
}

func newMetrics(ms ...*protocol.AgentMetric) *protocol.AgentMetricList {
	return &protocol.AgentMetricList{
		Metrics: ms,
	}
}
