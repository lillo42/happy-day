package common

type (
	ProblemDetails struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Details string `json:"message"`
		Status  int    `json:"statusCode"`
	}
)
