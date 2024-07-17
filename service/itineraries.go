package service

import (
	pb "content-service/generated/itineraries"
	"content-service/generated/user"
	"content-service/models"
	"content-service/storage/postgres"
	"context"
	"log/slog"
)

type ItineraryService struct {
	pb.UnimplementedItinerariesServiceServer
	ItineraryRepo *postgres.ItinerariesRepo
	Storyrepo     *postgres.TravelStoriesRepo
	UserClient    user.AuthServiceClient
	Logger        *slog.Logger
}

func (s *ItineraryService) CreateItinerary(ctx context.Context, in *pb.CreateItineraryRequest) (*pb.CreateItineraryResponse, error) {
	itinerary, err := s.ItineraryRepo.CreateItinerary(in)
	if err != nil {
		s.Logger.Error("sayohat rejasini tuzishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	for _, v := range in.Distinations {
		err = s.ItineraryRepo.CreateItineraryDestinations(models.ItineraryDestination{
			ItineraryId: itinerary.Id,
			Name:        v.Name,
			StartDate:   v.StartDate,
			EndDate:     v.EndDate,
		})
		if err != nil {
			s.Logger.Error("Xatolik intinerary_destinations tablega ma'lumot qo'shishda")
			return nil, err
		}
	}

	return itinerary, nil
}

func (s *ItineraryService) UpdateItinerary(ctx context.Context, in *pb.UpdateItineraryRequest) (*pb.UpdateItineraryResponse, error) {
	resp, err := s.ItineraryRepo.UpdateItinerary(in)
	if err != nil {
		s.Logger.Error("Sayohat rejasini yangilashda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *ItineraryService) DeleteItinerary(ctx context.Context, in *pb.DeleteItineraryRequest) (*pb.DeleteItineraryResponse, error) {
	resp, err := s.ItineraryRepo.DeleteItinerary(in.Id)
	if err != nil {
		s.Logger.Error("Sayohat rejasii o'chirishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	return resp, nil
}

func (s *ItineraryService) ListItineraries(ctx context.Context, in *pb.ListItinerariesRequest) (*pb.ListItinerariesResponse, error) {
	itinerary, err := s.ItineraryRepo.ListItineraries(in)
	if err != nil {
		s.Logger.Error("sayohat resjalarini ro'yxatini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	for _, v := range itinerary.Itineraries {
		author, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: v.Author.Id})
		if err != nil {
			s.Logger.Error("Sahoyat rejasini tuzgan inson malumotini olishda xatolik", slog.String("error", err.Error()))
			return nil, err
		}
		v.Author.Username = author.Username
	}

	return itinerary, nil
}

func (s *ItineraryService) GetItinerary(ctx context.Context, in *pb.GetItineraryRequest) (*pb.GetItineraryResponse, error) {
	itinerary, err := s.ItineraryRepo.GetItinerary(in.Id)
	if err != nil {
		s.Logger.Error("sayohat rejasi haqida to'liq ma'lumot olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	author, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: itinerary.Author.Id})
	if err != nil {
		s.Logger.Error("sahoyat rejasini tuzgan user haqida malumot olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	des, err := s.ItineraryRepo.GetItineraryDestinations(itinerary.Id)
	if err != nil {
		s.Logger.Error("Sayohat rejasidagi sayohat manzillarni olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	var destinations []*pb.Destination
	for _, v := range des {
		var des pb.Destination
		des.Name = v.Name
		des.StartDate = v.StartDate
		des.EndDate = v.EndDate
		act, err := s.ItineraryRepo.GetItineraryActivity(v.ID)
		if err != nil {
			s.Logger.Error("sayohat manzillarini activitylarini olishda xatolik", slog.String("error", err.Error()))
			return nil, err
		}
		des.Activities = act

		destinations = append(destinations, &des)
	}

	likeCount, err := s.Storyrepo.CountLikes(itinerary.Author.Id)
	if err != nil {
		s.Logger.Error("Likelar sonini topishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	commentCount, err := s.Storyrepo.CountComments(itinerary.Author.Id)
	if err != nil {
		s.Logger.Error("commentlar sonini topishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	itinerary.Author.Username = author.Username
	itinerary.Author.FullName = author.FullName
	itinerary.Destinations = destinations
	itinerary.CommentsCount = commentCount
	itinerary.LikesCount = likeCount

	return itinerary, nil
}

func (s *ItineraryService) LeaveComment(ctx context.Context, in *pb.LeaveCommentRequest) (*pb.LeaveCommentResponse, error) {
	resp, err := s.ItineraryRepo.CreateItineraryComments(in)
	if err != nil {
		s.Logger.Error("xatolik comment qoldirishda", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}
