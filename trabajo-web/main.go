package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "trabajo-web/db/sqlc"
	"trabajo-web/handlers"
	"trabajo-web/logic"
)

func main() {
	// obtener URL de conexión
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/outfits_db?sslmode=disable"
	}

	// conectar a la base de datos
	dbConn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Error configurando la base de datos: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	// capa DB (sqlc)
	queries := db.New(dbConn)

	// capa Lógica (services)
	usuariosSvc := logic.NewUsuariosService(queries)
	prendasSvc := logic.NewPrendasService(queries)
	looksSvc := logic.NewLooksService(queries)
	meGustaSvc := logic.NewMeGustaLooksService(queries)

	// capa HTTP (handlers)
	usuariosHdl := handlers.NewUsuariosHandler(usuariosSvc)
	prendasHdl := handlers.NewPrendasHandler(prendasSvc)
	looksHdl := handlers.NewLooksHandler(looksSvc)
	meGustaHdl := handlers.NewMeGustaLooksHandler(meGustaSvc)

	// configurar rutas
	router := handlers.SetupRoutes(usuariosHdl, prendasHdl, looksHdl, meGustaHdl)

	// iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
