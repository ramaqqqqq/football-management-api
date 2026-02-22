package ginmiddleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-test/lib/i18n"
	i18n_err "go-test/lib/i18n/errors"
)

type Response struct {
	Data     interface{} `json:"data"`
	Error    *Error      `json:"error"`
	Success  bool        `json:"success"`
	Metadata Meta        `json:"metadata"`
}

type Meta struct {
	RequestId string `json:"request_id"`
}

type Error struct {
	Code     string `json:"code"`
	Title    string `json:"message_title"`
	Message  string `json:"message"`
	Severity string `json:"message_severity"`
	Action   *Action `json:"action"`
}

type Action struct {
	NextState string `json:"next_state"`
}

func createSuccessResponse(data interface{}, reqID string) Response {
	return Response{
		Data:    data,
		Success: true,
		Metadata: Meta{RequestId: reqID},
	}
}

func createErrorResponse(err i18n_err.I18nError, reqID, lang string) Response {
	return Response{
		Data:    nil,
		Success: false,
		Error: &Error{
			Code:     err.Error(),
			Title:    i18n.Title(lang, err.Error()),
			Message:  i18n.Message(lang, err.Error()),
			Severity: "error",
			Action:   nil,
		},
		Metadata: Meta{RequestId: reqID},
	}
}

func getLanguage(c *gin.Context) string {
	if lang := c.GetHeader("X-User-Locale"); lang != "" {
		return lang
	}
	if lang := c.GetHeader("Accept-Language"); lang != "" {
		return lang
	}
	return "en-ID"
}

func GINSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, createSuccessResponse(data, GetRequestID(c)))
}

func GINCreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, createSuccessResponse(data, GetRequestID(c)))
}

func GINBadRequestResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, createErrorResponse(i18n_err.ErrBadRequest, GetRequestID(c), getLanguage(c)))
}

func GINUnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, createErrorResponse(i18n_err.ErrUnauthorized, GetRequestID(c), getLanguage(c)))
}

func GINForbiddenResponse(c *gin.Context) {
	c.JSON(http.StatusForbidden, createErrorResponse(i18n_err.ErrForbidden, GetRequestID(c), getLanguage(c)))
}

func GINInternalErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, createErrorResponse(i18n_err.ErrInternalServer, GetRequestID(c), getLanguage(c)))
}

func GINNotFoundResponse(c *gin.Context, err i18n_err.I18nError) {
	c.JSON(http.StatusNotFound, createErrorResponse(err, GetRequestID(c), getLanguage(c)))
}

func GINErrorResponse(c *gin.Context, err error) {
	var i18nErr i18n_err.I18nError
	if errors.As(err, &i18nErr) {
		statusCode := http.StatusInternalServerError
		switch i18nErr.Error() {
		case "err_team_not_found", "err_player_not_found", "err_match_not_found",
			"err_product_not_found", "err_order_not_found", "err_user_not_found", "err_merchant_not_found":
			statusCode = http.StatusNotFound
		case "err_invalid_credentials", "err_unauthorized", "err_invalid_token":
			statusCode = http.StatusUnauthorized
		case "err_forbidden":
			statusCode = http.StatusForbidden
		case "err_bad_request", "err_validation_failed", "err_invalid_request",
			"err_insufficient_stock", "err_jersey_number_taken", "err_match_already_has_result",
			"err_match_not_completed", "err_same_team_match":
			statusCode = http.StatusBadRequest
		case "err_email_already_exists":
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, createErrorResponse(i18nErr, GetRequestID(c), getLanguage(c)))
	} else {
		GINInternalErrorResponse(c)
	}
}
