// This file contains types that are used in the repository layer.
package repository

import (
	"github.com/google/uuid"
)

type Estate struct {
	ID             uuid.UUID `db:"estate_id"`
	Width          int       `db:"width"`
	Length         int       `db:"length"`
	Count          int       `db:"count"`
	Min            int       `db:"min"`
	Max            int       `db:"max"`
	Median         int       `db:"median"`
	PatrolDistance int       `db:"patrol_distance"`
	PatrolRoute    string    `db:"patrol_route"`
}

type Tree struct {
	ID       uuid.UUID `db:"tree_id"`
	EstateID uuid.UUID `db:"estate_id"`
	X        int       `db:"x"`
	Y        int       `db:"y"`
	Height   int       `db:"height"`
}
