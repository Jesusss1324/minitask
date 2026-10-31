package main

import (
	"log"
	"net/http"

	"minitask/config"
	"minitask/http/handlers"
	"minitask/repository"
	"minitask/services"
	"minitask/storage"
)

func main() {
	// 1. Cargar configuración
	cfg := config.Load()

	// 2. Inicializar conexión MySQL
	db, err := storage.NewMySQL(cfg.DSN())
	if err != nil {
		log.Fatal("Error conectando con MySQL: ", err)
	}
	defer db.Close()

	// 3. Construir las dependencias
	repo := repository.NewTaskRepository(db)
	svc := services.NewTaskService(repo)
	handler := handlers.NewTaskHandler(svc)

	// 4. Configurar rutas HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/tasks", handler.Create)

	// 5. Iniciar servidor
	log.Printf("Servidor iniciado en http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
