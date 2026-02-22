package handler

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginmiddleware "go-test/lib/middleware/gin"
	"go-test/src/v1/contract"
)

// CreateTeamHandler godoc
//
// @Summary		Create team
// @Description	Create a new football team
// @Tags		teams
// @Accept		multipart/form-data
// @Produce		json
// @Param		name			formData	string	true	"team name"
// @Param		logo			formData	file	false	"logo image"
// @Param		year_founded	formData	int		true	"year founded (1800-2100)"
// @Param		address			formData	string	false	"team address"
// @Param		city			formData	string	true	"team city"
// @Success		201		{object}	ginmiddleware.Response{data=contract.TeamResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/teams [post]
func CreateTeamHandler(svc TeamService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		name := c.PostForm("name")
		city := c.PostForm("city")
		yearFoundedStr := c.PostForm("year_founded")
		if name == "" || city == "" || yearFoundedStr == "" {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		yearFounded, err := strconv.Atoi(yearFoundedStr)
		if err != nil || yearFounded < 1800 || yearFounded > 2100 {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		req := contract.CreateTeamRequest{
			Name:        name,
			YearFounded: yearFounded,
			Address:     c.PostForm("address"),
			City:        city,
		}

		if file, err := c.FormFile("logo"); err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			savePath := filepath.Join("uploads", "teams", uuid.New().String()+ext)
			os.MkdirAll(filepath.Join("uploads", "teams"), 0755)
			if saveErr := c.SaveUploadedFile(file, savePath); saveErr == nil {
				req.Logo = "/" + savePath
			}
		}

		resp, err := svc.CreateTeam(ctx, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINCreatedResponse(c, resp)
	}
}

// GetTeamHandler godoc
//
// @Summary		Get team by ID
// @Description	Get a football team by its ID
// @Tags		teams
// @Produce		json
// @Param		id	path		int	true	"team ID"
// @Success		200	{object}	ginmiddleware.Response{data=contract.TeamResponse}
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/teams/{id} [get]
func GetTeamHandler(svc TeamService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		resp, err := svc.GetTeam(ctx, id)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// GetAllTeamsHandler godoc
//
// @Summary		Get all teams
// @Description	Get list of all football teams
// @Tags		teams
// @Produce		json
// @Success		200	{object}	ginmiddleware.Response{data=[]contract.TeamResponse}
// @Failure		500	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/teams [get]
func GetAllTeamsHandler(svc TeamService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		resp, err := svc.GetAllTeams(ctx)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// UpdateTeamHandler godoc
//
// @Summary		Update team
// @Description	Update a football team by ID
// @Tags		teams
// @Accept		multipart/form-data
// @Produce		json
// @Param		id				path		int		true	"team ID"
// @Param		name			formData	string	false	"team name"
// @Param		logo			formData	file	false	"logo image"
// @Param		year_founded	formData	int		false	"year founded (1800-2100)"
// @Param		address			formData	string	false	"team address"
// @Param		city			formData	string	false	"team city"
// @Success		200		{object}	ginmiddleware.Response{data=contract.TeamResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		404		{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/teams/{id} [put]
func UpdateTeamHandler(svc TeamService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		req := contract.UpdateTeamRequest{
			Name:    c.PostForm("name"),
			Address: c.PostForm("address"),
			City:    c.PostForm("city"),
		}

		if yearFoundedStr := c.PostForm("year_founded"); yearFoundedStr != "" {
			yearFounded, err := strconv.Atoi(yearFoundedStr)
			if err != nil || yearFounded < 1800 || yearFounded > 2100 {
				ginmiddleware.GINBadRequestResponse(c)
				return
			}
			req.YearFounded = yearFounded
		}

		if file, err := c.FormFile("logo"); err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			savePath := filepath.Join("uploads", "teams", uuid.New().String()+ext)
			os.MkdirAll(filepath.Join("uploads", "teams"), 0755)
			if saveErr := c.SaveUploadedFile(file, savePath); saveErr == nil {
				req.Logo = "/" + savePath
			}
		}

		resp, err := svc.UpdateTeam(ctx, id, req)
		if err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, resp)
	}
}

// DeleteTeamHandler godoc
//
// @Summary		Delete team
// @Description	Soft delete a football team by ID
// @Tags		teams
// @Produce		json
// @Param		id	path		int	true	"team ID"
// @Success		200	{object}	ginmiddleware.Response
// @Failure		404	{object}	ginmiddleware.Response
// @Security	BearerAuth
// @Router		/v1/teams/{id} [delete]
func DeleteTeamHandler(svc TeamService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		if err := svc.DeleteTeam(ctx, id); err != nil {
			ginmiddleware.GINErrorResponse(c, err)
			return
		}

		ginmiddleware.GINSuccessResponse(c, nil)
	}
}
