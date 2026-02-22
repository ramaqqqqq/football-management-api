package response

import (
	"context"
	"errors"
	"net/http"

	i18n_err "go-test/lib/i18n/errors"
	"go-test/lib/middleware/request"
)

func JSONSuccessResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	JSONResponse(ctx, w, createSuccessResponse(data, request.GetRequestID(ctx)), http.StatusOK)
}

func JSONCreatedResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	JSONResponse(ctx, w, createSuccessResponse(data, request.GetRequestID(ctx)), http.StatusCreated)
}

func JSONUnauthorizedResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrUnauthorized, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), http.StatusUnauthorized)
}

func JSONInternalErrorResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrInternalServer, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), http.StatusInternalServerError)
}

func JSONBadRequestResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrBadRequest, request.GetRequestID(ctx), request.GetLanguage(ctx), nil),
		http.StatusBadRequest)
}

func JSONUnprocessableEntity(ctx context.Context, w http.ResponseWriter, err i18n_err.I18nError, action *Action) {
	JSONResponse(ctx, w, createErrorResponse(err, request.GetRequestID(ctx), request.GetLanguage(ctx), action), http.StatusUnprocessableEntity)
}

func JSONNotFoundResponse(ctx context.Context, w http.ResponseWriter, err i18n_err.I18nError) {
	JSONResponse(ctx, w, createErrorResponse(err, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), http.StatusNotFound)
}

func JSONForbiddenResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrForbidden, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), http.StatusForbidden)
}

func JSONErrorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	var i18nErr i18n_err.I18nError
	if errors.As(err, &i18nErr) {
		statusCode := http.StatusInternalServerError

		switch i18nErr.Error() {
		case "err_product_not_found", "err_order_not_found", "err_user_not_found", "err_merchant_not_found":
			statusCode = http.StatusNotFound
		case "err_invalid_credentials", "err_unauthorized", "err_invalid_token":
			statusCode = http.StatusUnauthorized
		case "err_forbidden":
			statusCode = http.StatusForbidden
		case "err_bad_request", "err_validation_failed", "err_invalid_request", "err_insufficient_stock":
			statusCode = http.StatusBadRequest
		case "err_email_already_exists":
			statusCode = http.StatusConflict
		}

		JSONResponse(ctx, w, createErrorResponse(i18nErr, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), statusCode)
	} else {
		JSONInternalErrorResponse(ctx, w)
	}
}
