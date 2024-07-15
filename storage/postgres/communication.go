package postgres

import (
	pb "content-service/generated/communication"
	"database/sql"
	"fmt"
	"log/slog"
)

type CommenicationRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewCommunicationRepo(db *sql.DB, logger *slog.Logger) *CommenicationRepo {
	return &CommenicationRepo{
		DB:     db,
		Logger: logger,
	}
}

func (repo *CommenicationRepo) SendMessage(req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	var resp pb.SendMessageResponse

	err := repo.DB.QueryRow(`
		INSERT INTO messages (
			sender_id,
			recipient_id,
			content
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			sender_id,
			recipient_id,
			content,
			created_at
	`, req.SendeId, req.RecipientId, req.Content).Scan(&resp.Id, &resp.SenderId, &resp.RecipientId, &resp.Content, &resp.CreatedAt)

	if err != nil {
		repo.Logger.Error("Error in send message", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *CommenicationRepo) GetMessages(req *pb.ListMessageRequest) (*pb.ListMessageResponse, error) {
	var resp []*pb.Message
	offset := (req.Page - 1) * req.Limit

	rows, err := repo.DB.Query(`
		SELECT
			id,
			sender_id,
			recipient_id,
			content,
			created_at
		FROM
			messages
		OFFSET $1
		LIMIT $2
	`, offset, req.Limit)

	if err != nil {
		repo.Logger.Error("Error in get all messages", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var sender pb.Sender
		var recipient pb.Recipient
		var message pb.Message

		err = rows.Scan(&message.Id, &sender.Id, &recipient.Id, &message.Content, &message.CreatedAt)
		if err != nil {
			repo.Logger.Error("Error in scan message", slog.String("error", err.Error()))
			return nil, err
		}
		message.Sender = &sender
		message.Recipient = &recipient

		resp = append(resp, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			messages
	`, ).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting messages", slog.String("error", err.Error()))
		return nil, err
	}


	return &pb.ListMessageResponse{
		Message: resp,
		Limit: req.Limit,
		Page: req.Page,
		Total: total,
	}, nil
}

func (repo *CommenicationRepo) CreateTravelTips(req *pb.AddTravelTipsRequest) (*pb.AddTravelTipsResponse, error) {
	var resp pb.AddTravelTipsResponse

	err := repo.DB.QueryRow(`
		INSERT INTO travel_tips (
			title,
			content,
			category,
			author_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id,
			title,
			content,
			category,
			athor_id,
			created_at
	`, req.Title, req.Content, req.Category, req.AuthorId).
	Scan(&resp.Id, &resp.Title, resp.Content, &resp.Category, &resp.AuthorId, &resp.CreatedAt)

	if err != nil {
		repo.Logger.Error("Error in add travel tips", slog.String("error", err.Error()))
		return nil, err
	}
	return &resp, nil
}

func (repo *CommenicationRepo) GetTravelTips(req *pb.GetTravelTipsRequest) (*pb.GetTravelTipsResponse, error) {
	var resp []*pb.Tip
	var args []interface{}
	offset := (req.Page - 1) * req.Limit
	query := `
		SELECT
			id,
			title,
			category,
			author_id,
			created_at
		WHERE
			true `
	ind := 1

	if req.Catygory != "" {
		query += fmt.Sprintf(" AND catygory = $%d", ind)
		ind++
		args = append(args, req.Catygory)
	} 
		query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", ind, ind+1)
		args = append(args, offset, req.Limit)
	

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		repo.Logger.Error("Error in get all travels", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var tip pb.Tip 
		var author pb.Author
		err = rows.Scan(&tip.Id, &tip.Title, &tip.Category, &author.Id, &tip.CreatedAt)
		if err != nil {
			repo.Logger.Error("Error in add get all tips", slog.String("error", err.Error()))
			return nil, err
		}

		tip.Author = &author

		resp = append(resp, &tip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM
			travel_tips
	`).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting travel_tips", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.GetTravelTipsResponse{
		Tips: resp,
		Limit: req.Limit,
		Page: req.Page,
		Total: total,
	}, nil
}

func (repo *CommenicationRepo) CountMessages() (*int32, error) {
	var total int32
	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			messages 
	`, ).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting messages", slog.String("error", err.Error()))
		return nil, err
	}

	return &total, nil
}