package handler

import (
	"shopsvc/internal/core/port"
	"shopsvc/pkg/validator"
)

// handler is the user handler
type handler struct {
	service   port.Service
	validator validator.Validator
}

// New creates a new user handler
func New(
	service port.Service,
	validator validator.Validator,
) port.Handler {
	return &handler{
		service:   service,
		validator: validator,
	}
}
