package main

import (
	"log"
	"net/http"

	"minitask/config"
	"minitask/http/handlers"
	"minitask/http/middleware"
	"minitask/repository"
	"minitask/services"
	"minitask/storage"
)

func main() {
	cfg := config.Load()

	db, err := storage.NewMySQL(cfg.DSN())
	if err != nil {
		log.Fatal("Error conectando con MySQL: ", err)
	}
	defer db.Close()

	repo := repository.NewTaskRepository(db)
	svc := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/tasks", handler.Create)
	mux.HandleFunc("GET /api/v1/tasks", handler.List)
	mux.HandleFunc("GET /api/v1/tasks/{id}", handler.Get)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", handler.Update)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", handler.Delete)
	mux.HandleFunc("GET /api/v1/tasks/summary", handler.Summary)
	corsHandler := middleware.CORS(cfg.CORSAllowedOrigin)(mux)

	log.Printf("Servidor iniciado en http://localhost:%s", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler))
}
