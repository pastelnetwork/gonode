package services

import (
	"context"
	"net/http"
)

type ErrorHandler func(context.Context, http.ResponseWriter, error)
