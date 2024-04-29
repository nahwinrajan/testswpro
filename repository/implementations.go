package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"
)

const (
	// *** Estate ***
	queryGetEstateByID = `SELECT * FROM Estate WHERE estate_id = $1`
	queryInsertEstate  = `
		INSERT INTO Estate (width, length)
		VALUES ($1, $2)
		RETURNING estate_id
	`
	queryUpdateEstateStats = `
		UPDATE estates
		SET
			count = $1
			min = $2
			max = $3
			median = $4
			patrol_distance = $5
			patrol_route = $6
		WHERE
			estate_id = $7
	`

	// *** Tree ***
	queryGetTreeByEstateID = `SELECT * FROM Tree WHERE estate_id = $1`
	queryInsertTree        = `
		INSERT INTO Tree (estate_id, x, y, height)
		VALUES ($1, $2, $3, $4)
		RETURNING tree_id
	`
	queryDeleteTree = `DELETE FROM Tree WHERE tree_id = $1`
)

// *** Estate ***
func (r *Repository) GetEstateByID(ctx context.Context, estateID string) (output Estate, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM Estate WHERE estate_id = $1", estateID).Scan(output)
	return
}

func (r *Repository) InsertEstate(ctx context.Context, width, length int) (estateID string, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		queryInsertEstate,
		width,
		length,
	).Scan(&estateID)

	return
}

func (r *Repository) UpdateEstate(ctx context.Context, count, min, max, median, patrol_distance int, patrol_route string) error {
	_, err := r.Db.ExecContext(
		ctx,
		queryUpdateEstateStats,
		count,
		min,
		max,
		median,
		patrol_distance,
		patrol_route,
	)

	return err
}

// *** Tree ***
func (r *Repository) GetAllTreesInEstate(ctx context.Context, estateID string) ([]Tree, error) {
	trees := make([]Tree, 0)

	rows, err := r.Db.QueryContext(ctx, queryGetTreeByEstateID, estateID)
	if err != nil {
		return trees, err
	}
	defer rows.Close()

	for rows.Next() {
		var tree Tree
		if err := rows.Scan(&tree.ID, &tree.EstateID, &tree.X, &tree.Y, &tree.Height); err != nil {
			return trees, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

func (r *Repository) InsertTree(ctx context.Context, estateID string, x, y, height int) (treeID string, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		queryInsertEstate,
		estateID,
		x,
		y,
		height,
	).Scan(&treeID)

	// extra steps for checking specific db error
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// Tree already exists at the specified location, return an error
			return "", errors.New("tree already exists at the specified location")
		}

		// Return other errors as is
		return "", err
	}

	return
}

func (r *Repository) DeleteTree(ctx context.Context, treeID string) error {
	_, err := r.Db.ExecContext(ctx, queryDeleteTree, treeID)

	return err
}
