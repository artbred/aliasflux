package postgres

import (
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

var database *sqlx.DB
var once sync.Once

func Connection() *sqlx.DB {
	if database == nil {
		once.Do(Init)
	}

	return database
}

func Init() {
	var err error

	url := config.ConnectionURLBuilder("postgres")

	database, err = sqlx.Connect("pgx", url)
	if err != nil {
		common.Logger.WithError(err).Errorf("error connecting to postgresql %s", url)
		time.Sleep(5 * time.Second)
		Init()
	}

	if err = database.Ping(); err != nil {
		common.Logger.WithError(err).Errorf("error ping to postgresql %s", url)
		time.Sleep(5 * time.Second)
		Init()
	}

	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(20)

	common.Logger.Infof("connected to postgresql %s", url)
}

func init() {
	Init()
}
