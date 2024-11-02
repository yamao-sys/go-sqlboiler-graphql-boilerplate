package graph

import "app/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authService services.AuthService
	todoService services.TodoService
}

func NewResolver(authService services.AuthService, todoService services.TodoService) *Resolver {
	return &Resolver{
		authService: authService,
		todoService: todoService,
	}
}
