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
func (r *JSONFileRepository) ReadByID(ctx context.Context, id string) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// 1. Lê o arquivo e trata se der erro na leitura/unmarshal
	tasks, err := r.readFromFile()
	if err != nil {
		return domain.Task{}, err
	}
	// 2. Procura a tarefa pelo ID no próprio slice
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	// 4. Caso nao encontre
	return domain.Task{}, errors.New("task not found")
}

func (r *JSONFileRepository) UpdateTaskByID(
	ctx context.Context,
	id string,
	t domain.Task,
) (domain.Task, error) {

	// 1. Adquirir bloqueio de escrita
	r.mu.Lock()
	defer r.mu.Unlock()

	// 2. Ler as tarefas diretamente do arquivo
	tasks, err := r.readFromFile()
	if err != nil {
		return domain.Task{}, err
	}

	// 3. Procurar a tarefa pelo ID
	for i := range tasks {
		if tasks[i].ID == id {

			// 4. Atualizar somente os campos desejados
			tasks[i].Title = t.Title
			tasks[i].Status = t.Status

			// 5. Persistir a lista atualizada
			if err := r.WriteAllTasks(tasks); err != nil {
				return domain.Task{}, err
			}

			return tasks[i], nil
		}
	}

	return domain.Task{}, errors.New("task not found")
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
