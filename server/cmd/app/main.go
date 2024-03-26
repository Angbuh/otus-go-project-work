package main

import (
	"log"
	"my_notes_project/internal/api"
	"my_notes_project/internal/core"
	"my_notes_project/internal/database"

	"github.com/sirupsen/logrus"
)

func main() {
	//Присваиваем переменной функцию гетконфиг, обрабатываем ошибку
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	// Создаем новый логгер
	logger := logrus.New()

	// Парсим логлевел и записываем его как строку,возвращает уровень и ошибку
	lvl, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}

	// Устанавливаем уровень логгирования
	logger.SetLevel(lvl)

	// Настраиваем путь до базы данных
	db, err := database.NewSQLiteDatabase(config.DatabasePath, logger)
	if err != nil {
		panic(err)
	}
	// Закрываем базу данных
	defer db.CloseSQLiteDatabase()

	// Создаем новый апи и кор
	core := core.NewTheCore(db, logger)
	restAPI := api.NewRestAPI(core, logger)

	// Обрабатываем хендлеры на ошибку
	err = restAPI.HandlersInit()
	if err != nil {
		panic(err)
	}

	// Для прослушивания и приема входящих запросов
	log.Fatal(restAPI.Listen("0.0.0.0:8080"))

}
