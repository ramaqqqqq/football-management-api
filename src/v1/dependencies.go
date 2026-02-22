package v1

import (
	"context"

	"go-test/lib/atomic"
	atomicSQLX "go-test/lib/atomic/sqlx"
	ginmiddleware "go-test/lib/middleware/gin"
	"go-test/lib/provider"
	"go-test/src/app"
	goalRepo "go-test/src/repository/goal"
	matchRepo "go-test/src/repository/match"
	playerRepo "go-test/src/repository/player"
	teamRepo "go-test/src/repository/team"
	userRepo "go-test/src/repository/user"
	"go-test/src/v1/service"

	"github.com/sirupsen/logrus"
)

type APIRepositories struct {
	AtomicSessionProvider atomic.AtomicSessionProvider
	UserRepo              *userRepo.UserRepository
	TeamRepo              *teamRepo.TeamRepository
	PlayerRepo            *playerRepo.PlayerRepository
	MatchRepo             *matchRepo.MatchRepository
	GoalRepo              *goalRepo.GoalRepository
}

type APIServices struct {
	AuthService   *service.AuthService
	TeamService   *service.TeamService
	PlayerService *service.PlayerService
	MatchService  *service.MatchService
}

type APIDepedencies struct {
	Repositories  *APIRepositories
	Services      *APIServices
	JWTMiddleware *ginmiddleware.GinJWTMiddleware
}

func instantiateAPIRepositories(ctx context.Context) *APIRepositories {
	var r APIRepositories
	var err error

	r.AtomicSessionProvider = atomicSQLX.NewSqlxAtomicSessionProvider(app.DB())

	r.UserRepo, err = userRepo.InitUserRepository(ctx, app.DB())
	if err != nil {
		logrus.WithContext(ctx).Fatal("init user repo err: ", err)
	}

	r.TeamRepo, err = teamRepo.InitTeamRepository(ctx, app.DB())
	if err != nil {
		logrus.WithContext(ctx).Fatal("init team repo err: ", err)
	}

	r.PlayerRepo, err = playerRepo.InitPlayerRepository(ctx, app.DB())
	if err != nil {
		logrus.WithContext(ctx).Fatal("init player repo err: ", err)
	}

	r.MatchRepo, err = matchRepo.InitMatchRepository(ctx, app.DB())
	if err != nil {
		logrus.WithContext(ctx).Fatal("init match repo err: ", err)
	}

	r.GoalRepo, err = goalRepo.InitGoalRepository(ctx, app.DB())
	if err != nil {
		logrus.WithContext(ctx).Fatal("init goal repo err: ", err)
	}

	return &r
}

func instantiateAPIServices(ctx context.Context, r *APIRepositories) *APIServices {
	cfg := app.Config().JWT

	privateKey, err := provider.LoadRSAPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		logrus.WithContext(ctx).Fatalf("load RSA private key: %v", err)
	}

	pswdProvider := &provider.Bcrypt{}
	pswdComparator := &provider.Bcrypt{}

	return &APIServices{
		AuthService: service.NewAuthService(
			r.UserRepo,
			r.AtomicSessionProvider,
			privateKey,
			pswdProvider,
			pswdComparator,
		),
		TeamService: service.NewTeamService(
			r.TeamRepo,
			r.AtomicSessionProvider,
		),
		PlayerService: service.NewPlayerService(
			r.PlayerRepo,
			r.TeamRepo,
			r.AtomicSessionProvider,
		),
		MatchService: service.NewMatchService(
			r.MatchRepo,
			r.TeamRepo,
			r.PlayerRepo,
			r.GoalRepo,
			r.AtomicSessionProvider,
		),
	}
}

func Dependencies(ctx context.Context) *APIDepedencies {
	cfg := app.Config().JWT

	publicKey, err := provider.LoadRSAPublicKey(cfg.PublicKeyPath)
	if err != nil {
		logrus.WithContext(ctx).Fatalf("load RSA public key: %v", err)
	}

	repositories := instantiateAPIRepositories(ctx)
	services := instantiateAPIServices(ctx, repositories)
	jwtMiddleware := ginmiddleware.NewGinJWTMiddleware(publicKey)

	return &APIDepedencies{
		Repositories:  repositories,
		Services:      services,
		JWTMiddleware: jwtMiddleware,
	}
}
