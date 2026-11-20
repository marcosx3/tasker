package dto

import (
	"github.com/marcosx3/tasker/internal/domain"
)

type TaskRequest struct {
	Title  string `json:"title" validate:"required,min:2"`
	Status string `json:"status" validate:"required,min:2"`
}

type TaskResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func (request TaskRequest) ToDomain() (domain.Task, error) {
	return domain.NewTask(request.Title, request.Status)
}

func ToTaskResponse(task domain.Task) TaskResponse {
	return TaskResponse{
		ID:     task.ID,
		Title:  task.Title,
		Status: task.Status,
	}
}
