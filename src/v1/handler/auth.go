package handler

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "go-test/lib/middleware/gin"
	apperrors "go-test/src/errors"
	"go-test/src/v1/contract"
)

// RegisterHandler godoc
//
// @Summary		Register user
// @Description	Register a new user account
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		body	body		contract.RegisterRequest	true	"register request"
// @Success		201		{object}	ginmiddleware.Response{data=contract.AuthResponse}
// @Failure		400		{object}	ginmiddleware.Response
// @Failure		409		{object}	ginmiddleware.Response
// @Router		/v1/auth/register [post]
func RegisterHandler(svc AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req contract.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		if err := contract.ValidatePassword(req.Password); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		authResp, err := svc.Register(ctx, req)
		if err != nil {
			switch err {
			case apperrors.ErrEmailAlreadyExists:
				ginmiddleware.GINErrorResponse(c, err)
			default:
				ginmiddleware.GINBadRequestResponse(c)
			}
			return
		}

		ginmiddleware.GINCreatedResponse(c, authResp)
	}
}

// LoginHandler godoc
//
// @Summary		Login user
// @Description	Login with email and password
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		body	body		contract.LoginRequest	true	"login request"
// @Success		200		{object}	ginmiddleware.Response{data=contract.AuthResponse}
// @Failure		401		{object}	ginmiddleware.Response
// @Router		/v1/auth/login [post]
func LoginHandler(svc AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req contract.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			ginmiddleware.GINBadRequestResponse(c)
			return
		}

		authResp, err := svc.Login(ctx, req)
		if err != nil {
			ginmiddleware.GINUnauthorizedResponse(c)
			return
		}

		ginmiddleware.GINSuccessResponse(c, authResp)
	}
}
