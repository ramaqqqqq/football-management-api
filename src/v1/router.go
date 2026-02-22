package v1

import (
	"net/http"

	"go-test/src/v1/handler"
	_ "go-test/swagger/v1"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(r *gin.Engine, deps *APIDepedencies) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", func(c *gin.Context) {
		if c.Param("any") == "/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		swaggerHandler(c)
	})

	// Auth
	auth := r.Group("/v1/auth")
	{
		auth.POST("/register", handler.RegisterHandler(deps.Services.AuthService))
		auth.POST("/login", handler.LoginHandler(deps.Services.AuthService))
	}

	// Authenticated
	authorized := r.Group("/v1", deps.JWTMiddleware.Authenticate())

	// Team
	teams := authorized.Group("/teams")
	{
		teams.GET("", handler.GetAllTeamsHandler(deps.Services.TeamService))
		teams.GET("/:id", handler.GetTeamHandler(deps.Services.TeamService))
		teams.GET("/:id/players", handler.GetPlayersByTeamHandler(deps.Services.PlayerService))
		teams.POST("", handler.CreateTeamHandler(deps.Services.TeamService))
		teams.PUT("/:id", handler.UpdateTeamHandler(deps.Services.TeamService))
		teams.DELETE("/:id", handler.DeleteTeamHandler(deps.Services.TeamService))
	}

	// Player
	players := authorized.Group("/players")
	{
		players.GET("", handler.GetAllPlayersHandler(deps.Services.PlayerService))
		players.GET("/:id", handler.GetPlayerHandler(deps.Services.PlayerService))
		players.POST("", handler.CreatePlayerHandler(deps.Services.PlayerService))
		players.PUT("/:id", handler.UpdatePlayerHandler(deps.Services.PlayerService))
		players.DELETE("/:id", handler.DeletePlayerHandler(deps.Services.PlayerService))
	}

	// Match
	matches := authorized.Group("/matches")
	{
		matches.GET("", handler.GetAllMatchesHandler(deps.Services.MatchService))
		matches.GET("/:id", handler.GetMatchHandler(deps.Services.MatchService))
		matches.GET("/:id/report", handler.GetMatchReportHandler(deps.Services.MatchService))
		matches.POST("", handler.CreateMatchHandler(deps.Services.MatchService))
		matches.PUT("/:id", handler.UpdateMatchHandler(deps.Services.MatchService))
		matches.DELETE("/:id", handler.DeleteMatchHandler(deps.Services.MatchService))
		matches.POST("/:id/result", handler.SubmitResultHandler(deps.Services.MatchService))
	}
}
