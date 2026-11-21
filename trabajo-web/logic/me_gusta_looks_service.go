package logic

import (
	"context"
	"database/sql"
	"errors"

	db "trabajo-web/db/sqlc"
)

var (
	ErrSelfLike     = errors.New("no podés darle me gusta a tu propio look")
	ErrLookNotFound = errors.New("el look no existe")
)

type MeGustaLooksService struct {
	queries *db.Queries
}

func NewMeGustaLooksService(queries *db.Queries) *MeGustaLooksService {
	return &MeGustaLooksService{queries: queries}
}

type CreateMeGustaParams struct {
	UsuarioID int32 `json:"usuario_id"`
	LookID    int32 `json:"look_id"`
}

func (s *MeGustaLooksService) AddMeGusta(ctx context.Context, params CreateMeGustaParams) error {
	if params.UsuarioID <= 0 || params.LookID <= 0 {
		return ValidationError{Msg: "ID de usuario y ID de look deben ser mayores a 0"}
	}
	look, err := s.queries.ObtenerLookConDetalle(ctx, params.LookID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrLookNotFound
	}
	if err != nil {
		return err
	}
	if look.UsuarioID == params.UsuarioID {
		return ErrSelfLike
	}

	return s.queries.DarMeGusta(ctx, db.DarMeGustaParams{
		UsuarioID: params.UsuarioID,
		LookID:    params.LookID,
	})
}

func (s *MeGustaLooksService) DeleteMeGusta(ctx context.Context, params CreateMeGustaParams) error {
	if params.UsuarioID <= 0 || params.LookID <= 0 {
		return ValidationError{Msg: "ID de usuario y ID de look deben ser mayores a 0"}
	}
	_, err := s.queries.QuitarMeGusta(ctx, db.QuitarMeGustaParams{
		UsuarioID: params.UsuarioID,
		LookID:    params.LookID,
	})
	return err
}

func (s *MeGustaLooksService) ToggleMeGusta(ctx context.Context, usuarioID, lookID int32) error {
	if usuarioID <= 0 || lookID <= 0 {
		return ValidationError{Msg: "ID de usuario y ID de look deben ser mayores a 0"}
	}

	look, err := s.queries.ObtenerLookConDetalle(ctx, lookID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrLookNotFound
	}
	if err != nil {
		return err
	}
	if look.UsuarioID == usuarioID {
		return ErrSelfLike
	}

	eliminadas, err := s.queries.QuitarMeGusta(ctx, db.QuitarMeGustaParams{
		UsuarioID: usuarioID,
		LookID:    lookID,
	})
	if err != nil {
		return err
	}
	if eliminadas > 0 {
		return nil
	}

	return s.queries.DarMeGusta(ctx, db.DarMeGustaParams{
		UsuarioID: usuarioID,
		LookID:    lookID,
	})
}

func (s *MeGustaLooksService) ContarMeGustaPorLook(ctx context.Context, lookID int32) (int64, error) {
	if lookID <= 0 {
		return 0, ValidationError{Msg: "ID de look inválido"}
	}
	return s.queries.ContarMeGustaPorLook(ctx, lookID)
}

func (s *MeGustaLooksService) UsuarioLeDioMeGusta(ctx context.Context, usuarioID, lookID int32) (bool, error) {
	if usuarioID <= 0 || lookID <= 0 {
		return false, ValidationError{Msg: "ID de usuario o look inválido"}
	}
	return s.queries.UsuarioLeDioMeGusta(ctx, db.UsuarioLeDioMeGustaParams{
		UsuarioID: usuarioID,
		LookID:    lookID,
	})
}
