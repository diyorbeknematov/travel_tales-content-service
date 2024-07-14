package postgres

import (
	pb "content-service/generated/destination"
	"database/sql"
	"fmt"
	"log/slog"
)

type DestinationRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewDestinationRepo(db *sql.DB, logger *slog.Logger) *DestinationRepo {
	return &DestinationRepo{
		DB:     db,
		Logger: logger,
	}
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
			description,
		WHERE
			true `
	if req.Query != "" {
		query += fmt.Sprintf(" And name = $%d", ind)
		ind++
		args = append(args, req.Query)
	}
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", ind, ind+1)
	args = append(args, offset, req.Limit)

	rows, err := repo.DB.Query(query, args...)

	if err != nil {
		repo.Logger.Error("Error in get all destinations", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var dest pb.Destination
		err = rows.Scan(&dest.Id, &dest.Name, &dest.Country, &dest.Description)
		if err != nil {
			repo.Logger.Error("Error in scan destination", slog.String("error", err.Error()))
			return nil, err
		}

		resp = append(resp, &dest)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.ListDetinationResponse{
		Destinations: resp,
	}, nil
}

func (repo *DestinationRepo) GetTravelDestinations(id string) (*pb.GetDestinationResponse, error) {
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
	`, id).Scan(&resp.Id, &resp.Name, &resp.Country, &resp.Country, &resp.Description, &resp.BestTimeToVisit, &resp.AverageCostPerDay, &resp.Currency, &resp.Language)

	if err != nil {
		repo.Logger.Error("Error in get destinations", slog.String("error", err.Error()))
		return nil, err
	}

	return &resp, nil
}
