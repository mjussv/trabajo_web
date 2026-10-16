package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	connStr = "postgres://postgres:mysecretpassword@localhost:5432/outfits_db?sslmode=disable"
)

func TestSuiteCompleta(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("La base de datos no está respondiendo: %v", err)
	}

	queries := New(db)

	uniqueEmail := fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano())

	var usuarioID int32
	var prendaID int32
	var lookID int32

	// TEST USUARIOS
	t.Run("Usuarios_CRUD", func(t *testing.T) {
		user, err := queries.CrearUsuario(ctx, CrearUsuarioParams{
			Nombre: "Jane Doe",
			Email:  uniqueEmail,
		})
		if err != nil {
			t.Fatalf("Error al crear usuario: %v", err)
		}
		usuarioID = user.ID

		obtained, err := queries.ObtenerUsuario(ctx, usuarioID)
		if err != nil || obtained.ID != usuarioID {
			t.Fatalf("Error al obtener usuario: %v", err)
		}
	})

	// TEST PRENDAS
	t.Run("Prendas_CRUD", func(t *testing.T) {
		prenda, err := queries.CrearPrenda(ctx, CrearPrendaParams{
			UsuarioID: usuarioID,
			Nombre:    "Campera de Cuero",
			Categoria: "Abrigos",
		})
		if err != nil {
			t.Fatalf("Error al crear prenda: %v", err)
		}
		prendaID = prenda.ID

		prendas, err := queries.ListarPrendasPorUsuario(ctx, usuarioID)
		if err != nil || len(prendas) == 0 {
			t.Fatalf("Error al listar prendas del usuario: %v", err)
		}
	})

	// TEST LOOKS
	t.Run("Looks_CRUD", func(t *testing.T) {
		look, err := queries.CrearLook(ctx, CrearLookParams{
			UsuarioID:        usuarioID,
			Nombre:           "Outfit Urbano Noche",
			Descripcion:      sql.NullString{String: "Estilo casual chic ideal para salidas nocturnas.", Valid: true},
			Ocasion:          "Noche",
			Temporada:        "Otoño",
			PrendaEstrellaID: sql.NullInt32{Int32: prendaID, Valid: true},
			UrlImagen:        "https://example.com/outfit.jpg",
			Estilo:           sql.NullString{String: "Urbano", Valid: true},
		})
		if err != nil {
			t.Fatalf("Error al crear look: %v", err)
		}
		lookID = look.ID

		detalle, err := queries.ObtenerLookConDetalle(ctx, lookID)
		if err != nil || detalle.ID != lookID {
			t.Fatalf("Error al obtener detalle del look: %v", err)
		}
	})

	// TEST ME GUSTA
	t.Run("MeGusta_Interaccion", func(t *testing.T) {
		err := queries.DarMeGusta(ctx, DarMeGustaParams{
			UsuarioID: usuarioID,
			LookID:    lookID,
		})
		if err != nil {
			t.Fatalf("Error al dar me gusta: %v", err)
		}

		dioLike, err := queries.UsuarioLeDioMeGusta(ctx, UsuarioLeDioMeGustaParams{
			UsuarioID: usuarioID,
			LookID:    lookID,
		})
		if err != nil || !dioLike {
			t.Fatalf("Se esperaba que el usuario le haya dado Me Gusta al look")
		}

		count, err := queries.ContarMeGustaPorLook(ctx, lookID)
		if err != nil || count != 1 {
			t.Fatalf("Se esperaba 1 Me Gusta registrado, obtenido: %d", count)
		}
	})

	// TEST LIMPIAR CLÓSET
	t.Run("LimpiarCloset", func(t *testing.T) {
		err := queries.LimpiarClosetPorUsuario(ctx, usuarioID)
		if err != nil {
			t.Fatalf("Error al limpiar el clóset: %v", err)
		}

		prendas, err := queries.ListarPrendasPorUsuario(ctx, usuarioID)
		if err != nil || len(prendas) != 0 {
			t.Fatalf("Se esperaban 0 prendas tras limpiar el clóset")
		}
	})
}
