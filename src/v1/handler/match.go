package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	ginmiddleware "go-test/lib/middleware/gin"
	"go-test/src/v1/contract"
)

// CreateMatchHandler godoc
//
// @Summary		Create match
// @Description	Create a new match schedule between two teams
// @Tags		matches
// @Accept		json
// @Produce		json
// @Param		body	body		contract.CreateMatchRequest	true	"create match request"
// @Success		201		{object}	ginmiddleware.Response{data=contract.MatchResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		404		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches [post]
func CreateMatchHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req contract.CreateMatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.CreateMatch(ctx, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINCreatedResponse(c, resp)
	}
}

// GetMatchHandler godoc
//
// @Summary		Get match by ID
// @Description	Get a match by its ID including goals
// @Tags		matches
// @Produce		json
// @Param		id	path		int	true	"match ID"
// @Success		200	{object}	ginmiddleware.Response{data=contract.MatchResponse}
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches/{id} [get]
func GetMatchHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.GetMatch(ctx, id)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// GetAllMatchesHandler godoc
//
// @Summary		Get all matches
// @Description	Get list of all matches
// @Tags		matches
// @Produce		json
// @Success		200	{object}	ginmiddleware.Response{data=[]contract.MatchResponse}
// @Failure		500	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches [get]
func GetAllMatchesHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		resp, err := svc.GetAllMatches(ctx)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// UpdateMatchHandler godoc
//
// @Summary		Update match
// @Description	Update a match schedule by ID (only scheduled matches can be updated)
// @Tags		matches
// @Accept		json
// @Produce		json
// @Param		id		path		int							true	"match ID"
// @Param		body	body		contract.UpdateMatchRequest	true	"update match request"
// @Success		200		{object}	ginmiddleware.Response{data=contract.MatchResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		404		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches/{id} [put]
func UpdateMatchHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		var req contract.UpdateMatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.UpdateMatch(ctx, id, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// DeleteMatchHandler godoc
//
// @Summary		Delete match
// @Description	Soft delete a match and its goals by ID
// @Tags		matches
// @Produce		json
// @Param		id	path		int	true	"match ID"
// @Success		200	{object}	ginmiddleware.Response
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches/{id} [delete]
func DeleteMatchHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		if err := svc.DeleteMatch(ctx, id); err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, nil)
	}
}

// SubmitResultHandler godoc
//
// @Summary		Submit match result
// @Description	Submit final score and goals for a match
// @Tags		matches
// @Accept		json
// @Produce		json
// @Param		id		path		int							true	"match ID"
// @Param		body	body		contract.SubmitResultRequest	true	"submit result request"
// @Success		200		{object}	ginmiddleware.Response{data=contract.MatchResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		404		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches/{id}/result [post]
func SubmitResultHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		var req contract.SubmitResultRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.SubmitResult(ctx, id, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// GetMatchReportHandler godoc
//
// @Summary		Get match report
// @Description	Get detailed match report including top scorer and team win statistics
// @Tags		matches
// @Produce		json
// @Param		id	path		int	true	"match ID"
// @Success		200	{object}	ginmiddleware.Response{data=contract.MatchReportResponse}
// @Failure		400	{object}	ginmiddleware.Response
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/matches/{id}/report [get]
func GetMatchReportHandler(svc MatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.GetMatchReport(ctx, id)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}
