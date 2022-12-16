package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/rss"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type server struct {
	db  storage.PostsInterface
	api *api.API
}

func main() {
	log.Println("Start GoNews")
	var srv server

	e := godotenv.Load("../.env") //Загрузить файл .env из корня приложения
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, dbHost, dbPort, dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postsChan := make(chan storage.Post)
	errChan := make(chan error)
	db, err := postgres.New(ctx, connString, postsChan, errChan)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()

	srv.db = db
	srv.api = api.New(db, errChan)

	configFile := "./config.json"
	sourceRss, err := rss.NewSourceRss(configFile, postsChan, errChan)

	// Получаем список последних публикаций из базы данных (если они есть). Чтобы не дублировать данные в бд.
	lastPubTimeMap, _ := srv.db.GetLastPubDateForSources(sourceRss.SourceRssList)
	sourceRss.LastPubTimeFromDB = lastPubTimeMap

	if err != nil {
		log.Fatal(err)
	}

	sourceRss.RunGetSourcesInfo()
	srv.db.RunInsertPosts()
	runErrorsCheck(errChan)

	err = http.ListenAndServe(":8081", srv.api.Router())
	if err != nil {
		errChan <- err
		return
	}
}

// runErrorsCheck читает ошибки из канало ошибок
func runErrorsCheck(errChan chan error) {
	go func() {
		for {
			select {
			case err := <-errChan:
				log.Println(err)
			}
		}
	}()
}
