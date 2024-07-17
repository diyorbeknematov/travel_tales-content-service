package service

import (
	pb "content-service/generated/destination"
	"content-service/generated/user"
	"content-service/storage/postgres"
	rdb "content-service/storage/redis"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type DestinationService struct {
	pb.UnimplementedTravelDestinationServiceServer
	DestinationRepo *postgres.DestinationRepo
	UserClient      *user.AuthServiceClient
	RedisClient     *rdb.RedisClient
	Logger          *slog.Logger
}

func (s *DestinationService) ListTravelDestnations(ctx context.Context, in *pb.ListDetinationRequest) (*pb.ListDetinationResponse, error) {
	resp, err := s.DestinationRepo.GetDestinations(in)
	if err != nil {
		s.Logger.Error("sayohat manzillarni ro'yxatini olish")
		return nil, err
	}

	for _, v := range resp.Destinations {
		activities, err := s.DestinationRepo.GetDestinationActivities(v.Id)
		if err != nil {
			s.Logger.Error("Xatolik sayohat manzilinig activitylarini olishda", slog.String("error", err.Error()))
			return nil, err
		}

		v.PopularActivities = activities
	}

	return resp, nil
}

func (s *DestinationService) GetTravelDestination(ctx context.Context, in *pb.GetDestinationRequest) (*pb.GetDestinationResponse, error) {
	resp, err := s.DestinationRepo.GetTravelDestination(in.Id)
	if err != nil {
		s.Logger.Error("Sayohat manzili haqaida malumot olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	activities, err := s.DestinationRepo.GetDestinationActivities(resp.Id)
	if err != nil {
		s.Logger.Error("Sayohat manzilining activitysini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	atractions, err := s.DestinationRepo.GetDestinationAttractions(resp.Id)
	if err != nil {
		s.Logger.Error("Sayohat manzilining atractions olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	resp.PopularActivities = activities
	resp.TopAttractions = atractions

	return resp, nil
}

func (s *DestinationService) GetTrendDestinations(ctx context.Context, in *pb.GetTrendDestinationRequest) (*pb.GetTrendDestinationResponse, error) {
	const cacheKey = "trending_destinations"
	// Redis cachedan olishga harakat qilamiz
	destinationsJSON, err := s.RedisClient.R.Get(ctx, cacheKey).Bytes()
	if err == redis.Nil {
		// Agar cacheda bo'lmasa, DBdan olamiz
		destinations, err := s.DestinationRepo.GetTrendingDestinations(int(in.Limit))
		if err != nil {
			s.Logger.Error("xatolik top sayohat manzillarini olishda", slog.String("error", err.Error()))
			return nil, err
		}

		// Cachega qo'shamiz
		destinationsJSON, _ = json.Marshal(destinations)
		s.RedisClient.R.Set(ctx, cacheKey, destinationsJSON, 10*time.Minute)

		return &pb.GetTrendDestinationResponse{
			Destinations: destinations.Destinations,
			Total:        destinations.Total,
		}, nil
	} else if err != nil {
		s.Logger.Error("xatolik  top sayohat manzillarini olishda redisdan", slog.String("error", err.Error()))
		return nil, err
	}

	// Cachedan olingan ma'lumotlarni deserializatsiya qilamiz
	var trendDestinations pb.GetTrendDestinationResponse
	json.Unmarshal(destinationsJSON, &trendDestinations)

	return &trendDestinations, nil
}
