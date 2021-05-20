package main

import (
	"context"
	"log"
	"os"

	"github.com/arsura/moonbase-service/pkg/models/pgsql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	pg       *pgsql.Repositories
}

func main() {
	server := fiber.New()
	server.Use(cors.New())

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		errorLog.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		pg: &pgsql.Repositories{
			Currencies: &pgsql.DB{Conn: pool},
		},
	}

	app.routes(server)

	if err := server.Listen(":8080"); err != nil {
		errorLog.Printf("Unable to start server: %v\n", err)
		os.Exit(1)
	}
}
