package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nahwinrajan/testswpro/generated"
	"github.com/nahwinrajan/testswpro/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPostEstate(t *testing.T) {
	tests := []struct {
		name          string
		payload       generated.CreateEstateRequestBody
		mockRepoErr   error
		callRepoLayer bool
		expectedCode  int
		expectedID    string
	}{
		{
			name: "Positive Flow",
			payload: generated.CreateEstateRequestBody{
				Width:  100,
				Length: 200,
			},
			mockRepoErr:   nil,
			callRepoLayer: true,
			expectedCode:  http.StatusCreated,
			expectedID:    "mocked_estate_id",
		},
		{
			name: "Invalid Payload - Width is zero",
			payload: generated.CreateEstateRequestBody{
				Width:  0,
				Length: 200,
			},
			mockRepoErr:   nil,
			callRepoLayer: false,
			expectedCode:  http.StatusBadRequest,
			expectedID:    "",
		},
		{
			name: "Invalid Payload - Length is zero",
			payload: generated.CreateEstateRequestBody{
				Width:  100,
				Length: 0,
			},
			mockRepoErr:   nil,
			callRepoLayer: false,
			expectedCode:  http.StatusBadRequest,
			expectedID:    "",
		},
		{
			name: "Repository Error",
			payload: generated.CreateEstateRequestBody{
				Width:  100,
				Length: 200,
			},
			mockRepoErr:   errors.New("repository error"),
			callRepoLayer: true,
			expectedCode:  http.StatusBadRequest,
			expectedID:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()

			// Create a new instance of the mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepositorier(ctrl)

			// Create a new server instance with the mock repository
			srv := Server{
				repository: mockRepo,
			}

			payloadBytes, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(payloadBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			if tc.callRepoLayer {
				mockRepo.EXPECT().
					InsertEstate(gomock.Any(), tc.payload.Width, tc.payload.Length).
					Return(tc.expectedID, tc.mockRepoErr).
					Times(1)
			}

			c := e.NewContext(req, rec)
			err = srv.PostEstate(c)
			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusCreated {
				var resp generated.CreateEstateResponse
				err = json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.NotEmpty(t, resp.Id)
			}
		})
	}
}

func TestGetEstateIdStats(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		mockRepoErr     error
		callRepoLayer   bool
		expectedCode    int
		expectedStats   generated.EstateStatsResponse
		expectedMessage string
	}{
		{
			name:          "Positive Flow",
			id:            "valid_estate_id",
			mockRepoErr:   nil,
			callRepoLayer: true,
			expectedCode:  http.StatusOK,
			expectedStats: generated.EstateStatsResponse{
				Count:  10,
				Min:    5,
				Median: 15,
				Max:    20,
			},
			expectedMessage: "",
		},
		{
			name:            "Empty ID",
			id:              "",
			mockRepoErr:     nil,
			callRepoLayer:   false,
			expectedCode:    http.StatusNotFound,
			expectedStats:   generated.EstateStatsResponse{},
			expectedMessage: "resource not found",
		},
		{
			name:            "Repository Error",
			id:              "valid_estate_id",
			mockRepoErr:     errors.New("repository error"),
			callRepoLayer:   true,
			expectedCode:    http.StatusNotFound,
			expectedStats:   generated.EstateStatsResponse{},
			expectedMessage: "resource not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepositorier(ctrl)
			srv := Server{
				repository: mockRepo,
			}

			req := httptest.NewRequest(http.MethodGet, "/estate/"+tc.id+"/stats", nil)
			rec := httptest.NewRecorder()

			if tc.callRepoLayer {
				expectedEstate := repository.Estate{
					Count:  10,
					Min:    5,
					Median: 15,
					Max:    20,
				}
				mockRepo.EXPECT().
					GetEstateByID(gomock.Any(), tc.id).
					Return(expectedEstate, tc.mockRepoErr).
					Times(1)
			}

			c := e.NewContext(req, rec)
			err := srv.GetEstateIdStats(c, tc.id)
			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				var resp generated.EstateStatsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, tc.expectedStats, resp)
			} else {
				var respErr generated.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &respErr)
				require.NoError(t, err)

				require.Equal(t, tc.expectedMessage, respErr.Message)
			}
		})
	}
}

func TestGetEstateIdDronePlan(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		mockRepoErr     error
		callRepoLayer   bool
		expectedCode    int
		expectedPlan    generated.EstateDronePlanResponse
		expectedMessage string
	}{
		{
			name:          "Positive Flow",
			id:            "valid_estate_id",
			mockRepoErr:   nil,
			callRepoLayer: true,
			expectedCode:  http.StatusOK,
			expectedPlan: generated.EstateDronePlanResponse{
				Distance: 5000,
			},
			expectedMessage: "",
		},
		{
			name:            "Empty ID",
			id:              "",
			mockRepoErr:     nil,
			callRepoLayer:   false,
			expectedCode:    http.StatusNotFound,
			expectedPlan:    generated.EstateDronePlanResponse{},
			expectedMessage: "resource not found",
		},
		{
			name:            "Repository Error",
			id:              "valid_estate_id",
			mockRepoErr:     errors.New("repository error"),
			callRepoLayer:   true,
			expectedCode:    http.StatusNotFound,
			expectedPlan:    generated.EstateDronePlanResponse{},
			expectedMessage: "resource not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepositorier(ctrl)
			srv := Server{
				repository: mockRepo,
			}

			req := httptest.NewRequest(http.MethodGet, "/estate/"+tc.id+"/drone-plan", nil)
			rec := httptest.NewRecorder()

			if tc.callRepoLayer {
				expectedEstate := repository.Estate{
					PatrolDistance: 5000,
				}
				mockRepo.EXPECT().
					GetEstateByID(gomock.Any(), tc.id).
					Return(expectedEstate, tc.mockRepoErr).
					Times(1)
			}

			c := e.NewContext(req, rec)
			err := srv.GetEstateIdDronePlan(c, tc.id)
			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				var resp generated.EstateDronePlanResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, tc.expectedPlan, resp)
			} else {
				var respErr generated.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &respErr)
				require.NoError(t, err)

				require.Equal(t, tc.expectedMessage, respErr.Message)
			}
		})
	}
}

func TestPostEstateIdTree(t *testing.T) {
	validTrees := []repository.Tree{
		{X: 2, Y: 1, Height: 5},
		{X: 3, Y: 1, Height: 3},
		{X: 4, Y: 1, Height: 4},
	}

	simpleValidEstate := repository.Estate{
		ID:     "valid_estate_id",
		Width:  5,
		Length: 1,
	}

	errInvalidValueFormat := errors.New("invalid value or format")
	errResourceNotFound := errors.New("resource not found")

	tests := []struct {
		name           string
		id             string
		payload        generated.CreateTreeRequestBody
		expectedEstate repository.Estate
		// REMEMBER WE ARE MOCKING THE DB!
		estateTrees        []repository.Tree
		callRepoInsertTree bool
		callCalculate      bool
		mockInsertTreeErr  error
		mockCalculateErr   error
		mockRepoEstateErr  error
		expectedCode       int
		expectError        error
		expectedTreeID     string
	}{
		{
			name: "Positive Flow",
			id:   "valid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      1,
				Height: 4,
			},
			expectedEstate: simpleValidEstate,
			// REMEMBER WE ARE MOCKING THE DB!
			// this is what we will return from the mock,
			estateTrees:        validTrees,
			callRepoInsertTree: true,
			callCalculate:      true,
			mockInsertTreeErr:  nil,
			mockCalculateErr:   nil,
			mockRepoEstateErr:  nil,
			expectedCode:       http.StatusCreated,
			expectError:        nil,
			expectedTreeID:     "mocked_tree_id",
		},
		{
			name: "Invalid Estate ID",
			id:   "",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      1,
				Height: 4,
			},
			expectedCode: http.StatusNotFound,
			expectError:  errResourceNotFound,
		},
		{
			name: "Repository Error on GetEstateByID",
			id:   "invalid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      1,
				Height: 4,
			},
			mockRepoEstateErr: errors.New("repository error"),
			expectedCode:      http.StatusNotFound,
			expectError:       errResourceNotFound,
		},
		{
			name: "Invalid Payload - Height is Zero",
			id:   "valid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      1,
				Height: 0,
			},
			expectedEstate:    simpleValidEstate,
			mockRepoEstateErr: nil,
			expectedCode:      http.StatusBadRequest,
			expectError:       errInvalidValueFormat,
		},
		{
			name: "Invalid Payload - X Coordinate Out of Range",
			id:   "valid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      10, // Exceeds estate width
				Y:      1,
				Height: 4,
			},
			expectedEstate: simpleValidEstate,
			expectedCode:   http.StatusBadRequest,
			expectError:    errInvalidValueFormat,
		},
		{
			name: "Invalid Payload - Y Coordinate Out of Range",
			id:   "valid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      10, // Exceeds estate length
				Height: 4,
			},
			expectedEstate: simpleValidEstate,
			expectedCode:   http.StatusBadRequest,
			expectError:    errInvalidValueFormat,
		},
		{
			name: "Repository Error on InsertTree",
			id:   "valid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      1,
				Height: 4,
			},
			expectedEstate:     simpleValidEstate,
			callRepoInsertTree: true,
			mockInsertTreeErr:  errors.New("repository error"),
			expectedCode:       http.StatusBadRequest,
			expectError:        errInvalidValueFormat,
		},
		{
			name: "Repository Error on CalculateEstateMetadata",
			id:   "valid_estate_id",
			payload: generated.CreateTreeRequestBody{
				X:      4,
				Y:      1,
				Height: 4,
			},
			expectedEstate:     simpleValidEstate,
			estateTrees:        validTrees,
			callRepoInsertTree: true,
			callCalculate:      true,
			mockCalculateErr:   errors.New("faux error for mock calculate error"),
			expectedCode:       http.StatusBadRequest,
			expectError:        errInvalidValueFormat,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepositorier(ctrl)

			srv := Server{
				repository: mockRepo,
			}

			payloadBytes, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/estate/"+tc.id+"/tree", bytes.NewBuffer(payloadBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			if tc.id != "" {
				mockRepo.EXPECT().
					GetEstateByID(gomock.Any(), tc.id).
					Return(tc.expectedEstate, tc.mockRepoEstateErr).MaxTimes(2)
			}

			if tc.callRepoInsertTree {
				mockRepo.EXPECT().
					InsertTree(gomock.Any(), tc.id, tc.payload.X, tc.payload.Y, tc.payload.Height).
					Return(tc.expectedTreeID, tc.mockInsertTreeErr).
					Times(1)
			}

			if tc.callRepoInsertTree && tc.mockInsertTreeErr == nil {
				mockRepo.EXPECT().GetAllTreesInEstate(gomock.Any(), tc.id).Return(tc.estateTrees, nil)

				if tc.callCalculate && tc.mockCalculateErr == nil {
					mockRepo.EXPECT().UpdateEstate(
						gomock.Any(), // ctx
						gomock.Any(), // estate id
						gomock.Any(), // count
						gomock.Any(), // min
						gomock.Any(), // max
						gomock.Any(), // median
						gomock.Any(), // patrol distance
						gomock.Any(), // patrol route
					).Return(nil)
				} else if tc.callCalculate && tc.mockCalculateErr != nil {
					mockRepo.EXPECT().UpdateEstate(
						gomock.Any(), // ctx
						gomock.Any(), // estate id
						gomock.Any(), // count
						gomock.Any(), // min
						gomock.Any(), // max
						gomock.Any(), // median
						gomock.Any(), // patrol distance
						gomock.Any(), // patrol route
					).Return(tc.mockCalculateErr)

					mockRepo.EXPECT().
						DeleteTree(gomock.Any(), tc.expectedTreeID).
						Return(nil).
						Times(1)
				}
			}

			c := e.NewContext(req, rec)
			err = srv.PostEstateIdTree(c, tc.id)
			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusCreated {
				var resp generated.CreateTreeResponse
				err = json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.NotEmpty(t, resp.Id)
			} else {
				var errResp generated.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expectError.Error(), errResp.Message)
			}
		})
	}
}
