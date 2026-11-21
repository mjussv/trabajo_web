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
		uniqueNombre := fmt.Sprintf("Jane Doe %d", time.Now().UnixNano())
		user, err := queries.CrearUsuario(ctx, CrearUsuarioParams{
			Nombre:   uniqueNombre,
			Email:    uniqueEmail,
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("Error al crear usuario: %v", err)
		}
		usuarioID = user.ID

		obtained, err := queries.GetUsuariosByID(ctx, usuarioID)
		if err != nil || obtained.ID != usuarioID {
			t.Fatalf("Error al obtener usuario: %v", err)
		}

		updatedEmail := fmt.Sprintf("updated_%d@example.com", time.Now().UnixNano())
		updated, err := queries.ActualizarUsuario(ctx, ActualizarUsuarioParams{
			ID:     usuarioID,
			Nombre: uniqueNombre + "_edit",
			Email:  updatedEmail,
		})
		if err != nil || updated.Email != updatedEmail {
			t.Fatalf("Error al actualizar usuario: %v", err)
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

		updatedPrenda, err := queries.ActualizarPrenda(ctx, ActualizarPrendaParams{
			ID:        prendaID,
			UsuarioID: usuarioID,
			Nombre:    "Campera de Cuero Negra",
			Categoria: "Abrigos",
			Column5:   "https://example.com/campera_negra.jpg",
		})
		if err != nil || updatedPrenda.Nombre != "Campera de Cuero Negra" {
			t.Fatalf("Error al actualizar prenda: %v", err)
		}

		encontradas, err := queries.BuscarPrendas(ctx, sql.NullString{String: "Cuero Negra", Valid: true})
		if err != nil || len(encontradas) == 0 {
			t.Fatalf("Error al buscar prendas: %v", err)
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

		updatedLook, err := queries.ActualizarLook(ctx, ActualizarLookParams{
			ID:               lookID,
			UsuarioID:        usuarioID,
			Nombre:           "Outfit Urbano Renovado",
			Descripcion:      sql.NullString{String: "Actualizado", Valid: true},
			Ocasion:          "Noche",
			Temporada:        "Invierno",
			PrendaEstrellaID: sql.NullInt32{Int32: prendaID, Valid: true},
			UrlImagen:        "https://example.com/outfit2.jpg",
			Estilo:           sql.NullString{String: "Streetwear", Valid: true},
		})
		if err != nil || updatedLook.Nombre != "Outfit Urbano Renovado" {
			t.Fatalf("Error al actualizar look: %v", err)
		}

		looksBuscados, err := queries.BuscarLooks(ctx, BuscarLooksParams{
			UsuarioID: usuarioID,
			Column2:   sql.NullString{String: "Renovado", Valid: true},
		})
		if err != nil || len(looksBuscados) == 0 {
			t.Fatalf("Error al buscar looks: %v", err)
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

		rows, err := queries.QuitarMeGusta(ctx, QuitarMeGustaParams{
			UsuarioID: usuarioID,
			LookID:    lookID,
		})
		if err != nil || rows != 1 {
			t.Fatalf("Error al quitar me gusta: %v (rows=%d)", err, rows)
		}
	})

	// TEST LIMPIAR CLÓSET Y BORRADO
	t.Run("LimpiarCloset_Y_Borrado", func(t *testing.T) {
		if err := queries.EliminarLook(ctx, lookID); err != nil {
			t.Fatalf("Error al eliminar look: %v", err)
		}

		err := queries.LimpiarClosetPorUsuario(ctx, usuarioID)
		if err != nil {
			t.Fatalf("Error al limpiar el clóset: %v", err)
		}

		prendas, err := queries.ListarPrendasPorUsuario(ctx, usuarioID)
		if err != nil || len(prendas) != 0 {
			t.Fatalf("Se esperaban 0 prendas tras limpiar el clóset")
		}

		if err := queries.EliminarUsuario(ctx, usuarioID); err != nil {
			t.Fatalf("Error al eliminar usuario: %v", err)
		}
	})
}
