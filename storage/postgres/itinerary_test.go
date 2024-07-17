package postgres

import (
	pb "content-service/generated/itineraries"
	"content-service/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateItinerary(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	req := &pb.CreateItineraryRequest{
		Title:       "Test Itinerary",
		Description: "Test Description",
		StartDate:   time.Now().Format("2006-01-02 15:04:05"),
		EndDate:     time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		AthorId:     "975799c4-bd72-43c8-b0c5-93bd9461e033",
	}

	resp, err := repo.CreateItinerary(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Title)
}

func TestUpdateItinerary(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	req := &pb.UpdateItineraryRequest{
		Id:          "4d4b7261-20a3-4f8c-aa78-958f12d4a9db",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	resp, err := repo.UpdateItinerary(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Description, resp.Description)
}

func TestDeleteItinerary(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	resp, err := repo.DeleteItinerary("4d4b7261-20a3-4f8c-aa78-958f12d4a9db")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Itinerary deleted succesfully", resp.Message)
}

func TestGetItinerary(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	resp, err := repo.GetItinerary("26d176f2-4536-443f-bf53-5cf25d4ffc65")

	expected := &pb.GetItineraryResponse{
		Id:          "26d176f2-4536-443f-bf53-5cf25d4ffc65",
		Title:       "Test Itinerary",
		Description: "Test Description",
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expected.Id, resp.Id)
	assert.Equal(t, expected.Title, resp.Title)
	assert.Equal(t, expected.Description, resp.Description)
	fmt.Println(expected.Author)
	assert.NotNil(t, expected.Author)
}

func TestListItineraries(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	req := &pb.ListItinerariesRequest{
		Page:  1,
		Limit: 10,
	}

	resp, err := repo.ListItineraries(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.LessOrEqual(t, len(resp.Itineraries), int(req.Limit))
}

func TestCreateItineraryDestinations(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	req := models.ItineraryDestination{
		ItineraryId: "26d176f2-4536-443f-bf53-5cf25d4ffc65",
		Name:        "Destination Name",
		StartDate:   time.Now().Format("2006-01-02 15:04:05"),
		EndDate:     time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	err = repo.CreateItineraryDestinations(req)
	assert.NoError(t, err)
}

func TestCreateItineraryActivity(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	req := models.ItineraryActivity{
		DestinationId: "03b1ba1e-5048-4f88-9f42-282e3333ba81",
		Activity:      "Activity Name",
	}

	err = repo.CreateItineraryActivity(req)
	assert.NoError(t, err)
}

func TestGetItineraryDestinations(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	resp, err := repo.GetItineraryDestinations("26d176f2-4536-443f-bf53-5cf25d4ffc65")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetItineraryActivity(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	resp, err := repo.GetItineraryActivity("03b1ba1e-5048-4f88-9f42-282e3333ba81")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCreateItineraryComments(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	req := &pb.LeaveCommentRequest{
		AuthorId:    "9b0cf2c8-308c-4896-a737-511bff1bb991",
		ItineraryId: "26d176f2-4536-443f-bf53-5cf25d4ffc65",
		Content:     "Test comment",
	}

	resp, err := repo.CreateItineraryComments(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.AuthorId, resp.AuthorId)
	assert.Equal(t, req.ItineraryId, resp.ItineraryId)
	assert.Equal(t, req.Content, resp.Content)
}

func TestCountItinerary(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	total, err := repo.CountItinerary("9b0cf2c8-308c-4896-a737-511bff1bb991")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int32(0))
}

func TestCountItineraryComments(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	total, err := repo.CountItineraryComments("26d176f2-4536-443f-bf53-5cf25d4ffc65")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int32(0))
}

func TestTotalCountriesVisited(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	total, err := repo.TotalCountriesVisited("9b0cf2c8-308c-4896-a737-511bff1bb991")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int32(0))
}

func TestMostPopularItinerary(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewItinerariesRepo(db)

	resp, err := repo.MostPopularItinerary("9b0cf2c8-308c-4896-a737-511bff1bb991")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
