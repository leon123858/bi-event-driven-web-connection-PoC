package database

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DB struct {
	Client    *firestore.Client
	Ctx       context.Context
	ProjectID string
}

type TodoItem struct {
	ID          string `firestore:"id,omitempty"`
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	Completed   bool   `firestore:"completed"`
	Timestamp   int64  `firestore:"timestamp"`
}

func NewDB(projectID string, ctx context.Context) (*DB, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return nil, err
	}
	// Close client when done with
	// defer client.Close()
	return &DB{
		Client:    client,
		Ctx:       ctx,
		ProjectID: projectID,
	}, nil
}

func (db *DB) AddTodoItem(name, description string, completed bool) (string,error) {
	now := time.Now().UnixMicro()
	ref, _, err := db.Client.Collection("todoList").Add(db.Ctx, &TodoItem{
		Name:        name,
		Description: description,
		Completed:   completed,
		Timestamp:   now,
	})
	if err != nil {
		log.Printf("Failed adding item: %v", err)
		return "",err
	}
	return  ref.ID,nil
}

func (db *DB) RemoveTodoItem(id string) (string,error) {
	_, err := db.Client.Collection("todoList").Doc(id).Delete(db.Ctx)
	if err != nil {
		log.Printf("Failed to delete item with ID %s: %v", id, err)
		return "",err
	}
	return id,nil
}

func (db *DB) SetTodoItem(id string, description string, completed bool) (string,error) {
	_, err := db.Client.Collection("todoList").Doc(id).Update(db.Ctx, []firestore.Update{
		{
			Path:  "description",
			Value: description,
		}, {
			Path:  "completed",
			Value: completed,
		},
	})
	if err != nil {
		log.Printf("Failed to set item with ID %s: %v", id, err)
		return "", err
	}
	return id, nil
}

func (db *DB) GetTodoList(name string) (*[]TodoItem, error) {
	var todoList []TodoItem
	iter := db.Client.Collection("todoList").Where("name", "==", name).OrderBy("timestamp", firestore.Desc).Documents(db.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
			return nil, err
		}
		// convert firestore item  to TodoItem
		newItem := TodoItem{
			ID:          doc.Ref.ID,
			Name:        doc.Data()["name"].(string),
			Description: doc.Data()["description"].(string),
			Completed:   doc.Data()["completed"].(bool),
			Timestamp:   doc.Data()["timestamp"].(int64),
		}
		todoList = append(todoList, newItem)
	}
	return &todoList, nil
}
