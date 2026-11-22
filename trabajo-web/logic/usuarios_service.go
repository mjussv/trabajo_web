package logic

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	db "trabajo-web/db/sqlc"

	"golang.org/x/crypto/bcrypt"
)

// ValidationError identifica errores de validación de entrada del usuario (HTTP 400).
type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return e.Msg
}

// Hash dummy precalculado para evitar timing attacks cuando el usuario no existe.
var dummyBcryptHash = []byte("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi7mN.8V6pM0Y7n4o6u0/8lE5wR2J1e")

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@([a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

func validarEmail(email string) error {
	if email == "" {
		return ValidationError{Msg: "el email no puede estar vacío"}
	}
	if utf8.RuneCountInString(email) > 255 {
		return ValidationError{Msg: "el email no puede superar los 255 caracteres"}
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email || !emailRegex.MatchString(email) {
		return ValidationError{Msg: "ingresá una dirección de email válida (ej: usuario@dominio.com)"}
	}
	return nil
}

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
	Password string `json:"password"`
}

func (s *UsuariosService) CreateUsuarios(ctx context.Context, params CreateUsuariosParams) (db.Usuario, error) {
	nombreLimpio := strings.TrimSpace(params.Nombre)
	emailLimpio := strings.ToLower(strings.TrimSpace(params.Email))

	if nombreLimpio == "" {
		return db.Usuario{}, ValidationError{Msg: "el nombre no puede estar vacío"}
	}
	if utf8.RuneCountInString(nombreLimpio) > 100 {
		return db.Usuario{}, ValidationError{Msg: "el nombre no puede superar los 100 caracteres"}
	}
	if err := validarEmail(emailLimpio); err != nil {
		return db.Usuario{}, err
	}
	if strings.TrimSpace(params.Password) == "" {
		return db.Usuario{}, ValidationError{Msg: "la contraseña no puede estar vacía"}
	}
	if utf8.RuneCountInString(params.Password) < 6 {
		return db.Usuario{}, ValidationError{Msg: "la contraseña debe tener al menos 6 caracteres"}
	}
	if len(params.Password) > 72 {
		return db.Usuario{}, ValidationError{Msg: "la contraseña no puede superar los 72 caracteres"}
	}

	if _, err := s.queries.GetUsuariosByNombre(ctx, nombreLimpio); err == nil {
		return db.Usuario{}, ValidationError{Msg: "ese nombre de usuario ya está en uso"}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return db.Usuario{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.Usuario{}, errors.New("error al procesar la contraseña")
	}

	return s.queries.CrearUsuario(ctx, db.CrearUsuarioParams{
		Nombre:   nombreLimpio,
		Email:    emailLimpio,
		Password: string(hashed),
	})
}

func (s *UsuariosService) VerificarPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *UsuariosService) VerificarPasswordDummy(plainPassword string) {
	_ = bcrypt.CompareHashAndPassword(dummyBcryptHash, []byte(plainPassword))
}

func (s *UsuariosService) ListUsuarios(ctx context.Context) ([]db.Usuario, error) {
	return s.queries.ListarUsuarios(ctx)
}

func (s *UsuariosService) GetUsuariosByID(ctx context.Context, id int32) (db.Usuario, error) {
	if id <= 0 {
		return db.Usuario{}, ValidationError{Msg: "el ID de usuario debe ser mayor a 0"}
	}
	return s.queries.GetUsuariosByID(ctx, id)
}

func (s *UsuariosService) GetUsuariosByEmail(ctx context.Context, email string) (db.Usuario, error) {
	return s.queries.GetUsuariosByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
}

func (s *UsuariosService) GetUsuariosByNombre(ctx context.Context, nombre string) (db.Usuario, error) {
	return s.queries.GetUsuariosByNombre(ctx, strings.TrimSpace(nombre))
}

type UpdateUsuariosParams struct {
	ID     int32  `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func (s *UsuariosService) UpdateUsuarios(ctx context.Context, params UpdateUsuariosParams) (db.Usuario, error) {
	nombreLimpio := strings.TrimSpace(params.Nombre)
	emailLimpio := strings.ToLower(strings.TrimSpace(params.Email))

	if params.ID <= 0 {
		return db.Usuario{}, ValidationError{Msg: "el ID de usuario debe ser mayor a 0"}
	}
	if nombreLimpio == "" {
		return db.Usuario{}, ValidationError{Msg: "el nombre no puede estar vacío"}
	}
	if utf8.RuneCountInString(nombreLimpio) > 100 {
		return db.Usuario{}, ValidationError{Msg: "el nombre no puede superar los 100 caracteres"}
	}
	if err := validarEmail(emailLimpio); err != nil {
		return db.Usuario{}, err
	}

	if existente, err := s.queries.GetUsuariosByNombre(ctx, nombreLimpio); err == nil && existente.ID != params.ID {
		return db.Usuario{}, ValidationError{Msg: "ese nombre de usuario ya está en uso"}
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return db.Usuario{}, err
	}

	return s.queries.ActualizarUsuario(ctx, db.ActualizarUsuarioParams{
		ID:     params.ID,
		Nombre: nombreLimpio,
		Email:  emailLimpio,
	})
}

func (s *UsuariosService) DeleteUsuarios(ctx context.Context, id int32) error {
	if id <= 0 {
		return ValidationError{Msg: "ID de usuario inválido"}
	}
	return s.queries.EliminarUsuario(ctx, id)
}
