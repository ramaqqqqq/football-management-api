package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	ginmiddleware "go-test/lib/middleware/gin"
	"go-test/src/v1/contract"
)

// CreatePlayerHandler godoc
//
// @Summary		Create player
// @Description	Create a new player and assign to a team
// @Tags		players
// @Accept		json
// @Produce		json
// @Param		body	body		contract.CreatePlayerRequest	true	"create player request"
// @Success		201		{object}	ginmiddleware.Response{data=contract.PlayerResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		404		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/players [post]
func CreatePlayerHandler(svc PlayerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req contract.CreatePlayerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.CreatePlayer(ctx, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINCreatedResponse(c, resp)
	}
}

// GetPlayerHandler godoc
//
// @Summary		Get player by ID
// @Description	Get a player by their ID
// @Tags		players
// @Produce		json
// @Param		id	path		int	true	"player ID"
// @Success		200	{object}	ginmiddleware.Response{data=contract.PlayerResponse}
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/players/{id} [get]
func GetPlayerHandler(svc PlayerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.GetPlayer(ctx, id)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// GetAllPlayersHandler godoc
//
// @Summary		Get all players
// @Description	Get list of all players
// @Tags		players
// @Produce		json
// @Success		200	{object}	ginmiddleware.Response{data=[]contract.PlayerResponse}
// @Failure		500	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/players [get]
func GetAllPlayersHandler(svc PlayerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		resp, err := svc.GetAllPlayers(ctx)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// GetPlayersByTeamHandler godoc
//
// @Summary		Get players by team
// @Description	Get all players belonging to a specific team
// @Tags		teams
// @Produce		json
// @Param		id	path		int	true	"team ID"
// @Success		200	{object}	ginmiddleware.Response{data=[]contract.PlayerResponse}
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/teams/{id}/players [get]
func GetPlayersByTeamHandler(svc PlayerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.GetPlayersByTeam(ctx, teamID)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// UpdatePlayerHandler godoc
//
// @Summary		Update player
// @Description	Update a player by ID
// @Tags		players
// @Accept		json
// @Produce		json
// @Param		id		path		int								true	"player ID"
// @Param		body	body		contract.UpdatePlayerRequest	true	"update player request"
// @Success		200		{object}	ginmiddleware.Response{data=contract.PlayerResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		404		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/players/{id} [put]
func UpdatePlayerHandler(svc PlayerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		var req contract.UpdatePlayerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.UpdatePlayer(ctx, id, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// DeletePlayerHandler godoc
//
// @Summary		Delete player
// @Description	Soft delete a player by ID
// @Tags		players
// @Produce		json
// @Param		id	path		int	true	"player ID"
// @Success		200	{object}	ginmiddleware.Response
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/players/{id} [delete]
func DeletePlayerHandler(svc PlayerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		if err := svc.DeletePlayer(ctx, id); err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, nil)
	}
}
