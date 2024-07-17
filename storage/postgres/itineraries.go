package postgres

import (
	"content-service/generated/communication"
	pb "content-service/generated/itineraries"
	"content-service/models"
	"database/sql"
	"errors"
)

type ItinerariesRepo struct {
	DB *sql.DB
}

func NewItinerariesRepo(db *sql.DB) *ItinerariesRepo {
	return &ItinerariesRepo{
		DB: db,
	}
}

func (repo *ItinerariesRepo) CreateItinerary(req *pb.CreateItineraryRequest) (*pb.CreateItineraryResponse, error) {
	var resp pb.CreateItineraryResponse

	err := repo.DB.QueryRow(`
		INSERT INTO itineraries (
			title,
			description,
			start_date,
			end_date,
			author_id
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
			title = $1,
			description = $2
		WHERE
			id = $3 and deleted_at = 0
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
		return nil, err
	}

	// Nechta qator o'chirilganini tekshirish
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
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
			deleted_at = 0 AND id = $1
	`, id).Scan(&resp.Id, &resp.Title, &resp.Description, &resp.StartDate, &resp.EndDate, &author.Id, &resp.CreatedAt, &resp.UpdatedAt)

	if err != nil {
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
		return nil, err
	}

	for rows.Next() {
		var itinerary pb.Itinerary
		var author pb.Authors

		err = rows.Scan(&itinerary.Id, &itinerary.Title, &author.Id, &itinerary.StartDate, &itinerary.EndDate, &itinerary.CreatedAt)

		if err != nil {
			return nil, err
		}
		itinerary.Author = &author

		resp = append(resp, &itinerary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM
			itineraries
		WHERE
			deleted_at = 0
	`).Scan(&total)

	if err != nil {
		return nil, err
	}

	return &pb.ListItinerariesResponse{
		Itineraries: resp,
		Total:       total,
		Limit:       req.Limit,
		Page:        req.Page,
	}, nil
}

func (repo *ItinerariesRepo) CreateItineraryDestinations(req models.ItineraryDestination) error {
	_, err := repo.DB.Exec(`
		INSERT INTO itinerary_destinations (
			itinerary_id,
			name,
			start_date,
			end_date
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`, req.ItineraryId, req.Name, req.StartDate, req.EndDate)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ItinerariesRepo) CreateItineraryActivity(req models.ItineraryActivity) error {
	_, err := repo.DB.Exec(`
		INSERT INTO itinerary_activities (
			destination_id,
			activity
		)
		VALUES (
			$1,
			$2
		)
	`, req.DestinationId, req.Activity)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ItinerariesRepo) GetItineraryDestinations(id string) ([]models.Result, error) {
	var destinations []models.Result
	rows, err := repo.DB.Query(`
		SELECT
			id,
			name,
			start_date,
			end_date
		FROM
			itinerary_destinations i_d
		JOIN 
			itineraries i ON i.id = i_d.itinerary_id
		WHERE
			i.deleted_at = 0 and i_d.itinerary_id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res models.Result

		err = rows.Scan(res.ID, &res.Name, &res.StartDate, &res.EndDate)
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return destinations, nil
}

func (repo *ItinerariesRepo) GetItineraryActivity(id string) ([]string, error) {
	var destinations []string

	rows, err := repo.DB.Query(`
		SELECT
			activity
		FROM
			itinerary_activities 
		WHERE
			destination_id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res string

		err = rows.Scan(&res)
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return destinations, nil
}

func (repo *ItinerariesRepo) CreateItineraryComments(req *pb.LeaveCommentRequest) (*pb.LeaveCommentResponse, error) {
	var resp pb.LeaveCommentResponse

	err := repo.DB.QueryRow(`
		INSERT INTO itinerary_comments (
			author_id,
			itinerary_id,
			content
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
			itinerary_id,
			created_at
	`, req.AuthorId, req.ItineraryId, req.Content).Scan(&resp.Id, &resp.Content, &resp.AuthorId, &resp.ItineraryId, &resp.CreatedAt)

	return &resp, err
}

func (repo *ItinerariesRepo) CountItinerary(id string) (int32, error) {
	var total int32
	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			itineraries
		WHERE
			(deleted_at = 0) AND (author_id = $1 or id = $1) 
	`, id).Scan(&total)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func (repo *ItinerariesRepo) CountItineraryComments(id string) (int32, error) {
	var total int32
	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			itineraries
		WHERE
			(deleted_at = 0) AND (author_id = $1 or id = $1) 
	`, id).Scan(&total)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func (repo *ItinerariesRepo) TotalCountriesVisited(id string) (int32, error) {
	var total int32

	err := repo.DB.QueryRow(`
		SELECT 
			COUNT(DISTINCT id) 
		FROM 
			itinerary_destinations id
		JOIN 
			itineraries i ON id.itinerary_id = i.id
		WHERE 
			i.author_id = $1 AND i.deleted_at = 0;
	`, id).Scan(&total)

	return total, err
}

func (repo *ItinerariesRepo) MostPopularItinerary(id string) (*communication.MostPopularItinerary, error) {
	var resp communication.MostPopularItinerary

	err := repo.DB.QueryRow(`
		SELECT 
			id, 
			title, 
			likes_count 
		FROM 
			itineraries 
		WHERE 
			author_id = $1 AND deleted_at = 0 
		ORDER BY 
			likes_count DESC 
		LIMIT 1;
	`, id).Scan(&resp.Id, &resp.Title, &resp.LikesCount)

	return &resp, err
}
