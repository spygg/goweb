package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/solenovex/it/common"
	"github.com/solenovex/it/controller"
	"github.com/solenovex/it/middleware"

	_ "net/http/pprof"
)

const (
	host     = ""
	port     = 5432
	user     = ""
	password = ""
	dbname   = "demo"
)

func init() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	common.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ctx := context.Background()
	err = common.Db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected!")
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &middleware.BasicAuthMiddleware{},
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./wwwroot"))))
	controller.RegisterRoutes()

	log.Println("Server starting...")
	go http.ListenAndServe(":8000", nil)
	server.ListenAndServe()
}
