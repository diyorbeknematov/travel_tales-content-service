package postgres

import (
	pb "content-service/generated/communication"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewCommunicationRepo(db)
	req := &pb.SendMessageRequest{
		SendeId:    "e1b9af75-931d-4d3b-acd7-a00e2571fa92",
		RecipientId: "9b0cf2c8-308c-4896-a737-511bff1bb991",
		Content:     "Hello, world!",
	}

	resp, err := repo.SendMessage(req)
	assert.NoError(t, err)
	assert.Equal(t, req.SendeId, resp.SenderId)
	assert.Equal(t, req.RecipientId, resp.RecipientId)
	assert.Equal(t, req.Content, resp.Content)
}

func TestGetMessages(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewCommunicationRepo(db)
	req := &pb.ListMessageRequest{
		Page:  1,
		Limit: 10,
	}

	resp, err := repo.GetMessages(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Message)
}

func TestCreateTravelTips(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewCommunicationRepo(db)
	req := &pb.AddTravelTipsRequest{
		Title:    "Tip",
		Content:  "This is a test tip.",
		Category: "Category",
		AuthorId: "a779fc79-0cd2-47fe-a5b4-1c451e426c29",
	}

	resp, err := repo.CreateTravelTips(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Content, resp.Content)
	assert.Equal(t, req.Category, resp.Category)
	assert.Equal(t, req.AuthorId, resp.AuthorId)
}

func TestGetTravelTips(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewCommunicationRepo(db)
	req := &pb.GetTravelTipsRequest{
		Page:     1,
		Limit:    10,
		Catygory: "Category",
	}

	resp, err := repo.GetTravelTips(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Tips)
}

func TestCountMessages(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewCommunicationRepo(db)

	total, err := repo.CountMessages()
	assert.NoError(t, err)
	assert.NotNil(t, total)
	assert.True(t, *total >= 0)
}
