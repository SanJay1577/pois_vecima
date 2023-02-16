package store

import (
	"strings"

	"pois/models"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"

	"github.com/jmoiron/sqlx"
)

type Alias interface {
	FindAlias(string) ([]models.Alias, error)
	CreateAlias(models.Alias) (models.Alias, error)
	DeleteAlias(string, string) error
}

// SQLAlias is a DB-backed concrete store.
type SQLAlias struct {
	DB  *sqlx.DB
	Log *zaplogger.Logger
}

//FindAlias fetch alias names based on channel name
func (s *SQLAlias) FindAlias(channelname string) ([]models.Alias, error) {
	alias := []models.Alias{}
	err := s.DB.Select(&alias, "SELECT channelname, aliasname FROM channel_alias where channelname= $1",
		channelname)
	return alias, err

}

//FindAllAlias fetch all channel alias name from DB
func (s *SQLAlias) FindAllAlias() ([]models.Alias, error) {
	alias := []models.Alias{}
	err := s.DB.Select(&alias, "SELECT channelname, aliasname FROM channel_alias")
	return alias, err

}

// CreateAlias create channel alias in DB.
func (s *SQLAlias) CreateAlias(request models.Alias) ([]models.Alias, error) {

	alias := []models.Alias{}

	channelname := strings.TrimSpace(request.Channelname)
	aliasname := strings.TrimSpace(request.AliasName)

	err := s.DB.Select(
		&alias,
		`INSERT INTO channel_alias(channelname, aliasname)
		VALUES ($1, $2)
		RETURNING channelname, aliasname`, channelname, aliasname,
	)

	return alias, err
}

// DeleteAlias deletes a channel alias row from DB
func (s *SQLAlias) DeleteAlias(channel string, alias string) error {

	_, err := s.DB.Exec(`DELETE FROM channel_alias	WHERE channelname = $1 and aliasname= $2`, channel, alias)
	if err != nil {
		return err
	}

	return err
}
