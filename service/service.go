package service

import (
	"context"
	"log"

	constants "github.com/tejiriaustin/lema/constants"
	"github.com/tejiriaustin/lema/env"
	"github.com/tejiriaustin/lema/logger"
)

type (
	Container struct {
		UserService UserServiceInterface
		PostService PostServiceInterface
	}

	Pager struct {
		Page    int64
		PerPage int64
	}
)

func NewService(lemaLogger logger.Logger, conf *env.Environment) *Container {
	log.Println("Creating Service Container...")
	return &Container{
		UserService: NewUserService(lemaLogger),
		PostService: NewPostService(lemaLogger),
	}
}

func GetPageNumberFromContext(ctx context.Context) int64 {
	n, ok := ctx.Value(constants.ContextKeyPageNumber).(int64)
	if !ok {
		return 0
	}
	return n
}

func GetPerPageLimitFromContext(ctx context.Context) int64 {
	l, ok := ctx.Value(constants.ContextKeyPerPageLimit).(int64)
	if !ok {
		return 0
	}
	return l
}
