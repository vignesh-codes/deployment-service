package svc

import (
	adapter "deployment-service/apps/repository/adapter"
	"deployment-service/constants"
	"deployment-service/logger"
	model_event_logger "deployment-service/models/model.eventlogger"

	"go.uber.org/zap"
)

type EventLoggerService struct {
	repository *adapter.Repository
}

func (svc EventLoggerService) LogEvent(payload model_event_logger.Event) error {
	logger.EventLogger.Info(constants.SERVICE_NAME, zap.Any("payload", payload))
	svc.repository.Kubernetes.ListDeployments("default")
	return nil
}
