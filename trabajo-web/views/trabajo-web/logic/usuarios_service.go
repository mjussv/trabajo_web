package logic

import (
	"context"
	"errors"
	"strings"

	db "trabajo-web/db/sqlc"
)

type UsuariosService struct {
	queries *db.Queries
}

func NewUsuariosService(queries *db.Queries) *UsuariosService {
	return &UsuariosService{
		queries: queries,
	}
}

type CreateUsuariosParams struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"` // Mantenemos la struct para compatibilidad con el Handler
}

func (s *UsuariosService) CreateUsuarios(ctx context.Context, params CreateUsuariosParams) (db.Usuario, error) {
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Usuario{}, errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(params.Email) == "" {
		return db.Usuario{}, errors.New("el email no puede estar vacío")
	}

	return s.queries.CrearUsuario(ctx, db.CrearUsuarioParams{
		Nombre:   params.Nombre,
		Email:    params.Email,
		Password: params.Password,
	})
}

func (s *UsuariosService) ListUsuarios(ctx context.Context) ([]db.Usuario, error) {
	return s.queries.ListarUsuarios(ctx)
}

func (s *UsuariosService) GetUsuariosByID(ctx context.Context, id int32) (db.Usuario, error) {
	if id <= 0 {
		return db.Usuario{}, errors.New("el ID de usuario debe ser mayor a 0")
	}
	return s.queries.GetUsuariosByID(ctx, id)
}

func (s *UsuariosService) GetUsuariosByEmail(ctx context.Context, email string) (db.Usuario, error) {
	return s.queries.GetUsuariosByEmail(ctx, email)
}

func (s *UsuariosService) GetUsuariosByNombre(ctx context.Context, nombre string) (db.Usuario, error) {
	return s.queries.GetUsuariosByNombre(ctx, nombre)
}

type UpdateUsuariosParams struct {
	ID     int32  `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func (s *UsuariosService) UpdateUsuarios(ctx context.Context, params UpdateUsuariosParams) (db.Usuario, error) {
	if params.ID <= 0 {
		return db.Usuario{}, errors.New("el ID de usuario debe ser mayor a 0")
	}
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Usuario{}, errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(params.Email) == "" {
		return db.Usuario{}, errors.New("el email no puede estar vacío")
	}

	return s.queries.ActualizarUsuario(ctx, db.ActualizarUsuarioParams{
		ID:     params.ID,
		Nombre: params.Nombre,
		Email:  params.Email,
	})
}

func (s *UsuariosService) DeleteUsuarios(ctx context.Context, id int32) error {
	if id <= 0 {
		return errors.New("ID de usuario inválido")
	}
	return s.queries.EliminarUsuario(ctx, id)
}
