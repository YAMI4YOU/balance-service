package configure

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func MustInitDB(ctx context.Context, dbURL string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		panic("unable to connect to database")
	}

	for i := 0; i < 3; i++ {
		if err := conn.Ping(ctx); err != nil {
			if i == 2 {
				panic("unable to connect to database")
			}
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	log.Println("Successfully connected to database")

	return conn
}
