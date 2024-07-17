package postgres

import (
	pb "content-service/generated/destination"
	"database/sql"
	"fmt"
)

type DestinationRepo struct {
	DB *sql.DB
}

func NewDestinationRepo(db *sql.DB) *DestinationRepo {
	return &DestinationRepo{
		DB: db,
	}
}

func (repo *DestinationRepo) CreateDestination(req *pb.AddDestinationRequest) (*pb.AddDestionationResponse, error) {
	var resp pb.AddDestionationResponse

	err := repo.DB.QueryRow(`
		INSERT INTO destinations (
			name,
			country,
			description,
			best_time_to_visit,
			average_cost_per_day,
			currency,
			language
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		RETURNING
			id,
			name,
			country,
			description,
			best_time_to_visit,
			average_cost_per_day,
			currency,
			language,
			created_at
	`, req.Name, req.Country, req.Description, req.BestTimeToVisit, req.AverageCostPerDay, req.Currency, req.Language).
		Scan(&resp.Id, &resp.Name, &resp.Country, &resp.Description, &resp.BestTimeToVisit, &resp.AverageCostPerDay, &resp.Currency, &resp.Language, &resp.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (repo *DestinationRepo) GetDestinations(req *pb.ListDetinationRequest) (*pb.ListDetinationResponse, error) {
	var resp []*pb.Destination
	offset := (req.Page - 1) * req.Limit
	var args []interface{}
	ind := 1
	query := `
		SELECT
			id,
			name,
			country,
			description
		FROM
			destinations
		WHERE
			deleted_at = 0 `
	if req.Query != "" {
		query += fmt.Sprintf(" And name = $%d", ind)
		ind++
		args = append(args, req.Query)
	}
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", ind, ind+1)
	args = append(args, offset, req.Limit)

	rows, err := repo.DB.Query(query, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dest pb.Destination
		err = rows.Scan(&dest.Id, &dest.Name, &dest.Country, &dest.Description)
		if err != nil {
			return nil, err
		}

		resp = append(resp, &dest)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM
			destinations
		WHERE
			deleted_at = 0
	`).Scan(&total)

	if err != nil {
		return nil, err
	}

	return &pb.ListDetinationResponse{
		Destinations: resp,
		Total:        total,
		Limit:        req.Limit,
		Page:         req.Page,
	}, nil
}

func (repo *DestinationRepo) GetTravelDestination(id string) (*pb.GetDestinationResponse, error) {
	var resp pb.GetDestinationResponse

	err := repo.DB.QueryRow(`
		SELECT
			id, 
			name,
			country,
			description,
			best_time_to_visit,
			average_cost_per_day,
			currency,
			language
		FROM
			destinations
		WHERE
			deleted_at = 0 and id = $1
	`, id).Scan(&resp.Id, &resp.Name, &resp.Country, &resp.Description, &resp.BestTimeToVisit, &resp.AverageCostPerDay, &resp.Currency, &resp.Language)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (repo *DestinationRepo) GetDestinationActivities(id string) ([]string, error) {
	var resp []string

	rows, err := repo.DB.Query(`
		SELECT
			activity
		FROM
			destination_activities da
		INNER JOIN
			destinations d ON da.destination_id = d.id
		WHERE
			da.destination_id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var activity string

		err = rows.Scan(&activity)
		if err != nil {
			return nil, err
		}

		resp = append(resp, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (repo *DestinationRepo) GetDestinationAttractions(id string) ([]string, error) {
	rows, err := repo.DB.Query(`
        SELECT 
            attraction 
        FROM 
            top_attractions
        WHERE 
            destination_id = $1
    `, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attractions []string
	for rows.Next() {
		var attraction string
		if err := rows.Scan(&attraction); err != nil {
			return nil, err
		}
		attractions = append(attractions, attraction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attractions, nil
}

func (repo *DestinationRepo) GetTrendingDestinations(limit int) (*pb.GetTrendDestinationResponse, error) {
    rows, err := repo.DB.Query(`
        SELECT 
            id, 
            name, 
            country, 
            popularity_score
        FROM 
            destinations 
        ORDER BY 
            popularity_score DESC 
        LIMIT $1
    `, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var destinations []*pb.TrendDestination
    for rows.Next() {
        var dest pb.TrendDestination
        if err := rows.Scan(&dest.Id, &dest.Name, &dest.Country, &dest.PopularityScore); err != nil {
            return nil, err
        }
        destinations = append(destinations, &dest)
    }

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			destinations
		WHERE
			deleted_at = 0
	`).Scan(&total)

	if err != nil {
		return nil, err
	}

    return &pb.GetTrendDestinationResponse{
		Destinations: destinations,
		Total: total,
	}, nil
}