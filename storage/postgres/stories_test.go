package postgres

import (
	pb "content-service/generated/stories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestCreateTravelStory(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	req := &pb.CreateTravelStoryRequest{
		Title:    "Title",
		Content:  "Content",
		Location: "Location",
		AuthorId: "6f645314-23f1-482e-bf83-417439ee582b",
	}

	resp, err := repo.CreateTravelStory(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Content, resp.Content)
	assert.Equal(t, req.Location, resp.Location)
	assert.Equal(t, req.AuthorId, resp.AuthorId)
	assert.NotZero(t, resp.Id)
	assert.NotZero(t, resp.CreatedAt)
}

func TestUpdateTravelStory(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	req := &pb.UpdateTravelStoryRequest{
		Id:      "4905038e-a27f-4386-9771-99cc10d2c6cd", 
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	resp, err := repo.UpdateTravelStory(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Content, resp.Content)
	assert.Equal(t, req.Id, resp.Id)
	assert.NotZero(t, resp.UpdatedAt)
}

func TestDeleteTravelStory(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	resp, err := repo.DeleteTravelStory("4905038e-a27f-4386-9771-99cc10d2c6cd")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "story successfully deleted", resp.Message)
}

func TestGetTravelStories(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	req := &pb.ListTravelStoryRequest{
		Page:  1,
		Limit: 2,
	}

	resp, err := repo.GetTravelStories(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, req.Page, resp.Page)
	assert.Equal(t, req.Limit, resp.Limit)
	assert.NotZero(t, resp.Total)
	assert.NotEmpty(t, resp.Stories)
}

func TestGetTravelStory(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	id := "fbff2f5f-decc-445a-a7d2-aa5df278f534"

	resp, err := repo.GetTravelStory(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, id, resp.Id)
	assert.NotEmpty(t, resp.Title)
	assert.NotEmpty(t, resp.Content)
	assert.NotEmpty(t, resp.Location)
	assert.NotZero(t, resp.CreatedAt)
}

func TestAddComment(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	req := &pb.AddCommentRequest{
		Content:  "Test Comment",
		AuthorId: "6f645314-23f1-482e-bf83-417439ee582b",
		StoryId:  "fbff2f5f-decc-445a-a7d2-aa5df278f534", 
	}

	resp, err := repo.AddComment(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, req.Content, resp.Content)
	assert.Equal(t, req.AuthorId, resp.AuthorId)
	assert.Equal(t, req.StoryId, resp.StoryId)
	assert.NotZero(t, resp.Id)
	assert.NotZero(t, resp.CreatedAt)
}

func TestGetComments(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	req := &pb.ListCommentsRequest{
		Page:  1,
		Limit: 2,
	}

	resp, err := repo.GetComments(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, req.Page, resp.Page)
	assert.Equal(t, req.Limit, resp.Limit)
	assert.NotZero(t, resp.Total)
	assert.NotEmpty(t, resp.Comments)
}

func TestAddLike(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	req := &pb.AddLikeRequest{
		UserId:  "6f645314-23f1-482e-bf83-417439ee582b",
		StoryId: "fbff2f5f-decc-445a-a7d2-aa5df278f534",
	}

	resp, err := repo.AddLike(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, req.UserId, resp.UserId)
	assert.Equal(t, req.StoryId, resp.StoryId)
	assert.NotZero(t, resp.LikedAt)
}

func TestCountStories(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	id := "6f645314-23f1-482e-bf83-417439ee582b" 

	total, err := repo.CountStories(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotZero(t, total)
}

func TestCountComments(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	id := "6f645314-23f1-482e-bf83-417439ee582b" 

	total, err := repo.CountComments(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotZero(t, total)
}

func TestCountLikes(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	id := "6f645314-23f1-482e-bf83-417439ee582b" 

	total, err := repo.CountLikes(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotZero(t, total)
}

func TestMostPopularStory(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewTravelStoriesRepo(db)

	id := "6f645314-23f1-482e-bf83-417439ee582b"

	resp, err := repo.MostPopularStory(id)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, resp.Id)
	assert.NotEmpty(t, resp.Title)
	assert.NotZero(t, resp.LikesCount)
}
