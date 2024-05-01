package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/nahwinrajan/testswpro/uuidgen"
)

const (
	// *** Estate ***
	queryGetEstateByID = `SELECT 
		estate_id, width, length, count, min, max, median, patrol_distance, patrol_route
	 FROM estates 
	 WHERE estate_id = $1`

	queryInsertEstate = `
		INSERT INTO estates (estate_id, width, length)
		VALUES ($1, $2, $3)
	`
	queryUpdateEstateStats = `
		UPDATE estates
		SET
			count = $2,
			min = $3,
			max = $4,
			median = $5,
			patrol_distance = $6,
			patrol_route = $7,
			updated_at = now()
		WHERE
			estate_id = $1
	`

	// *** Tree ***
	queryGetTreeByEstateID = `SELECT
		tree_id, estate_id, x, y, height
	 FROM trees 
	 WHERE estate_id = $1`

	queryInsertTree = `
		INSERT INTO trees (estate_id, tree_id, x, y, height)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING tree_id
	`
	queryDeleteTree = `DELETE FROM trees WHERE tree_id = $1`
)

// *** Estate ***
func (rp *Repository) GetEstateByID(ctx context.Context, estateID string) (Estate, error) {
	var estate Estate

	err := rp.db.QueryRowContext(ctx, queryGetEstateByID, estateID).Scan(
		&estate.ID,
		&estate.Width,
		&estate.Length,
		&estate.Count,
		&estate.Min,
		&estate.Max,
		&estate.Median,
		&estate.PatrolDistance,
		&estate.PatrolRoute,
	)
	// *** DEBUGGING
	fmt.Printf("---***--- [Repo.GetEstateByID] returned estate: %+v\n", estate)
	fmt.Printf("---***--- [Repo.GetEstateByID] returned error: %+v\n", err)
	// *** END DEBUGGING
	return estate, err
}

func (rp *Repository) InsertEstate(ctx context.Context, width, length int) (string, error) {
	uuidEstateID, err := uuidgen.NewRandom()
	if err != nil {
		return "", err
	}

	_, err = rp.db.ExecContext(
		ctx,
		queryInsertEstate,
		uuidEstateID.String(),
		width,
		length,
	)

	return uuidEstateID.String(), err
}

func (rp *Repository) UpdateEstate(
	ctx context.Context,
	estateID string,
	count, min, max, median, patrolDistance int,
	patrolRoute string,
) error {
	_, err := rp.db.ExecContext(
		ctx,
		queryUpdateEstateStats,
		estateID,
		count,
		min,
		max,
		median,
		patrolDistance,
		patrolRoute,
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

func (rp *Repository) InsertTree(ctx context.Context, estateID string, x, y, height int) (string, error) {
	uuidTreeID, err := uuidgen.NewRandom()
	if err != nil {
		return "", err
	}

	_, err = rp.db.ExecContext(
		ctx,
		queryInsertTree,
		estateID,
		uuidTreeID.String(),
		x,
		y,
		height,
	)

	// *** DEBUGGING
	fmt.Printf("---***--- [Repo.InsertTree] returned tree_id:%s, error: %+v\n", uuidTreeID.String(), err)
	// // *** END DEBUGGING

	// extra steps for checking specific db error
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// *** DEBUGGING
			fmt.Printf("---***--- [Repo.InsertTree] pq error: %+v\n", pqErr)
			// // *** END DEBUGGING
			// Tree already exists at the specified location, return an error
			return "", errors.New("tree already exists at the specified location")
		}

		// Return other errors as is
		return "", err
	}

	return uuidTreeID.String(), nil
}

func (rp *Repository) DeleteTree(ctx context.Context, treeID string) error {
	_, err := rp.db.ExecContext(ctx, queryDeleteTree, treeID)

	return err
}
