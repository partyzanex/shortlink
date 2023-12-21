package service

import "github.com/pkg/errors"

var (
	ErrServiceIsStarted = errors.New("service is started")
	ErrInitialized      = errors.New("service is initialized")
	ErrInitialize       = errors.New("service cannot initialize")
)
