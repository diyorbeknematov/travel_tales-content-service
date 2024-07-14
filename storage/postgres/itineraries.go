package postgres

import (
	pb "content-service/generated/itineraries"
	"database/sql"
	"errors"
	"log/slog"
)

type ItinerariesRepo struct {
	DB *sql.DB
	Logger *slog.Logger
}

func NewItinerariesRepo(db *sql.DB, logger *slog.Logger) *ItinerariesRepo {
	return &ItinerariesRepo{
		DB: db,
		Logger: logger,
	}
}

func(repo *ItinerariesRepo) CreateItinerary(req *pb.CreateItineraryRequest) (*pb.CreateItineraryResponse, error) {
	var resp pb.CreateItineraryResponse

	err := repo.DB.QueryRow(`
		INSERT INTO itineraries (
			title,
			description,
			start_date,
			end_date,
			athor_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING
			id,
			title,
			description,
			start_date,
			end_date,
			author_id,
			created_at
	`, req.Title, req.Description, req.StartDate, req.EndDate, req.AthorId).
	Scan(&resp.Id, &resp.Title, &resp.Description, &resp.StartDate, &resp.EndDate, &resp.AuthorId, &resp.CreatedAt)

	if err != nil {
		repo.Logger.Error("Error in created itineraries", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *ItinerariesRepo) UpdateItinerary(req *pb.UpdateItineraryRequest) (*pb.UpdateItineraryResponse, error) {
	var resp pb.UpdateItineraryResponse

	err := repo.DB.QueryRow(`
		UPDATE 
			itineraries
		SET
			title,
			description 
		WHERE
			id = $1 and deleted_at = 0
		RETURNING
			id,
			title,
			description,
			start_date,
			end_date,
			author_id,
			updated_at
	`, req.Title, req.Description, req.Id).Scan(&resp.Id, &resp.Title, &resp.Description, &resp.StartDate, &resp.EndDate, &resp.AuthorId, &resp.UpdatedAt)

	if err != nil {
		repo.Logger.Error("Error in updated itinerary", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}

func (repo *ItinerariesRepo) DeleteItinerary(id string) (*pb.DeleteItineraryResponse, error) {
	// Itinerary-ni o'chirish uchun SQL so'rovi
	res, err := repo.DB.Exec(`
		DELETE FROM itineraries
		WHERE id = $1
	`, id)

	if err != nil {
		repo.Logger.Error("Error in deleting itinerary", slog.String("error", err.Error()))
		return nil, err
	}

	// Nechta qator o'chirilganini tekshirish
	rowsAffected, err := res.RowsAffected()
	if err != nil  {
		repo.Logger.Error("Error in getting rows affected", slog.String("error", err.Error()))
		return nil, err
	}

	if rowsAffected == 0 {
		repo.Logger.Error("Itinerary not found or already deleted")
		return &pb.DeleteItineraryResponse{
			Message: "Itinerary not found or already deleted",
		}, errors.New("itinerary not found or already deleted")
	}

	return &pb.DeleteItineraryResponse{
		Message: "Itinerary deleted succesfully",
	}, nil
}

func (repo *ItinerariesRepo) GetItinerary(id string) (*pb.GetItineraryResponse, error) {
	var resp pb.GetItineraryResponse
	var author pb.Author
	err := repo.DB.QueryRow(`
		SELECT
			id,
			title,
			description,
			start_date,
			end_date,
			author_id,
			created_at,
			updated_at
		FROM
			itineraries
		WHERE
			deleted_at = 0 id = $1
	`, id).Scan(&resp.Id, &resp.Title, &resp.Description, &resp.StartDate, &resp.EndDate, &author.Id, resp.CreatedAt, &resp.UpdatedAt)

	if err != nil {
		repo.Logger.Error("Error in get itinerary", slog.String("error", err.Error()))
		return nil, err
	}
	resp.Author = &author

	return &resp, nil
}

func (repo *ItinerariesRepo) ListItineraries(req *pb.ListItinerariesRequest) (*pb.ListItinerariesResponse, error) {
	var resp []*pb.Itinerary
	offset := (req.Page - 1) * req.Limit
	rows, err := repo.DB.Query(`
		SELECT
			id,
			title,
			author_id,
			start_date,
			end_date,
			created_at
		FROM
			itineraries
		WHERE
			deleted_at = 0
		OFFSET $1
		LIMIT $2
	`, offset, req.Limit)

	if err != nil {
		repo.Logger.Error("Error in get all itineraries", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var itinerary pb.Itinerary
		var author pb.Authors

		err = rows.Scan(&itinerary.Id, &itinerary.Title, &author.Id, &itinerary.StartDate, &itinerary.EndDate, &itinerary.CreatedAt)

		if err != nil {
			repo.Logger.Error("Error in scan itinerary", slog.String("error", err.Error()))
			return nil, err
		}
		itinerary.Author = &author

		resp = append(resp, &itinerary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.ListItinerariesResponse{
		Itineraries: resp,
	}, nil
}
