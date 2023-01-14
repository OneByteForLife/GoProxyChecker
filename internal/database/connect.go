package database

import (
	"GoProxyChecker/internal/config"
	"context"
	"fmt"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func ConnectToDatabase() *pgxpool.Pool {
	conf := config.ReadConfig()

	// max_conns = 50 - ~175000 req/s
	urlConnect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=50", conf.Database.Username, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Dbname, conf.Database.Sslmode)

	logrus.Infof("Connecting by url - %s", urlConnect)

	dbPool, err := pgxpool.Connect(context.Background(), urlConnect)
	if err != nil {
		logrus.Errorf("Err connect to database - %s", err)
		return nil
	}

	return dbPool
}
