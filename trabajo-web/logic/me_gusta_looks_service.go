package logic

import (
	"context"

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
	return s.queries.DarMeGusta(ctx, db.DarMeGustaParams{
		UsuarioID: params.UsuarioID,
		LookID:    params.LookID,
	})
}

func (s *MeGustaLooksService) DeleteMeGusta(ctx context.Context, params CreateMeGustaParams) error {
	return s.queries.QuitarMeGusta(ctx, db.QuitarMeGustaParams{
		UsuarioID: params.UsuarioID,
		LookID:    params.LookID,
	})
}
