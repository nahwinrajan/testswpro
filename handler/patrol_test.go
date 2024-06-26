package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/nahwinrajan/testswpro/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPatrol(t *testing.T) {
	tests := []struct {
		name             string
		estate           repository.Estate
		trees            []repository.Tree
		expectedMin      int
		expectedMax      int
		expectedMedian   int
		expectedDistance int
		expectedPath     string
		expectedError    error
	}{
		{
			name: "Valid estate with trees",
			estate: repository.Estate{
				Width:  5,
				Length: 1,
			},
			trees: []repository.Tree{
				{X: 2, Y: 1, Height: 5},
				{X: 3, Y: 1, Height: 3},
				{X: 4, Y: 1, Height: 4},
			},
			expectedMin:      3,
			expectedMax:      5,
			expectedMedian:   4,
			expectedDistance: 64,
			expectedPath:     "1,1,1,ew,0,10;2,2,1,vu,6,16;3,2,1,ew,6,26;4,3,1,vd,2,28;5,3,1,ew,2,38;6,4,1,vu,1,39;7,4,1,ew,1,49;8,5,1,ew,1,59;",
			expectedError:    nil,
		},
		{
			name:   "Empty estate",
			estate: repository.Estate{},
			trees: []repository.Tree{
				{X: 2, Y: 1, Height: 5},
				{X: 3, Y: 1, Height: 3},
				{X: 4, Y: 1, Height: 4},
			},
			expectedError: errors.New("invalid estate value"),
		},
		{
			name: "Empty tree list",
			estate: repository.Estate{
				Width:  5,
				Length: 1,
			},
			trees:         []repository.Tree{},
			expectedError: errors.New("no trees found in estate"),
		},
		{
			name: "Invalid estate width",
			estate: repository.Estate{
				Width:  0,
				Length: 1,
			},
			trees: []repository.Tree{
				{X: 2, Y: 1, Height: 5},
				{X: 3, Y: 1, Height: 3},
				{X: 4, Y: 1, Height: 4},
			},
			expectedError: errors.New("invalid estate value"),
		},
		{
			name: "Invalid estate length",
			estate: repository.Estate{
				Width:  5,
				Length: 0,
			},
			trees: []repository.Tree{
				{X: 2, Y: 1, Height: 5},
				{X: 3, Y: 1, Height: 3},
				{X: 4, Y: 1, Height: 4},
			},
			expectedError: errors.New("invalid estate value"),
		},
		{
			name:          "Empty estate and tree list",
			estate:        repository.Estate{},
			trees:         []repository.Tree{},
			expectedError: errors.New("invalid estate value"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// we can make do with empty server as the function
			// only operate on received parameter
			srv := Server{}
			min, max, median, distance, path, err := srv.patrol(tc.estate, tc.trees)

			require.Equal(t, tc.expectedMin, min)
			require.Equal(t, tc.expectedMax, max)
			require.Equal(t, tc.expectedMedian, median)
			require.Equal(t, tc.expectedDistance, distance)
			require.Equal(t, tc.expectedPath, path)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCalculateEstateMetadata(t *testing.T) {
	tests := []struct {
		name                string
		estate              repository.Estate
		trees               []repository.Tree
		mockGetEstateErr    error
		mockGetAllTreesErr  error
		mockUpdateEstateErr error
		expectedUpdateErr   error
	}{
		{
			name: "Valid estate with trees",
			estate: repository.Estate{
				Width:  5,
				Length: 1,
			},
			trees: []repository.Tree{
				{X: 2, Y: 1, Height: 5},
				{X: 3, Y: 1, Height: 3},
				{X: 4, Y: 1, Height: 4},
			},
			mockGetEstateErr:    nil,
			mockGetAllTreesErr:  nil,
			mockUpdateEstateErr: nil,
			expectedUpdateErr:   nil,
		},
		// Add more test cases here as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new instance of the mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepositorier(ctrl)

			// Create a new server instance with the mock repository
			srv := Server{
				repository: mockRepo,
			}

			// Set up mock expectations
			mockRepo.EXPECT().GetEstateByID(gomock.Any(), gomock.Any()).Return(tc.estate, tc.mockGetEstateErr).Times(1)
			mockRepo.EXPECT().GetAllTreesInEstate(gomock.Any(), gomock.Any()).Return(tc.trees, tc.mockGetAllTreesErr).Times(1)
			mockRepo.EXPECT().UpdateEstate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockUpdateEstateErr).Times(1)

			// Call the function under test
			err := srv.calculateEstateMetadata(context.Background(), "estate_id")

			// Assert the result
			require.Equal(t, tc.expectedUpdateErr, err)
		})
	}
}
