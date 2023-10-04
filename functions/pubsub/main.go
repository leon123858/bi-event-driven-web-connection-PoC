package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type PubSubInfo struct {
	ProjectID string
	Client    *pubsub.Client
	Ctx       context.Context
}

type Notice struct {
	UserId string `json:"userId"`
	Msg    string `json:"msg"`
}

func NewPubSub(projectId string) (*PubSubInfo, error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return &PubSubInfo{
		ProjectID: projectId,
		Client:    client,
		Ctx:       ctx,
	}, nil
}

// create topic when not exists
func (info *PubSubInfo) CreateTopic(topicId string) error {
	topic := info.Client.Topic(topicId)
	exists, err := topic.Exists(info.Ctx)
	if err != nil {
		return err
	}
	if !exists {
		_, err := info.Client.CreateTopic(info.Ctx, topicId)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove topic when exists
func (info *PubSubInfo) RemoveTopic(topicId string) error {
	topic := info.Client.Topic(topicId)
	exists, err := topic.Exists(info.Ctx)
	if err != nil {
		return err
	}
	if exists {
		if err := topic.Delete(info.Ctx); err != nil {
			return err
		}
	}
	return nil
}

func (info *PubSubInfo) PublishNotice(channelId int64, userId, message string) error {
	// Get a topic reference.
	topic := info.Client.Topic("notice-" + fmt.Sprint(channelId))

	// create bytes[] from Notice struct
	notice := Notice{
		UserId: userId,
		Msg:    message,
	}
	data, err := json.Marshal(notice)
	if err != nil {
		return err
	}

	// Publish the message.
	ctx := context.Background()
	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
