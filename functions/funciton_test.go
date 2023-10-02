package function

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
)

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
			Completed:   false,
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
