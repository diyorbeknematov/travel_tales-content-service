package service

import (
	pb "content-service/generated/stories"
	"content-service/generated/user"
	"content-service/models"
	"content-service/storage/postgres"
	"context"
	"log/slog"
)

type TravelStoriesService struct {
	pb.UnimplementedTravelStoriesServiceServer
	StoriyRepo *postgres.TravelStoriesRepo
	UserClient user.AuthServiceClient
	Logger     *slog.Logger
}

func (s *TravelStoriesService) CreateTravelStory(cxt context.Context, in *pb.CreateTravelStoryRequest) (*pb.CreateTravelStoryResponse, error) {
	resp, err := s.StoriyRepo.CreateTravelStory(in)
	if err != nil {
		return nil, err
	}

	for _, v := range in.Tags {
		err := s.StoriyRepo.CreateStoryTags(models.StoryTag{
			StoryId: resp.Id,
			Tag:     v,
		})
		if err != nil {
			return nil, err
		}
	}

	resp.Tags = in.Tags

	return resp, nil
}

func (s *TravelStoriesService) UpdateTravelStory(ctx context.Context, in *pb.UpdateTravelStoryRequest) (*pb.UpdateTravelStoryResponse, error) {
	resp, err := s.StoriyRepo.UpdateTravelStory(in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TravelStoriesService) DeleteTravelStory(ctx context.Context, in *pb.DeleteTravelStoryRequest) (*pb.DeleteTravelStoryResponse, error) {
	resp, err := s.StoriyRepo.DeleteTravelStory(in.TravelStoryId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TravelStoriesService) ListTravelStory(ctx context.Context, in *pb.ListTravelStoryRequest) (*pb.ListTravelStoryResponse, error) {
	travelStories, err := s.StoriyRepo.GetTravelStories(in)
	if err != nil {
		s.Logger.Error("Xatolik sayohat hikoyalarini olshida", slog.String("error", err.Error()))
		return nil, err
	}

	for _, v := range travelStories.Stories {
		author, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: v.Author.Id})
		if err != nil {
			s.Logger.Error("Xatolik hikoyalarning userlarini olishda", slog.String("error", err.Error()))
			return nil, err
		}
		v.Author.Username = author.Username
	}

	return travelStories, nil
}

func (s *TravelStoriesService) GetTravelStory(ctx context.Context, in *pb.GetTravelStoryRequest) (*pb.GetTravelStoryResponse, error) {
	resp, err := s.StoriyRepo.GetTravelStory(in.StoryId)
	if err != nil {
		s.Logger.Error("hikoya haqida to'liq ma'lumot olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	author, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: resp.Author.Id})
	if err != nil {
		s.Logger.Error("Xatolik hikoyaning authorni olishda", slog.String("error", err.Error()))
		return nil, err
	}

	commentCount, err := s.StoriyRepo.CountComments(resp.Id)
	if err != nil {
		s.Logger.Error("sayohat hikoyasiga yozilgan commentlar sonini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	likeConunt, err := s.StoriyRepo.CountLikes(resp.Id)
	if err != nil {
		s.Logger.Error("sayohat hikoyasiga bosilgan likelar sonini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	resp.Author.Username = author.Username
	resp.Author.FullName = author.FullName
	resp.CommentsCount = commentCount
	resp.LikesCount = likeConunt

	return resp, nil
}

func (s *TravelStoriesService) AddCommment(ctx context.Context, in *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	resp, err := s.StoriyRepo.AddComment(in)
	if err != nil {
		s.Logger.Error("Xatolik hikoyaga izoh qoldirishda", slog.String("error", err.Error()))
		return nil, err
	}
	return resp, nil
}

func (s *TravelStoriesService) ListComments(ctx context.Context, in *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	comments, err := s.StoriyRepo.GetComments(in)
	if err != nil {
		s.Logger.Error("Xatolik hikoyalarga yozilgan izohlarni olishda", slog.String("error", err.Error()))
		return nil, err
	}

	for _, comment := range comments.Comments {
		author, err := s.UserClient.UserInfo(ctx, &user.UserInfoRequest{Id: comment.Author.Id})
		if err != nil {
			s.Logger.Error("Xatolik commentlarning authorlarini olishda", slog.String("error", err.Error()))
			return nil, err
		}
		comment.Author.Username = author.Username
	}

	return comments, nil
}

func (s *TravelStoriesService) AddLike(ctx context.Context, in *pb.AddLikeRequest) (*pb.AddLikeResponse, error) {
	resp, err := s.StoriyRepo.AddLike(in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *TravelStoriesService) CountComments(ctx context.Context, in *pb.CountCommentsRequest) (*pb.CountCommentsResponse, error) {
	resp, err := s.StoriyRepo.CountComments(in.UserId)
	if err != nil {
		s.Logger.Error("xatolik hikoyalarga yozilgan commentlar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.CountCommentsResponse{
		CountComments: int32(resp),
	}, nil
}

func (s *TravelStoriesService) CountLikes(ctx context.Context, in *pb.CountLikesRequest) (*pb.CountLikesResponse, error) {
	resp, err := s.StoriyRepo.CountLikes(in.UserId)
	if err != nil {
		s.Logger.Error("xatolik hikoyalarga bosilgan likelar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.CountLikesResponse{
		CountLikes: int32(resp),
	}, nil
}

func (s *TravelStoriesService) CountStories(ctx context.Context, in *pb.CountStoriesRequest) (*pb.CountStoriesResponse, error) {
	resp, err := s.StoriyRepo.CountComments(in.UserId)
	if err != nil {
		s.Logger.Error("xatolik hikoyalarga yozilgan commentlar sonini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.CountStoriesResponse{
		CountStories: int32(resp),
	}, nil
}