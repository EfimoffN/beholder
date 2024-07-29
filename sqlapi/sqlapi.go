package sqlapi

import (
	"context"
	"errors"
	"log"

	"github.com/EfimoffN/beholder/types"
	"github.com/jmoiron/sqlx"
)

var (
	ErrSessionsNotFound = errors.New("Sessions not found by ID")

	ErrClientFoundAlotOf = errors.New("Many clients have been detected")
)

type API struct {
	db *sqlx.DB
}

func NewAPI(db *sqlx.DB) *API {
	return &API{
		db: db,
	}
}

func (api *API) GetSessionsByID(sessionid string) (*types.PrjSessionRow, error) {
	return getSessionByID(context.Background(), api.db, sessionid)
}

func getSessionByID(ctx context.Context, db TxContext, sessionid string) (*types.PrjSessionRow, error) {
	session := []types.PrjSessionRow{}

	err := db.SelectContext(ctx, &session, `SELECT * FROM prj_session WHERE prj_session.sessionid = $1`, sessionid)
	if err != nil {
		log.Println("getClientSessionsByID api.db.Select failed with an error: ", err.Error())
		return nil, err
	}

	if len(session) == 0 {
		return nil, ErrSessionsNotFound
	}

	return &session[0], nil
}
