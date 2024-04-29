// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetEstateByID(ctx context.Context, estateID string) (output Estate, err error)
	InsertEstate(ctx context.Context, width, length int) (estateID string, err error)
	UpdateEstate(ctx context.Context, count, min, max, median, patrol_distance int, patrol_route string) error
	GetAllTreesInEstate(ctx context.Context, estateID string) ([]Tree, error)
	InsertTree(ctx context.Context, estateID string, x, y, height int) (treeID string, err error)
	DeleteTree(ctx context.Context, treeID string) error
}
