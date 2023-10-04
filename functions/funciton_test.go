package function

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
)

func setup() {
	println("setup start")
	messageQueue.CreateTopic("notice-1")
}

func teardown() {
	list, _ := db.GetTodoList("test")
	for _, item := range *(list) {
		db.RemoveTodoItem(item.ID)
	}
	messageQueue.RemoveTopic("notice-1")
	println("teardown complete")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_getTodoList(t *testing.T) {
	type args struct {
		ctx  context.Context
		data Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"basic", args{context.Background(), Request{
			Name:      "test",
			ChannelId: int64(1),
			UserId:    "test",
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert the struct to bytes.
			jsonBytes, err := json.Marshal(tt.args.data)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg := MessagePublishedData{
				Message: PubSubMessage{
					Data: []byte(jsonBytes),
				},
			}

			e := event.New()
			e.SetDataContentType("application/json")
			e.SetData(e.DataContentType(), msg)
			if err := getTodoList(tt.args.ctx, e); (err != nil) != tt.wantErr {
				t.Errorf("getTodoList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addTodoItem(t *testing.T) {
	ptrBool := new(bool)
	*ptrBool = false

	type args struct {
		ctx  context.Context
		data Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"basic", args{context.Background(), Request{
			Name:        "test",
			ChannelId:   int64(1),
			UserId:      "test",
			Description: "test",
			Completed:   ptrBool,
		}}, false}, {"basic", args{context.Background(), Request{
			Name:        "test",
			ChannelId:   int64(1),
			UserId:      "test2",
			Description: "test2",
			Completed:   ptrBool,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert the struct to bytes.
			jsonBytes, err := json.Marshal(tt.args.data)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg := MessagePublishedData{
				Message: PubSubMessage{
					Data: []byte(jsonBytes),
				},
			}

			e := event.New()
			e.SetDataContentType("application/json")
			e.SetData(e.DataContentType(), msg)
			if err := addTodoItem(tt.args.ctx, e); (err != nil) != tt.wantErr {
				t.Errorf("addTodoItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_removeTodoItem(t *testing.T) {
	type args struct {
		ctx  context.Context
		data Request
	}

	list, _ := db.GetTodoList("test")
	beginLen := len(*list)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"basic", args{context.Background(), Request{
			Name:        "test",
			ChannelId:   int64(1),
			UserId:      "test",
			Description: "test",
		}}, true}, {"basic", args{context.Background(), Request{
			Name:      "test",
			ChannelId: int64(1),
			UserId:    "test",
			ID:        (*list)[0].ID,
		}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert the struct to bytes.
			jsonBytes, err := json.Marshal(tt.args.data)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg := MessagePublishedData{
				Message: PubSubMessage{
					Data: []byte(jsonBytes),
				},
			}

			e := event.New()
			e.SetDataContentType("application/json")
			e.SetData(e.DataContentType(), msg)
			if err := removeTodoItem(tt.args.ctx, e); (err != nil) != tt.wantErr {
				t.Errorf("removeTodoItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	list, _ = db.GetTodoList("test")
	finalLen := len(*list)

	if beginLen != finalLen+1 {
		t.Errorf("removeTodoItem() error = %v, wantErr %v", beginLen, finalLen+1)
	}
}

func Test_updateTodoItem(t *testing.T) {
	type args struct {
		ctx  context.Context
		data Request
	}

	// add item
	db.AddTodoItem("test", "test", false)
	list, _ := db.GetTodoList("test")

	ptrBool := new(bool)
	*ptrBool = true

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"basic", args{context.Background(), Request{
			ChannelId:   int64(1),
			UserId:      "test",
			ID:          (*list)[0].ID,
			Description: "new description",
		}}, true}, {"basic", args{context.Background(), Request{
			ChannelId:   int64(1),
			UserId:      "test",
			ID:          (*list)[0].ID,
			Description: "new description",
			Completed:   ptrBool,
		}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert the struct to bytes.
			jsonBytes, err := json.Marshal(tt.args.data)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg := MessagePublishedData{
				Message: PubSubMessage{
					Data: []byte(jsonBytes),
				},
			}

			e := event.New()
			e.SetDataContentType("application/json")
			e.SetData(e.DataContentType(), msg)
			if err := updateTodoItem(tt.args.ctx, e); (err != nil) != tt.wantErr {
				t.Errorf("updateTodoItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	list, _ = db.GetTodoList("test")

	if (*list)[0].Description != "new description" {
		t.Errorf("updateTodoItem() error = %v, wantErr %v", (*list)[0].Description, "new description")
	}
	if (*list)[0].Completed != true {
		t.Errorf("updateTodoItem() error = %v, wantErr %v", (*list)[0].Completed, true)
	}
}
