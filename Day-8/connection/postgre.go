package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	var err error
	databaseUrl := "postgres://postgres:Rakhapostgre@localhost:5432/Personal-Web"

	Conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Success connect to database")
}

// postgres://{user}:{password}@{host}:{port}/{database}
// user = user postgres
// password = password postgres
// host = host postgres
// port = port postgres
// database = database postgres
