package service

import (
	pb "content-service/generated/communication"
	"content-service/generated/user"
	"content-service/storage/postgres"
	"context"
	"log/slog"
)

type CommunicationService struct {
	pb.UnimplementedCommunicationServiceServer
	CommunicationRepo *postgres.CommenicationRepo
	StoryRepo 		  *postgres.TravelStoriesRepo
	ItineraryRepo 	  *postgres.ItinerariesRepo
	UserClient        user.AuthServiceClient
	Logger            *slog.Logger
}

func (s *CommunicationService) SendMessageUser(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	resp, err := s.CommunicationRepo.SendMessage(in)
	if err != nil {
		s.Logger.Error("Userga xabar jo'natishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *CommunicationService) ListMessage(ctx context.Context, in *pb.ListMessageRequest) (*pb.ListMessageResponse, error) {
	resp, err := s.CommunicationRepo.GetMessages(in)
	if err != nil {
		s.Logger.Error("xabarlar ro'yxatini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	for _, v := range resp.Message {
		sender, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: v.Sender.Id})
		if err != nil {
			s.Logger.Error("xabar jo'natuvchini ma'lumotlarni olishda xatolik", slog.String("error", err.Error()))
			return nil, err
		}

		recipient, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: v.Recipient.Id})
		if err != nil {
			s.Logger.Error("Xabar jo'natuvchini ma'lumotlarini olishda xatolik", slog.String("error", err.Error()))
			return nil, err
		}

		v.Sender.Username = sender.Username
		v.Recipient.Username = recipient.Username
	}

	return resp, nil
}

func (s *CommunicationService) AddTravelTips(ctx context.Context, in *pb.AddTravelTipsRequest) (*pb.AddTravelTipsResponse, error) {
	resp, err := s.CommunicationRepo.CreateTravelTips(in)
	if err != nil {
		s.Logger.Error("yangi sayohat maslahatlarini qo'shishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *CommunicationService) GetTravelTips(ctx context.Context, in *pb.GetTravelTipsRequest) (*pb.GetTravelTipsResponse, error) {
	tips, err := s.CommunicationRepo.GetTravelTips(in)
	if err != nil {
		s.Logger.Error("sayohat manzillarini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	for _, v := range tips.Tips {
		author, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: v.Author.Id})
		if err != nil {
			s.Logger.Error("maslahat bergan odamning malumotlarini olishda xatolik", slog.String("error", err.Error()))
			return nil, err
		}

		v.Author.Username = author.Username
	}

	return tips, nil
}

func (s *CommunicationService) GetUserStatics(ctx context.Context, in *pb.GetUserStaticsRequest) (*pb.GetUserStaticsResponse, error) {
	totalStories, err := s.StoryRepo.CountStories(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning hikoyalar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	totalIntineraries, err := s.ItineraryRepo.CountItinerary(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning sayohat rejalarini sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	totalCountryViseted, err := s.ItineraryRepo.TotalCountriesVisited(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning sayohat qilgan davlatlar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	totalLikesReceived, err := s.StoryRepo.CountLikes(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning hikoyalariga bosilgan likelar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	totalCommentsReceived, err := s.StoryRepo.CountComments(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning hikoyalariga yozilgan commentlar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	mostPopularItinerary, err := s.ItineraryRepo.MostPopularItinerary(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning mashhur sayohat rejarlarini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	mostPopularStory, err := s.StoryRepo.MostPopularStory(in.UserId)
	if err != nil {
		s.Logger.Error("Xatolik userning mashhur sayohat hikoyalarini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.GetUserStaticsResponse{
		TotalStories: totalStories,
		TotalItineraries: totalIntineraries,
		TotalCountriesVisited: totalCountryViseted,
		TotalLikesReceived: totalLikesReceived,
		TotalCommentsReceived: totalCommentsReceived,
		MostPopularStory: mostPopularStory,
		MostPopularItinerary: mostPopularItinerary,
	}, nil
}