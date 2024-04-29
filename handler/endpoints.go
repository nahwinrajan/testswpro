package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nahwinrajan/testswpro/generated"
)

func (srv *Server) PostEstate(ectx echo.Context) error {
	var payload generated.CreateEstateRequestBody
	var respBadReq generated.ErrorResponse
	respBadReq.Message = "invalid value or format"

	defer ectx.Request().Body.Close()
	err := ectx.Bind(&payload)
	if err != nil {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateEstate] failed read payload, err:%s", err)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	// validation
	switch {
	case payload.Width < 1 || payload.Width > 50000:
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateEstate] invalid value width:%+v, payload:%+v", payload.Width, payload)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	case payload.Length < 1 || payload.Length > 50000:
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateEstate] invalid value length:%+v", payload.Length)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	id, err := srv.Repository.InsertEstate(ectx.Request().Context(), payload.Width, payload.Length)
	if err != nil {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateEstate] failed to insert payload:%+v, err:%s", payload, err)
		// the problem statement does not specify this error, nor in api docs..
		// but it is possibility, should have been InternalServer error but..
		respBadReq.Message = "failed to create resource"
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	var resp generated.CreateEstateResponse
	resp.Id = id
	return ectx.JSON(http.StatusCreated, resp)
}

func (srv *Server) PostEstateIdTree(ectx echo.Context, id int) error {
	var payload generated.CreateTreeRequestBody
	var respBadReq generated.ErrorResponse
	respBadReq.Message = "invalid value or format"

	strEstateID := ectx.Param("id")
	if len(strEstateID) == 0 {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateTree] param tree_id not passed")
		respBadReq.Message = "resource not found"
		return ectx.JSON(http.StatusNotFound, respBadReq)
	}

	estate, err := srv.Repository.GetEstateByID(
		ectx.Request().Context(),
		strEstateID,
	)
	if err != nil {
		respBadReq.Message = "resource not found"
		return ectx.JSON(http.StatusNotFound, respBadReq)
	}

	defer ectx.Request().Body.Close()
	err = ectx.Bind(&payload)
	if err != nil {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateTree] failed to read payload, err:%s", err)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	switch {
	case payload.Height < 1 || payload.Height > 30:
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateTree] invalid value height:%+v", payload.Height)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	case payload.X < 1 || payload.X > estate.Width:
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateTree] invalid value x:%+v", payload.X)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	case payload.Y < 1 || payload.Y > estate.Length:
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateTree] invalid value y:%+v", payload.Y)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	strTreeID, err := srv.Repository.InsertTree(
		ectx.Request().Context(),
		strEstateID,
		payload.X,
		payload.Y,
		payload.Height,
	)
	// TODO: cater for Internal Server error
	if err != nil {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[CreateTree] failed to insert payload:%+v, err:%s", payload, err)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	var resp generated.CreateTreeResponse
	resp.Id = strTreeID

	// TODO: calculate estate metadata (i.e: count, max, min, median, distance, routes)
	err = srv.patrol(ectx.Request().Context(), strEstateID)
	if err != nil {
		// TODO: what should we do ?
		// ideally, we have this on background, with multiple retry and then
		// if still fail output alert to slack, manual rectify issue, and re-trigger calculation
		ectx.Logger().Errorf("[CreateTree] failed to calculate stats and distance, err:%s", err)
		srv.Repository.DeleteTree(ectx.Request().Context(), strTreeID)
		return ectx.JSON(http.StatusBadRequest, respBadReq)
	}

	return ectx.JSON(http.StatusCreated, resp)
}

func (srv *Server) GetEstateIdStats(ectx echo.Context, id int) error {
	var respBadReq generated.ErrorResponse
	respBadReq.Message = "invalid value or format"

	strEstateID := ectx.Param("id")
	if len(strEstateID) == 0 {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[GetEstateStats] param estate_id not passed")
		respBadReq.Message = "resource not found"
		return ectx.JSON(http.StatusNotFound, respBadReq)
	}

	estate, err := srv.Repository.GetEstateByID(
		ectx.Request().Context(),
		strEstateID,
	)
	if err != nil {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[GetEstateStats] failed to read estate_id:%+v, err:%s", strEstateID, err)
		respBadReq.Message = "resource not found"
		return ectx.JSON(http.StatusNotFound, respBadReq)
	}

	var resp generated.EstateStatsResponse
	resp.Count = estate.Count
	resp.Min = estate.Min
	resp.Median = estate.Median
	resp.Max = estate.Max

	return ectx.JSON(http.StatusOK, resp)
}

func (srv *Server) GetEstateIdDronePlan(ectx echo.Context, id int) error {
	var respBadReq generated.ErrorResponse
	respBadReq.Message = "invalid value or format"

	strEstateID := ectx.Param("id")
	if len(strEstateID) == 0 {
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[DronePlan] param estate_id not passed")
		respBadReq.Message = "resource not found"
		return ectx.JSON(http.StatusNotFound, respBadReq)
	}

	estate, err := srv.Repository.GetEstateByID(
		ectx.Request().Context(),
		strEstateID,
	)
	if err != nil {
		respBadReq.Message = "resource not found"
		// TODO: change log level according log level company guideline (info, error, etc)
		ectx.Logger().Errorf("[DronePlan] failed retrieve information estate_id:%s, err:%s", strEstateID, err)
		return ectx.JSON(http.StatusNotFound, respBadReq)
	}

	var resp generated.EstateDronePlanResponse
	resp.Distance = estate.PatrolDistance

	// // if the request without query param
	// strMaxDistance := ectx.QueryParam("max_distance")
	// if len(strMaxDistance) == 0 {
	// 	// TODO: change log level according log level company guideline (info, error, etc)
	// 	ectx.Logger().Errorf("[DronePlan] max_distance not passed")
	// 	return ectx.JSON(http.StatusOK, resp)
	// }

	// maxDistance, err := strconv.Atoi(strMaxDistance)
	// if err != nil {
	// 	// TODO: change log level according log level company guideline (info, error, etc)
	// 	ectx.Logger().Errorf("[DronePlan] invalid param max_distance: %s, err:%s", strMaxDistance, err)
	// 	return ectx.JSON(http.StatusBadRequest, respBadReq)
	// }

	// resp.Rest.X, resp.Rest.Y, err = srv.calculateMaxDistance(estate, maxDistance)
	// if err != nil {
	// 	// TODO: should be Internal Server Error, but its not on problem spec
	// 	// TODO: change log level according log level company guideline (info, error, etc)
	// 	ectx.Logger().Errorf("[calculateMaxDistance] return error: %s", err)
	// 	return ectx.JSON(http.StatusBadRequest, respBadReq)
	// }

	return ectx.JSON(http.StatusOK, resp)
}

// NOTE: this version return the (x,y) where (x,y) is the closest distance to max distance OR exact match to maxDistance
// unfortunately there are edge cases where the difference can be right at the border (under 10m), plot(x,y) and the next move is vertical for tree height,
// func (srv *Server) calculateMaxDistance(estate repository.Estate, maxDistance int) (restX, restY int, err error) {
// 	if maxDistance >= estate.PatrolDistance {
// 		restX = estate.Width
// 		restY = estate.Length
// 		return
// 	}

// 	strRoutes := strings.Split(estate.PatrolRoute, ";")
// 	lenRoutes := len(strRoutes)
// 	// postX, postY, currDistance := make([]int, 0, lenRoutes), make([]int, 0, lenRoutes), make([]int, 0, lenRoutes)
// 	// directions := make([]string, 0, lenRoutes)
// 	// var idx int
// 	var x, y, distance int
// 	for i := 0; i < lenRoutes; i++ {
// 		// stepStrFormat: step_#,x,y,direction,step_distance,current_distance
// 		strs := strings.Split(strRoutes[i], ",")

// 		x, err = strconv.Atoi(strs[1])
// 		if err != nil {
// 			return
// 		}

// 		y, err = strconv.Atoi(strs[2])
// 		if err != nil {
// 			return
// 		}

// 		distance, err = strconv.Atoi(strs[5])
// 		if err != nil {
// 			return
// 		}

// 		// early exit if we found the target
// 		switch {
// 		case distance == maxDistance:
// 			return x, y, nil
// 		case distance <= maxDistance:
// 			// otherwise keep tracking the closest distance traveled to max distance
// 			restX = x
// 			restY = y
// 			// idx = i
// 		}

// 		// postX = append(postX, x)
// 		// postY = append(postY, y)
// 		// currDistance = append(currDistance, distance)
// 		// directions = append(directions, strs[3])
// 	}

// 	// TODO find the last mile: detect if the next step is vertical or horizontal
// 	// we've gon through the routes and no exact match
// 	// find exactly where the drone ran out of distance
// 	// stepStrFormat: step_#,x,y,direction,step_distance,current_distance
// 	// stepStrFormat = "%d,%d,%d,%s,%d,%d"

// 	// we should have update restX, restY to closest position without it running out of battery
// 	// there are edge cases next step distance more than allocated remaining max_distance (i.e:
// 	// climb 23m to monitor yet distance left 3m; move to next plot which is 10m but remaining distance 1m)
// 	return
// }
