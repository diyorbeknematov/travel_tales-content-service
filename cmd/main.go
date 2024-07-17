package main

import (
	"content-service/cmd/server"
	"content-service/config"
	"content-service/generated/communication"
	"content-service/generated/destination"
	"content-service/generated/itineraries"
	"content-service/generated/stories"
	"content-service/logs"
	"content-service/service"
	"content-service/storage/postgres"
	"content-service/storage/redis"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

func main() {
	logs.InitLogger()
	logs.Logger.Info("Starting the server ...")
	db, err := postgres.ConnectDB()
	if err != nil {
		logs.Logger.Error("Error connection th postgres", slog.String("error", err.Error()))
		log.Fatal(err)
	}
	defer db.Close()

	cfg := config.Load()
	listener, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		logs.Logger.Error("Error create to new listener", "error", err.Error())
		log.Fatal(err)
	}

	userClient, err := server.NewUserClient(cfg)
	if err != nil {
		logs.Logger.Error("Error in conn userclient", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	s := grpc.NewServer()
	stories.RegisterTravelStoriesServiceServer(s, &service.TravelStoriesService{
		StoriyRepo: postgres.NewTravelStoriesRepo(db),
		Logger:     logs.Logger,
		UserClient: userClient,
	})

	itineraries.RegisterItinerariesServiceServer(s, &service.ItineraryService{
		ItineraryRepo: postgres.NewItinerariesRepo(db),
		Storyrepo:     postgres.NewTravelStoriesRepo(db),
		UserClient:    userClient,
		Logger:        logs.Logger,
	})

	destination.RegisterTravelDestinationServiceServer(s, &service.DestinationService{
		DestinationRepo: postgres.NewDestinationRepo(db),
		UserClient:      &userClient,
		Logger:          logs.Logger,
		RedisClient:     redis.NewRedisClient(),
	})

	communication.RegisterCommunicationServiceServer(s, &service.CommunicationService{
		CommunicationRepo: postgres.NewCommunicationRepo(db),
		StoryRepo:         postgres.NewTravelStoriesRepo(db),
		ItineraryRepo:     postgres.NewItinerariesRepo(db),
		UserClient:        userClient,
		Logger:            logs.Logger,
	})

	logs.Logger.Info("server is running ", "PORT", cfg.GRPC_PORT)

	log.Printf("server is running on %v...", listener.Addr())
	if err := s.Serve(listener); err != nil {
		logs.Logger.Error("Faild server is running", "error", err.Error())
		log.Fatal(err)
	}
}
