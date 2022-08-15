package middlewares

import (
	"net/http"

	"happyday/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	errorMap = map[error]common.ProblemDetails{}

	InvalidBody = common.ProblemDetails{
		Type:    "/api/errors/invalid-body",
		Title:   "INT001",
		Details: "Invalid body",
		Status:  http.StatusBadRequest,
	}

	UnexpectedError = common.ProblemDetails{
		Type:    "/api/errors/unexpected-error",
		Title:   "HYD000",
		Details: "Unexpected error",
		Status:  http.StatusInternalServerError,
	}
)

var Logger *zap.Logger

func AddErrors(errors map[error]common.ProblemDetails) {
	for k, v := range errors {
		errorMap[k] = v
	}
}

func HandleError(context *gin.Context, err error) {
	problem, found := errorMap[err]

	if !found {
		if Logger != nil {
			Logger.Error("an unexpected error happen", zap.Error(err))
		}

		problem = UnexpectedError
	}

	HandleProblem(context, problem)
}

func HandleProblem(context *gin.Context, problem common.ProblemDetails) {
	context.Header(common.ContentType, common.Problem)
	context.JSON(problem.Status, problem)
}
