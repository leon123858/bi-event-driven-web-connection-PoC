package function

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/leon123858/bi-event-driven-web-connection-PoC/functions/database"
	"github.com/leon123858/bi-event-driven-web-connection-PoC/functions/pubsub"
)

var db *database.DB
var messageQueue *pubsub.PubSubInfo

func init() {
	projectId := "tw-rd-ca-leon-lin"

	functions.CloudEvent("GetTodoList", getTodoList)
	functions.CloudEvent("AddTodoItem", addTodoItem)
	functions.CloudEvent("RemoveTodoItem", removeTodoItem)

	// create golang wait group
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_db, err := database.NewDB(projectId, context.Background())
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		db = _db
	}()
	go func() {
		defer wg.Done()
		_mq, err := pubsub.NewPubSub(projectId)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		messageQueue = _mq
	}()
	wg.Wait()
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Request struct {
	Name        string `json:"name"`
	ChannelId   int64  `json:"channelId"`
	UserId      string `json:"userId"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func getTodoList(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	dataBuffer := string(msg.Message.Data) // Automatically decoded from base64.
	data := Request{}
	if err := json.Unmarshal([]byte(dataBuffer), &data); err != nil {
		return err
	}
	if data.Name == "" || data.ChannelId == 0 || data.UserId == "" {
		return errors.New("bad request")
	}

	list, err := db.GetTodoList(data.Name)
	if err != nil {
		return err
	}

	// Convert the array of structs to a JSON string.
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// publish to pubsub
	if err := messageQueue.PublishNotice(data.ChannelId, data.UserId, string(jsonBytes)); err != nil {
		return err
	}

	return nil
}

func addTodoItem(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	dataBuffer := string(msg.Message.Data) // Automatically decoded from base64.
	data := Request{}
	if err := json.Unmarshal([]byte(dataBuffer), &data); err != nil {
		return err
	}
	if data.Name == "" || data.ChannelId == 0 || data.UserId == "" {
		return errors.New("bad request")
	}

	if err := db.AddTodoItem(data.Name, data.Description, data.Completed); err != nil {
		return err
	}

	return nil
}

func removeTodoItem(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	name := string(msg.Message.Data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}

func updateTodoItem(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	name := string(msg.Message.Data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}
