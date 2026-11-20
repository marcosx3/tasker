package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcosx3/tasker/internal/handler"
	"github.com/marcosx3/tasker/internal/repository"
	"github.com/marcosx3/tasker/internal/service"
)

func main() {
	// 1. Inicializa o Repository (cria tasks.json na raiz se não existir)
	repo, err := repository.NewJSONFileRepository("tasks.json")
	if err != nil {
		log.Fatalf("Erro ao inicializar o repositório: %v", err)
	}

	// 2. Injeta o Repository no Service
	taskService := service.NewTaskService(repo)

	// 3. Injeta o Service no Handler
	taskHandler := handler.NewTaskHandler(taskService)

	// 4. Configura as rotas usando o ServeMux padrão do Go
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.PostTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)

	// 5. Inicia o servidor HTTP na porta 9000
	fmt.Println("Upando servidor Go!")
	fmt.Println("Servidor rodando na porta: 9000")

	if err := http.ListenAndServe(":9000", mux); err != nil {
		log.Fatalf("Erro ao subir o servidor: %v", err)
	}
}
