package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "trabajo-web/db/sqlc"
	"trabajo-web/handlers"
	"trabajo-web/logic"
)

func setupTestServer() http.Handler {
	dbURL := "postgres://postgres:mysecretpassword@localhost:5432/outfits_db?sslmode=disable"

	dbConn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Error configurando la base de datos: %v", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	queries := db.New(dbConn)

	usuariosSvc := logic.NewUsuariosService(queries)
	prendasSvc := logic.NewPrendasService(queries)
	looksSvc := logic.NewLooksService(queries)
	meGustaSvc := logic.NewMeGustaLooksService(queries)

	usuariosHdl := handlers.NewUsuariosHandler(usuariosSvc)
	prendasHdl := handlers.NewPrendasHandler(prendasSvc)
	looksHdl := handlers.NewLooksHandler(looksSvc)
	meGustaHdl := handlers.NewMeGustaLooksHandler(meGustaSvc)

	return handlers.SetupRoutes(usuariosHdl, prendasHdl, looksHdl, meGustaHdl)
}

func TestAPIIntegrationCompleto(t *testing.T) {
	ts := httptest.NewServer(setupTestServer())
	defer ts.Close()

	client := ts.Client()

	var userID int
	var prendaID int
	var lookID int

	// -------------------------------------------------------------
	// 1. CRUD USUARIOS
	// -------------------------------------------------------------
	t.Run("CRUD Usuarios", func(t *testing.T) {
		// POST /usuarios (crear)
		userBody, _ := json.Marshal(map[string]string{
			"nombre": "Justina Test",
			"email":  "justina@example.com",
		})
		resp, err := client.Post(ts.URL+"/usuarios", "application/json", bytes.NewBuffer(userBody))
		if err != nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
			t.Fatalf("Error en POST /usuarios: %v, status: %d", err, resp.StatusCode)
		}

		var user struct {
			ID int `json:"id"`
		}
		json.NewDecoder(resp.Body).Decode(&user)
		resp.Body.Close()

		if user.ID == 0 {
			t.Fatalf("Esperado ID de usuario válido, obtenido 0")
		}
		userID = user.ID

		// GET /usuarios (listar todos)
		resp, err = client.Get(ts.URL + "/usuarios")
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Error en GET /usuarios: %v", err)
		}
		resp.Body.Close()

		// GET /usuarios/{id} (buscar por ID)
		resp, err = client.Get(fmt.Sprintf("%s/usuarios/%d", ts.URL, userID))
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
		}
	})

	// -------------------------------------------------------------
	// 2. CRUD PRENDAS
	// -------------------------------------------------------------
	t.Run("CRUD Prendas", func(t *testing.T) {
		// POST /prendas (Crear)
		prendaBody, _ := json.Marshal(map[string]interface{}{
			"usuario_id": userID,
			"nombre":     "Campera Jean",
			"categoria":  "Abrigos",
		})
		resp, err := client.Post(ts.URL+"/prendas", "application/json", bytes.NewBuffer(prendaBody))
		if err != nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
			t.Fatalf("Error en POST /prendas: %v, status: %d", err, resp.StatusCode)
		}

		var prenda struct {
			ID int `json:"id"`
		}
		json.NewDecoder(resp.Body).Decode(&prenda)
		resp.Body.Close()
		prendaID = prenda.ID

		// GET /prendas (Listar todas / por usuario)
		resp, err = client.Get(ts.URL + "/prendas")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
		}
	})

	// -------------------------------------------------------------
	// 3. CRUD LOOKS
	// -------------------------------------------------------------
	t.Run("CRUD Looks", func(t *testing.T) {
		// POST /looks (Crear)
		lookBody, _ := json.Marshal(map[string]interface{}{
			"usuario_id":         userID,
			"nombre":             "Look Urbano",
			"descripcion":        "Casual",
			"ocasion":            "Dia",
			"temporada":          "Primavera",
			"prenda_estrella_id": prendaID,
			"url_imagen":         "http://example.com/img.jpg",
			"estilo":             "Casual",
		})
		resp, err := client.Post(ts.URL+"/looks", "application/json", bytes.NewBuffer(lookBody))
		if err != nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
			t.Fatalf("Error en POST /looks: %v, status: %d", err, resp.StatusCode)
		}

		var look struct {
			ID int `json:"id"`
		}
		json.NewDecoder(resp.Body).Decode(&look)
		resp.Body.Close()
		lookID = look.ID

		// GET /looks (listar todos)
		resp, err = client.Get(ts.URL + "/looks")
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Error en GET /looks: %v", err)
		}
		resp.Body.Close()

		// GET /looks/{id} (buscar por ID)
		resp, err = client.Get(fmt.Sprintf("%s/looks/%d", ts.URL, lookID))
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
		}
	})

	// -------------------------------------------------------------
	// 4. INTERACCIONES (ME GUSTA)
	// -------------------------------------------------------------
	t.Run("Interacciones MeGusta", func(t *testing.T) {
		likeBody, _ := json.Marshal(map[string]interface{}{
			"usuario_id": userID,
			"look_id":    lookID,
		})

		// POST /me-gusta (dar me gusta)
		resp, err := client.Post(ts.URL+"/me-gusta", "application/json", bytes.NewBuffer(likeBody))
		if err != nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
			t.Fatalf("Error en POST /me-gusta: %v", err)
		}
		resp.Body.Close()

		// DELETE /me-gusta (quitar me gusta)
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/me-gusta", bytes.NewBuffer(likeBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil || (resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK) {
			t.Fatalf("Error en DELETE /me-gusta: %v", err)
		}
		resp.Body.Close()
	})

	// -------------------------------------------------------------
	// 5. BORRADOS Y MANEJO DE ERRORES (404)
	// -------------------------------------------------------------
	t.Run("Eliminacion y Errores 404", func(t *testing.T) {
		// DELETE /looks/{id} (Eliminar Look creado)
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/looks/%d", ts.URL, lookID), nil)
		resp, err := client.Do(req)
		if err != nil || (resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK) {
			t.Fatalf("Error al eliminar look: %v", err)
		}
		resp.Body.Close()

		// Probar recurso que no existe --> 404 Not Found
		resp, err = client.Get(ts.URL + "/looks/99999")
		if err == nil {
			resp.Body.Close()
		}
	})
}
