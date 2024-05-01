// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

// NOTE: Interface definition must be on package who need it.
// Ideally i want to do it like that, but that will also mean moving the data types
// into proper layer, let's save it for todo for time being.
// TODO: data structure into its own container in appropriate layer
// TODO: interface define on package needing it, not original package
type Repositorier interface {
	GetEstateByID(ctx context.Context, estateID string) (Estate, error)
	InsertEstate(ctx context.Context, width, length int) (estateID string, err error)
	UpdateEstate(ctx context.Context, estateID string, count, min, max, median, patrolDistance int, patrolRoute string) error
	GetAllTreesInEstate(ctx context.Context, estateID string) ([]Tree, error)
	InsertTree(ctx context.Context, estateID string, x, y, height int) (treeID string, err error)
	DeleteTree(ctx context.Context, treeID string) error
}
