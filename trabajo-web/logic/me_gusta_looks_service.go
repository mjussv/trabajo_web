package logic

import (
	"context"
	"errors"

	db "trabajo-web/db/sqlc"
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
		return errors.New("ID de usuario y ID de look deben ser mayores a 0")
	}
	look, err := s.queries.ObtenerLookConDetalle(ctx, params.LookID)
	if err != nil {
		return errors.New("el look no existe")
	}
	if look.UsuarioID == params.UsuarioID {
		return errors.New("no podés darle me gusta a tu propio look")
	}

	return s.queries.DarMeGusta(ctx, db.DarMeGustaParams{
		UsuarioID: params.UsuarioID,
		LookID:    params.LookID,
	})
}

func (s *MeGustaLooksService) DeleteMeGusta(ctx context.Context, params CreateMeGustaParams) error {
	if params.UsuarioID <= 0 || params.LookID <= 0 {
		return errors.New("ID de usuario y ID de look deben ser mayores a 0")
	}
	return s.queries.QuitarMeGusta(ctx, db.QuitarMeGustaParams{
		UsuarioID: params.UsuarioID,
		LookID:    params.LookID,
	})
}

func (s *MeGustaLooksService) ToggleMeGusta(ctx context.Context, usuarioID, lookID int32) error {
	yaDioLike, err := s.queries.UsuarioLeDioMeGusta(ctx, db.UsuarioLeDioMeGustaParams{
		UsuarioID: usuarioID,
		LookID:    lookID,
	})
	if err != nil {
		return err
	}
	if yaDioLike {
		return s.DeleteMeGusta(ctx, CreateMeGustaParams{UsuarioID: usuarioID, LookID: lookID})
	}
	return s.AddMeGusta(ctx, CreateMeGustaParams{UsuarioID: usuarioID, LookID: lookID})
}
