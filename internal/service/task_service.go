package service

import (
	"context"

	"github.com/marcosx3/tasker/internal/domain"
)

type TaskService struct {
	repo domain.Repository
}

func NewTaskService(repo domain.Repository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	// Como domain.NewTask já validou e gerou o UUID, aqui persistimos o ponteiro
	if err := s.repo.SaveTask(ctx, &task); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	return s.repo.ReadALL(ctx)
}
