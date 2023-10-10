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
	functions.CloudEvent("UpdateTodoItem", updateTodoItem)

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
	ChannelId   string `json:"channelId"`
	UserId      string `json:"userId"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

type Response struct {
	Result bool        `json:"result"`
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
}

func publishResponse(funcName string, req Request, data interface{}) error {
	jsonBytes, err := json.Marshal(Response{Result: true, Type: funcName, Data: data})
	if err != nil {
		fmt.Println(err)
		return err
	}
	println(string(jsonBytes))
	if err := messageQueue.PublishNotice(req.ChannelId, req.UserId, string(jsonBytes)); err != nil {
		return err
	}
	return nil
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
	if data.Name == "" || data.ChannelId == "" || data.UserId == "" {
		return errors.New("bad request")
	}
	list, err := db.GetTodoList(data.Name)
	if err != nil {
		return err
	}

	// publish to pubsub
	if err := publishResponse("getTodoList", data, *list); err != nil {
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
	if data.Name == "" || data.ChannelId == "" || data.UserId == "" {
		return errors.New("bad request")
	}

	docId,err := db.AddTodoItem(data.Name, data.Description, *data.Completed);
	if  err != nil {
		return err
	}

	// publish to pubsub
	if err := publishResponse("addTodoItem", data, docId); err != nil {
		return err
	}

	return nil
}

func removeTodoItem(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	dataBuffer := string(msg.Message.Data) // Automatically decoded from base64.
	data := Request{}
	if err := json.Unmarshal([]byte(dataBuffer), &data); err != nil {
		return err
	}
	if data.ChannelId == "" || data.UserId == "" || data.ID == "" {
		return errors.New("bad request")
	}

	id,err := db.RemoveTodoItem(data.ID)
	if ; err != nil {
		return err
	}

	// publish to pubsub
	if err := publishResponse("removeTodoItem", data, id); err != nil {
		return err
	}

	return nil
}

func updateTodoItem(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	dataBuffer := string(msg.Message.Data) // Automatically decoded from base64.
	data := Request{}
	if err := json.Unmarshal([]byte(dataBuffer), &data); err != nil {
		return err
	}
	if data.ChannelId == "" || data.UserId == "" || data.ID == "" || data.Description == "" || data.Completed == nil {
		return errors.New("bad request")
	}
	
	id, err := db.SetTodoItem(data.ID, data.Description, *data.Completed)
	if ; err != nil {
		return err
	}

	// publish to pubsub
	if err := publishResponse("updateTodoItem", data, id); err != nil {
		return err
	}

	return nil
}
