package domain

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Task struct {
	ID     string
	Title  string
	Status string
}

func NewTask(title, status string) (Task, error) {
	title = strings.TrimSpace(title)
	status = strings.TrimSpace(strings.ToLower(status))
	if len(title) < 2 {
		return Task{}, errors.New("title must be has more than 2 caracteres")
	}
	if len(status) < 2 {
		return Task{}, errors.New("title status be has more than 2 caracteres")
	}
	return Task{
		ID:     uuid.NewString(),
		Title:  title,
		Status: status,
	}, nil
}

type Repository interface {
	SaveTask(ctx context.Context, t *Task) error
	ReadALL(ctx context.Context) ([]Task, error)
	ReadByID(ctx context.Context, id string) (Task, error)
	UpdateTaskByID(ctx context.Context, id string, t Task) (Task, error)
}
