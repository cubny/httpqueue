//go:build mocks

package mocks

import (
	// ensures this package is in vendors folder and fixes a bug in go:generate that appears because of use of reflection in mocks generation
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./app/timer/repo_mock.go -package=mocks -mock_names=Repo=Repo github.com/cubny/httpqueue/internal/app/timer Repo
//go:generate mockgen -destination=./app/timer/service_mock.go -package=mocks -mock_names=Service=Service github.com/cubny/httpqueue/internal/app/timer Service
//go:generate mockgen -destination=./app/timer/outbox_mock.go -package=mocks -mock_names=Outbox=Outbox github.com/cubny/httpqueue/internal/app/timer Outbox
//go:generate mockgen -destination=./app/timer/producer_mock.go -package=mocks -mock_names=Producer=Producer github.com/cubny/httpqueue/internal/app/timer Producer
//go:generate mockgen -destination=./app/timer/http_client_mock.go -package=mocks -mock_names=HttpClient=HttpClient github.com/cubny/httpqueue/internal/app/timer HttpClient

//region external
//go:generate mockgen -destination=./external/asynq/broker_mock.go -package=mocks -mock_names=Broker=Broker github.com/cubny/httpqueue/internal/infra/asynq/timer Broker
//go:generate mockgen -destination=./external/asynq/redisconnopt_mock.go -package=mocks -mock_names=RedisConnOpt=RedisConnOpt github.com/hibiken/asynq RedisConnOpt
//go:generate mockgen -destination=./external/asynq/handler_mock.go -package=mocks -mock_names=Handler=Handler github.com/hibiken/asynq Handler
//go:generate mockgen -destination=./external/redis/client_mock.go -package=mocks -mock_names=UniversalClient=RedisClient github.com/go-redis/redis/v8 UniversalClient
//go:generate mockgen -destination=./external/redis/pipeliner_mock.go -package=mocks -mock_names=Pipeliner=RedisPipeliner github.com/go-redis/redis/v8 Pipeliner
//endregion
