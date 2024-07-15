package models

type ItineraryDestination struct {
	ID          string
	ItineraryId string
	Name        string
	StartDate   string
	EndDate     string
}

type ItineraryActivity struct {
	ID            string
	DestinationId string
	Activity      string
}

type StoryTag struct {
	StoryId string
	Tag     string
}

type Result struct {
	ID        string
	Name      string
	StartDate string
	EndDate   string
}
