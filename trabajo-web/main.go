package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "trabajo-web/db/sqlc"
	"trabajo-web/handlers"
	"trabajo-web/logic"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/outfits_db?sslmode=disable"
	}

	dbConn, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("error configurando la base de datos: %w", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("error conectando a la base de datos: %w", err)
	}

	queries := db.New(dbConn)

	usuariosSvc := logic.NewUsuariosService(queries)
	prendasSvc := logic.NewPrendasService(queries)
	looksSvc := logic.NewLooksService(queries)
	meGustaSvc := logic.NewMeGustaLooksService(queries)

	usuariosHdl := handlers.NewUsuariosHandler(usuariosSvc)
	prendasHdl := handlers.NewPrendasHandler(prendasSvc, usuariosSvc, looksSvc)
	looksHdl := handlers.NewLooksHandler(looksSvc, usuariosSvc, prendasSvc)
	meGustaHdl := handlers.NewMeGustaLooksHandler(meGustaSvc)

	router := handlers.SetupRoutes(usuariosHdl, prendasHdl, looksHdl, meGustaHdl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	return srv.ListenAndServe()
}
