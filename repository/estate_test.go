package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetEstateByID(t *testing.T) {
	tests := []struct {
		name           string
		expectedEstate Estate
		expectedErr    error
	}{
		{
			name: "Valid estate ID",
			expectedEstate: Estate{
				ID:     "estate_id_value",
				Width:  10,
				Length: 20,
			},
			expectedErr: nil,
		},
		{
			name:           "Empty estate ID",
			expectedEstate: Estate{},
			expectedErr:    errors.New("estate ID is empty"),
		},
		{
			name:           "Estate not found",
			expectedEstate: Estate{},
			expectedErr:    sql.ErrNoRows,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database and repository
			dbmock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbmock.Close()

			repo := Repository{
				db: dbmock,
			}

			queryPattern := `SELECT .* FROM estates WHERE estate_id \= \$1`
			if tc.expectedErr != nil {
				mock.ExpectQuery(queryPattern).WithArgs(tc.expectedEstate.ID).WillReturnError(tc.expectedErr)
			} else {
				mock.ExpectQuery(queryPattern).WithArgs(tc.expectedEstate.ID).WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"estate_id", "width", "length", "count", "min", "max", "median", "patrol_distance", "patrol_route",
						}).
						AddRow(
							tc.expectedEstate.ID,
							tc.expectedEstate.Width,
							tc.expectedEstate.Length,
							tc.expectedEstate.Count,
							tc.expectedEstate.Min,
							tc.expectedEstate.Max,
							tc.expectedEstate.Median,
							tc.expectedEstate.PatrolDistance,
							tc.expectedEstate.PatrolRoute,
						),
				)
			}

			// Call the function under test
			estate, err := repo.GetEstateByID(context.Background(), tc.expectedEstate.ID)

			// Verify the result
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedEstate, estate)

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertEstate(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		length      int
		expectedID  string
		expectedErr error
	}{
		{
			name:        "Valid width and length",
			width:       10,
			length:      20,
			expectedID:  "lets-pretend-this-is-uuid",
			expectedErr: nil,
		},
		{
			name:        "Zero width",
			width:       0,
			length:      20,
			expectedID:  "",
			expectedErr: errors.New("width must be greater than zero"),
		},
		{
			name:        "Negative width",
			width:       -5,
			length:      20,
			expectedID:  "",
			expectedErr: errors.New("width must be greater than zero"),
		},
		{
			name:        "Zero length",
			width:       10,
			length:      0,
			expectedID:  "",
			expectedErr: errors.New("length must be greater than zero"),
		},
		{
			name:        "Negative length",
			width:       10,
			length:      -5,
			expectedID:  "",
			expectedErr: errors.New("length must be greater than zero"),
		},
		{
			name:        "Database error",
			width:       10,
			length:      20,
			expectedID:  "",
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database and repository
			dbmock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbmock.Close()

			repo := Repository{
				db: dbmock,
			}

			queryPattern := `INSERT INTO estates \(estate_id, width, length\) VALUES \(\$1, \$2, \$3\)`
			if tc.expectedErr != nil {
				mock.ExpectExec(queryPattern).WithArgs(sqlmock.AnyArg(), tc.width, tc.length).WillReturnError(tc.expectedErr)
			} else {
				mock.ExpectExec(queryPattern).WithArgs(sqlmock.AnyArg(), tc.width, tc.length).WillReturnResult(sqlmock.NewResult(0, 1))
			}

			// Call the function under test
			estateID, err := repo.InsertEstate(context.Background(), tc.width, tc.length)

			// Verify the result
			require.Equal(t, tc.expectedErr, err)
			_, err = uuid.Parse(estateID)
			require.NoError(t, err)

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateEstate(t *testing.T) {
	tests := []struct {
		name           string
		estateID       string
		count          int
		min            int
		max            int
		median         int
		patrolDistance int
		patrolRoute    string
		expectedErr    error
	}{
		{
			name:           "Valid parameters",
			estateID:       "estate_id_value",
			count:          5,
			min:            1,
			max:            10,
			median:         5,
			patrolDistance: 100,
			patrolRoute:    "route",
			expectedErr:    nil,
		},
		{
			name:           "Invalid estate ID",
			estateID:       "",
			count:          5,
			min:            1,
			max:            10,
			median:         5,
			patrolDistance: 100,
			patrolRoute:    "route",
			expectedErr:    errors.New("estate ID is empty"),
		},
		{
			name:           "Database error",
			estateID:       "estate_id_value",
			count:          5,
			min:            1,
			max:            10,
			median:         5,
			patrolDistance: 100,
			patrolRoute:    "route",
			expectedErr:    errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database and repository
			dbmock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbmock.Close()

			repo := Repository{
				db: dbmock,
			}

			queryPattern := `
				UPDATE estates
				SET
					count = \$2,
					min = \$3,
					max = \$4,
					median = \$5,
					patrol_distance = \$6,
					patrol_route = \$7,
					updated_at = now\(\)
				WHERE
					estate_id = \$1
			`
			if tc.expectedErr != nil {
				mock.ExpectExec(queryPattern).
					WithArgs(tc.estateID, tc.count, tc.min, tc.max, tc.median, tc.patrolDistance, tc.patrolRoute).
					WillReturnError(tc.expectedErr)
			} else {
				mock.ExpectExec(queryPattern).
					WithArgs(tc.estateID, tc.count, tc.min, tc.max, tc.median, tc.patrolDistance, tc.patrolRoute).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}

			// Call the function under test
			err = repo.UpdateEstate(
				context.Background(),
				tc.estateID,
				tc.count,
				tc.min,
				tc.max,
				tc.median,
				tc.patrolDistance,
				tc.patrolRoute,
			)

			// Verify the result
			require.Equal(t, tc.expectedErr, err)

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllTreesInEstate(t *testing.T) {
	tests := []struct {
		name          string
		estateID      string
		expectedTrees []Tree
		expectedErr   error
	}{
		{
			name:     "Valid estate ID with trees",
			estateID: "estate_id_value",
			expectedTrees: []Tree{
				{ID: "tree_id_1", EstateID: "estate_id_value", X: 1, Y: 2, Height: 10},
				{ID: "tree_id_2", EstateID: "estate_id_value", X: 3, Y: 4, Height: 12},
			},
			expectedErr: nil,
		},
		{
			name:          "Invalid estate ID",
			estateID:      "",
			expectedTrees: nil,
			expectedErr:   errors.New("estate ID is empty"),
		},
		{
			name:          "Database error",
			estateID:      "estate_id_value",
			expectedTrees: nil,
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database and repository
			dbmock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbmock.Close()

			repo := Repository{
				db: dbmock,
			}

			rows := sqlmock.NewRows([]string{"tree_id", "estate_id", "x", "y", "height"})
			for _, tree := range tc.expectedTrees {
				rows.AddRow(tree.ID, tree.EstateID, tree.X, tree.Y, tree.Height)
			}

			queryPattern := `SELECT .* FROM trees WHERE estate_id = \$1`
			if tc.expectedErr != nil {
				mock.ExpectQuery(queryPattern).WithArgs(tc.estateID).WillReturnError(tc.expectedErr)
			} else {
				mock.ExpectQuery(queryPattern).WithArgs(tc.estateID).WillReturnRows(rows)
			}

			// Call the function under test
			trees, err := repo.GetAllTreesInEstate(context.Background(), tc.estateID)

			// Verify the result
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedTrees, trees)

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertTree(t *testing.T) {
	tests := []struct {
		name          string
		estateID      string
		x             int
		y             int
		height        int
		expectedID    string
		expectedErr   error
		expectedPQErr *pq.Error
	}{
		{
			name:        "Valid tree insertion",
			estateID:    "estate_id_value",
			x:           1,
			y:           2,
			height:      10,
			expectedID:  "we'll-just-validate-it's-valid-uuid",
			expectedErr: nil,
		},
		{
			name:        "Unique constraint violation",
			estateID:    "estate_id_value",
			x:           1,
			y:           2,
			height:      10,
			expectedID:  "",
			expectedErr: errors.New("tree already exists at the specified location"),
			expectedPQErr: &pq.Error{
				Code: "23505",
			},
		},
		{
			name:          "Database error",
			estateID:      "estate_id_value",
			x:             1,
			y:             2,
			height:        10,
			expectedID:    "",
			expectedErr:   errors.New("database error"),
			expectedPQErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database and repository
			dbmock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbmock.Close()

			repo := Repository{
				db: dbmock,
			}

			queryPattern := `INSERT INTO trees \(estate_id, tree_id, x, y, height\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING tree_id`
			if tc.expectedErr != nil {
				mock.ExpectExec(queryPattern).WithArgs(tc.estateID, sqlmock.AnyArg(), tc.x, tc.y, tc.height).WillReturnError(tc.expectedErr)
			} else if tc.expectedPQErr != nil {
				mock.ExpectExec(queryPattern).WithArgs(tc.estateID, sqlmock.AnyArg(), tc.x, tc.y, tc.height).WillReturnError(tc.expectedPQErr)
			} else {
				mock.ExpectExec(queryPattern).WithArgs(tc.estateID, sqlmock.AnyArg(), tc.x, tc.y, tc.height).WillReturnResult(sqlmock.NewResult(0, 1))
			}

			// Call the function under test
			treeID, err := repo.InsertTree(context.Background(), tc.estateID, tc.x, tc.y, tc.height)

			// Verify the result
			require.Equal(t, tc.expectedErr, err)
			if tc.expectedPQErr != nil {
				require.NotNil(t, err)
				require.Equal(t, tc.expectedErr, err)
			}

			if treeID != "" {
				_, err = uuid.Parse(treeID)
				require.NoError(t, err)
			}

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteTree(t *testing.T) {
	tests := []struct {
		name        string
		treeID      string
		expectedErr error
	}{
		{
			name:        "Valid tree ID",
			treeID:      "tree_id_value",
			expectedErr: nil,
		},
		{
			name:        "Tree not found",
			treeID:      "non_existing_tree_id",
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database and repository
			dbmock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbmock.Close()

			repo := Repository{
				db: dbmock,
			}

			queryPattern := `DELETE FROM trees WHERE tree_id = \$1`
			if tc.expectedErr != nil {
				mock.ExpectExec(queryPattern).WithArgs(tc.treeID).WillReturnError(tc.expectedErr)
			} else {
				mock.ExpectExec(queryPattern).WithArgs(tc.treeID).WillReturnResult(sqlmock.NewResult(0, 1))
			}

			// Call the function under test
			err = repo.DeleteTree(context.Background(), tc.treeID)

			// Verify the result
			require.Equal(t, tc.expectedErr, err)

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
