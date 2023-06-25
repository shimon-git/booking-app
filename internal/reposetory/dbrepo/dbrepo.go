package dbrepo

import (
	"database/sql"

	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/reposetory"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *sql.DB) reposetory.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) reposetory.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
