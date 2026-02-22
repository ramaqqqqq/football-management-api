package apperrors

import (
	i18n_err "go-test/lib/i18n/errors"
)

var (
	// Authentication
	ErrInvalidCredentials = i18n_err.NewI18nError("err_invalid_credentials")
	ErrEmailAlreadyExists = i18n_err.NewI18nError("err_email_already_exists")
	ErrUnauthorized       = i18n_err.NewI18nError("err_unauthorized")
	ErrInvalidToken       = i18n_err.NewI18nError("err_invalid_token")
	ErrForbidden          = i18n_err.NewI18nError("err_forbidden")

	// Validation
	ErrValidationFailed = i18n_err.NewI18nError("err_validation_failed")
	ErrInvalidRequest   = i18n_err.NewI18nError("err_invalid_request")

	// Team
	ErrTeamNotFound = i18n_err.NewI18nError("err_team_not_found")

	// Player
	ErrPlayerNotFound    = i18n_err.NewI18nError("err_player_not_found")
	ErrJerseyNumberTaken = i18n_err.NewI18nError("err_jersey_number_taken")

	// Match
	ErrMatchNotFound         = i18n_err.NewI18nError("err_match_not_found")
	ErrMatchAlreadyHasResult = i18n_err.NewI18nError("err_match_already_has_result")
	ErrMatchNotCompleted     = i18n_err.NewI18nError("err_match_not_completed")
	ErrSameTeamMatch         = i18n_err.NewI18nError("err_same_team_match")
)
