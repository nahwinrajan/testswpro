// This file will run automated tests for API.
package tests

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// )

// const ApiUrl = "http://localhost:8080"

// func TestApi(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skip API tests")
// 	}

// 	testcases := getTestCases()
// 	ctx := context.Background()
// 	client := &http.Client{}

// 	for _, tc := range testcases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			for idx := range tc.Steps {
// 				step := &tc.Steps[idx]
// 				request, err := step.Request(t, ctx, &tc)
// 				request.Header.Set("Content-Type", "application/json")
// 				request.Header.Set("Accept", "application/json")
// 				require.NoError(t, err)

// 				// Send request
// 				response, err := client.Do(request)

// 				require.NoError(t, err)
// 				defer response.Body.Close()

// 				// Check response
// 				ReadJsonResult(t, response, step)
// 				step.Expect(t, ctx, &tc, response, step.Result)
// 			}
// 		})
// 	}
// }

// func getTestCases() []TestCase {
// 	return []TestCase{
// 		{
// 			Name: "Test Hello",
// 			Steps: []TestCaseStep{
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						return http.NewRequest("GET", ApiUrl+"/hello", nil)
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name: "Test Hello with name",
// 			Steps: []TestCaseStep{
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						return http.NewRequest("GET", ApiUrl+"/hello?id=123", nil)
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						require.Equal(t, "Hello User 123", data["message"])
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						return http.NewRequest("GET", ApiUrl+"/hello?id=456", nil)
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						step1 := tc.Steps[0]
// 						require.Equal(t, "Hello User 123", step1.Result["message"])
// 					},
// 				},
// 			},
// 		},
// 		//----- Test for API
// 		{
// 			Name: "Test 1",
// 			Steps: []TestCaseStep{
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						return http.NewRequest("POST", ApiUrl+"/estate", nil)
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name: "Test 2",
// 			Steps: []TestCaseStep{
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"length": 10,
// 							"width":  20,
// 						}
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name: "Test 3",
// 			Steps: []TestCaseStep{
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"length": 10,
// 							"width":  20,
// 						}
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						fmt.Println("Data:", data)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"height": 10,
// 							"x":      5,
// 							"y":      5,
// 						}
// 						id := tc.Steps[0].Result["id"].(string)
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate/"+id+"/tree", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"height": 20,
// 							"x":      6,
// 							"y":      5,
// 						}
// 						id := tc.Steps[0].Result["id"].(string)
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate/"+id+"/tree", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name: "Test 4",
// 			Steps: []TestCaseStep{
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"length": 5,
// 							"width":  1,
// 						}
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"height": 10,
// 							"x":      2,
// 							"y":      1,
// 						}
// 						id := tc.Steps[0].Result["id"].(string)
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate/"+id+"/tree", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"height": 20,
// 							"x":      3,
// 							"y":      1,
// 						}
// 						id := tc.Steps[0].Result["id"].(string)
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate/"+id+"/tree", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						req := map[string]int{
// 							"height": 10,
// 							"x":      4,
// 							"y":      1,
// 						}
// 						id := tc.Steps[0].Result["id"].(string)
// 						body, err := json.Marshal(req)
// 						require.NoError(t, err)
// 						return http.NewRequest("POST", ApiUrl+"/estate/"+id+"/tree", bytes.NewReader(body))
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						RequireIsUUID(t, data["id"].(string))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						id := tc.Steps[0].Result["id"].(string)
// 						return http.NewRequest("GET", ApiUrl+"/estate/"+id+"/stats", nil)
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						fmt.Println("Data:", data)
// 						require.Equal(t, 3, int(data["count"].(float64)))
// 						require.Equal(t, 10, int(data["min"].(float64)))
// 						require.Equal(t, 20, int(data["max"].(float64)))
// 						require.Equal(t, 10, int(data["median"].(float64)))
// 					},
// 				},
// 				{
// 					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
// 						id := tc.Steps[0].Result["id"].(string)
// 						return http.NewRequest("GET", ApiUrl+"/estate/"+id+"/drone-plan", nil)
// 					},
// 					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
// 						require.Equal(t, http.StatusOK, resp.StatusCode)
// 						require.Equal(t, 82, int(data["distance"].(float64)))
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// type TestCase struct {
// 	Name  string
// 	Steps []TestCaseStep
// }

// type RequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
// type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

// type TestCaseStep struct {
// 	Request RequestFunc
// 	Expect  ExpectFunc
// 	Result  map[string]any
// }

// func ResponseContains(t *testing.T, resp *http.Response, text string) {
// 	body, err := io.ReadAll(resp.Body)
// 	bodyStr := string(body)
// 	require.NoError(t, err)
// 	require.Contains(t, bodyStr, text)
// }

// func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
// 	var result map[string]any
// 	err := json.NewDecoder(resp.Body).Decode(&result)
// 	step.Result = result
// 	require.NoError(t, err)
// }

// func RequireIsUUID(t *testing.T, value string) {
// 	_, err := uuid.Parse(value)
// 	require.NoError(t, err)
// }
