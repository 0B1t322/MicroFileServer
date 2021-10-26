package auth

import "github.com/go-kit/kit/endpoint"

type Auther interface {
	AuthMiddleware() 	endpoint.Middleware
	IsAdmin()			endpoint.Middleware
}