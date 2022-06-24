package usecases

import (
	"db_project/app/models"
	"db_project/app/repositories"
	"db_project/utils/errors"
)

type IServiceUseCase interface {
	Clear() (err error)
	Status() (status *models.ForumStatus, err error)
}

type ServiceUseCase struct {
	serviceRepository repositories.IServiceRepository
}

func CreateServiceUseCase(serviceRepository repositories.IServiceRepository) IServiceUseCase {
	return &ServiceUseCase{serviceRepository: serviceRepository}
}

func (usecase *ServiceUseCase) Clear() (err error) {
	err = usecase.serviceRepository.Clear()
	if err != nil {
		err = errors.ServerInternal
		return
	}
	return
}

func (usecase *ServiceUseCase) Status() (status *models.ForumStatus, err error) {
	status, err = usecase.serviceRepository.Status()
	if err != nil {
		err = errors.ServerInternal
		return
	}
	return
}
