package main

import (
	"log"
	"my_notes_project/internal/api"
	"my_notes_project/internal/database"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	logger := logrus.New()

	lvl, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}

	logger.SetLevel(lvl)

	db, err := database.NewSQLiteDatabase("./database/noteuser.db")
	if err != nil {
		panic(err)
	}
	defer db.CloseSQLiteDatabase()

	restAPI := api.NewRestAPI(database.DBRepository(&db), logger)

	err = restAPI.HandlersInit()
	if err != nil {
		panic(err)
	}

	log.Fatal(restAPI.Listen("0.0.0.0:8080"))

}
