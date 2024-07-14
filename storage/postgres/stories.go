package postgres

import (
	pb "content-service/generated/stories"
	"database/sql"
	"log/slog"
)

type TravelStoriesRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewTravelStoriesRepo(db *sql.DB, logger *slog.Logger) *TravelStoriesRepo {
	return &TravelStoriesRepo{
		DB:     db,
		Logger: logger,
	}
}

func (repo *TravelStoriesRepo) CreateTravelStory(req *pb.CreateTravelStoryRequest) (*pb.CreateTravelStoryResponse, error) {
	var resp pb.CreateTravelStoryResponse

	err := repo.DB.QueryRow(`
		INSERT INTO stories (
			title,
			content,
			location,
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
			location,
			author_id,
			created_at
	`, req.Title, req.Content, req.Location, req.AuthorId).
		Scan(&resp.Id, &resp.Title, &resp.Content,
			&resp.Location, &resp.AuthorId, &resp.CreatedAt)

	if err != nil {
		repo.Logger.Error("Error in travel story", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *TravelStoriesRepo) UpdateTravelStory(req *pb.UpdateTravelStoryRequest) (*pb.UpdateTravelStoryResponse, error) {
	var resp pb.UpdateTravelStoryResponse
	err := repo.DB.QueryRow(`
		UPDATE
			stories
		SET
			title = $1,
			content = $2
		WHERE
			id = $3 and deleted_at = 0
		RETURNING
			id,
			title,
			content,
			location,
			author_id,
			updated_at
	`, req.Title, req.Content, req.Id).Scan(&resp.Id, &resp.Title, &resp.Content,
		&resp.Location, &resp.AuthorId, &resp.UpdatedAt)

	if err != nil {
		repo.Logger.Error("Error in updated story", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *TravelStoriesRepo) DeleteTravelStory(id string) (*pb.DeleteTravelStoryResponse, error) {
	res, err := repo.DB.Exec(`
		UPDATE 
		stories
		SET
			deleted_at = DATE_PART('epoch', CURRENT_TIMESTAMP)::INT
		WHERE
			deleted_at = 0 and id = $1
	`)

	if err != nil {
		repo.Logger.Error("Error in story deletion", slog.String("error", err.Error()))
		return &pb.DeleteTravelStoryResponse{
			Message: "Error in travel story deletion",
		}, nil
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		repo.Logger.Error("Error getting rows affected", slog.String("error", err.Error()))
	}
	if rowsAffected == 0 {
		err := sql.ErrNoRows
		repo.Logger.Error("No story found to delete", slog.String("id", id))
		return nil, err
	}
	return &pb.DeleteTravelStoryResponse{
		Message: "story successfully deleted",
	}, nil
}

func (repo *TravelStoriesRepo) GetTravelStories(req *pb.ListTravelStoryRequest) (*pb.ListTravelStoryResponse, error) {
	var resp []*pb.TravelStory	
	offset := (req.Page - 1) * req.Limit

	rows, err := repo.DB.Query(`
		SELECT
			id,
			title,
			author_id,
			location,
			created_at
		FROM
			stories
		WHERE
			deleted_at = 0
		OFFSET $1
		LIMIT &2
	`, offset, req.Limit)

	if err != nil {
		repo.Logger.Error("Error get all stories", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var story pb.TravelStory
		var author pb.Authors

		err := rows.Scan(&story.Id, &story.Title, &author.Id, &story.Location, &story.CreatedAt)
		if err != nil {
			repo.Logger.Error("Error scan in variable", slog.String("error", err.Error()))
			return nil, err
		}

		story.Author = &author

		resp = append(resp, &story)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.ListTravelStoryResponse{
		Stories: resp,
	}, nil
}

func (repo *TravelStoriesRepo) GetTravelStory(id string) (*pb.GetTravelStoryResponse, error) {
	var resp pb.GetTravelStoryResponse
	var author pb.Author

	err := repo.DB.QueryRow(`
		SELECT
			id,
			title,
			content,
			location,
			author_id,
			created_at,
			uptadet_at
		FROM
			stories
		WHERE
			deleted_at = 0 AND id = $1
	`, id).Scan(&resp.Id, &resp.Title, &resp.Content, &resp.Location, &author.Id, &resp.CreatedAt, &resp.UpdatedAt)
	
	if err != nil {
		repo.Logger.Error("Error in get travel story", slog.String("error", err.Error()))
		return nil, err
	}
	resp.Author = &author

	return &resp, err
}

func (repo *TravelStoriesRepo) AddComment(req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	var resp pb.AddCommentResponse

	err := repo.DB.QueryRow(`
		INSERT INTO comments (
			content,
			author_id,
			story_id
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			content,
			author_id,
			story_id,
			created_at
	`, req.Content, req.AuthorId, req.StoryId).Scan(&resp.Id, &resp.Content, &resp.AuthorId, &resp.StoryId, &resp.CreatedAt)

	if err != nil {
		repo.Logger.Error("Error in added comment", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *TravelStoriesRepo) GetComments(req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	var resp []*pb.Comment

	rows, err := repo.DB.Query(`
		SELECT
			id,
			content,
			author_id,
			crated_at
		FROM
			comments c
		INNER JOIN
			stories s ON c.story_id = s.id 
		WHERE
			s.deleted_at = 0
	`)

	if err != nil {
		repo.Logger.Error("Error in get comments", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment pb.Comment
		var author pb.Authors
		
		err := rows.Scan(&comment.Id, &comment.Content, &author.Id, &comment.CreatedAt)
		if err != nil {
			repo.Logger.Error("Error scanning comment", slog.String("error", err.Error()))
			return nil, err
		}
		comment.Author = &author
		resp = append(resp, &comment)
	}

	if err = rows.Err(); err != nil {
		repo.Logger.Error("Rows error", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.ListCommentsResponse{Comments: resp}, nil
}

func (repo *TravelStoriesRepo) AddLike(req *pb.AddLikeRequest) (*pb.AddLikeResponse, error) {
	var resp pb.AddLikeResponse
	err := repo.DB.QueryRow(`
		INSERT INTO likes (
			user_id,
			story_id
		)
		VALUES (
			$1,
			$2
		)
		RETURNING
			user_id,
			story_id,
			created_at
	`, req.UserId, req.StoryId).Scan(&resp.UserId, &resp.StoryId, &resp.LikedAt)

	if err != nil {
		repo.Logger.Error("Error in scan like", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *TravelStoriesRepo) CountStories() (*int32, error) {
	var total int32
	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			stories
		WHEEW 
			deleted_at = 0
	`).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting stories", slog.String("error", err.Error()))
		return nil, err
	}

	return &total, nil
}

func (repo *TravelStoriesRepo) CountComments(id string) (*int32, error) {
	var total int32
	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			comments c
		JOIN 
			stories s ON c.story_id = s.id
		WHEEW 
			deleted_at = 0 and c.author_id = $1
	`, id).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting comments", slog.String("error", err.Error()))
		return nil, err
	}

	return &total, nil
}

func (repo *TravelStoriesRepo) CountLikes(id string) (*int32, error) {
	var total int32
	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			likes l
		JOIN 
			stories s ON l.story_id = s.id
		WHEEW 
			deleted_at = 0 and l.author_id = $1
	`, id).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting likes", slog.String("error", err.Error()))
		return nil, err
	}

	return &total, nil
}