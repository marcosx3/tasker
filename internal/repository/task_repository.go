package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/marcosx3/tasker/internal/domain"
)

type JSONFileRepository struct {
	mu       sync.RWMutex
	filePath string
}

func NewJSONFileRepository(filepath string) (*JSONFileRepository, error) {
	repo := &JSONFileRepository{filePath: filepath}
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		if err := repo.WriteAllTasks([]domain.Task{}); err != nil {
			return nil, err
		}
	}
	return repo, nil
}

func (r *JSONFileRepository) WriteAllTasks(tasks []domain.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := r.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpFile, r.filePath)
}

func (r *JSONFileRepository) SaveTask(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readFromFile()
	if err != nil {
		return err
	}

	tasks = append(tasks, *task)
	return r.WriteAllTasks(tasks)
}

// ReadALL cumpre o contrato de domain.Repository com trava de leitura (RLock)
func (r *JSONFileRepository) ReadALL(ctx context.Context) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.readFromFile()
}

// Método interno sem lock para ser reutilizado tanto por SaveTask quanto por ReadALL
func (r *JSONFileRepository) readFromFile() ([]domain.Task, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []domain.Task{}, nil
	}

	var tasks []domain.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
