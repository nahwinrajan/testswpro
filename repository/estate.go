package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"
)

const (
	// *** Estate ***
	queryGetEstateByID = `SELECT * FROM estates WHERE estate_id = $1`
	queryInsertEstate  = `
		INSERT INTO estates (width, length)
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
	queryGetTreeByEstateID = `SELECT * FROM trees WHERE estate_id = $1`
	queryInsertTree        = `
		INSERT INTO trees (estate_id, x, y, height)
		VALUES ($1, $2, $3, $4)
		RETURNING tree_id
	`
	queryDeleteTree = `DELETE FROM trees WHERE tree_id = $1`
)

// *** Estate ***
func (rp *Repository) GetEstateByID(ctx context.Context, estateID string) (output Estate, err error) {
	err = rp.db.QueryRowContext(ctx, "SELECT * FROM Estate WHERE estate_id = $1", estateID).Scan(output)
	return
}

func (rp *Repository) InsertEstate(ctx context.Context, width, length int) (estateID string, err error) {
	err = rp.db.QueryRowContext(
		ctx,
		queryInsertEstate,
		width,
		length,
	).Scan(&estateID)

	return
}

func (rp *Repository) UpdateEstate(ctx context.Context, count, min, max, median, patrol_distance int, patrol_route string) error {
	_, err := rp.db.ExecContext(
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
func (rp *Repository) GetAllTreesInEstate(ctx context.Context, estateID string) ([]Tree, error) {
	trees := make([]Tree, 0)

	rows, err := rp.db.QueryContext(ctx, queryGetTreeByEstateID, estateID)
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

func (rp *Repository) InsertTree(ctx context.Context, estateID string, x, y, height int) (treeID string, err error) {
	err = rp.db.QueryRowContext(
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

func (rp *Repository) DeleteTree(ctx context.Context, treeID string) error {
	_, err := rp.db.ExecContext(ctx, queryDeleteTree, treeID)

	return err
}
