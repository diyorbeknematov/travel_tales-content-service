package postgres

import (
	pb "content-service/generated/destination"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDestination(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewDestinationRepo(db)
	req := &pb.AddDestinationRequest{
		Name:              "Test Destination",
		Country:           "Test Country",
		Description:       "Test Description",
		BestTimeToVisit:   "Test Time",
		AverageCostPerDay: 100,
		Currency:          "USD",
		Language:          "English",
	}

	resp, err := repo.CreateDestination(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.Country, resp.Country)
}

func TestGetDestinations(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewDestinationRepo(db)
	req := &pb.ListDetinationRequest{
		Page:  1,
		Limit: 10,
		Query: "Test Destination",
	}

	resp, err := repo.GetDestinations(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Total, int32(2))
}

func TestGetTravelDestination(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewDestinationRepo(db)
	id := "1"

	resp, err := repo.GetTravelDestination(id)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// assert.Equal(t, resp)
}

func TestGetDestinationActivities(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewDestinationRepo(db)
	id := "a068856b-f64a-4e68-8a3c-e37769bc760a"

	resp, err := repo.GetDestinationActivities(id)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetDestinationAttractions(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewDestinationRepo(db)
	id := ""

	resp, err := repo.GetDestinationAttractions(id)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetTrendingDestinations(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewDestinationRepo(db)
	limit := 5

	resp, err := repo.GetTrendingDestinations(limit)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
