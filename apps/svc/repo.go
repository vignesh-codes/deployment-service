package svc

import (
	adapter "deployment-service/apps/repository/adapter"
)

var ServiceRepo *ServiceRepository

type ServiceRepository struct {
	EventLoggerService *EventLoggerService
}

func NewServiceRepo(repository *adapter.Repository) *ServiceRepository {
	return &ServiceRepository{
		EventLoggerService: &EventLoggerService{repository},
	}
}
