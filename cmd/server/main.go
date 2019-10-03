//go:generate protoc -I=./../../internal/model/ --go_out=plugins=grpc:./../../internal/model/ ./../../internal/model/model.proto
package main

import (
	"flag"
	"log"

	"gitlab.com/toby3d/test/internal/db"
	"gitlab.com/toby3d/test/internal/handler"
	"gitlab.com/toby3d/test/internal/server"
	"gitlab.com/toby3d/test/internal/store"
)

var (
	flagAddr = flag.String("addr", ":2368", "set specific address and port for server instance")
	flagDB   = flag.String(
		"db", `host=/var/run/postgresql dbname=testing sslmode=disable`,
		"set specific parameters for connecting to database",
	)
)

func main() {
	flag.Parse()

	dataBase, err := db.Open(*flagDB)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer dataBase.Close()

	if err = db.AutoMigrate(dataBase); err != nil {
		log.Fatalln(err.Error())
	}

	srv, err := server.NewServer(
		*flagAddr, handler.NewHandler(store.NewCartStore(dataBase), store.NewProductStore(dataBase)),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err = srv.Start(); err != nil {
		log.Fatalln(err.Error())
	}
}
