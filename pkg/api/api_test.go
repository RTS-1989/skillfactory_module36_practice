package api

import (
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAPI_postsHandler(t *testing.T) {
	e := godotenv.Load("../../.env") //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("test_db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, dbHost, dbPort, dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postsChan := make(chan storage.Post)
	errChan := make(chan error)

	db, _ := postgres.New(ctx, connString, postsChan, errChan)
	api := New(db, errChan)

	req := httptest.NewRequest(http.MethodGet, "/news/10", nil)
	// Создаём объект для записи ответа обработчика.
	resp := httptest.NewRecorder()
	api.router.ServeHTTP(resp, req)

	if !(resp.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", resp.Code, http.StatusOK)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в массив заказов.
	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве ровно один элемент.
	const wantLen = 10
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
}

func TestNew(t *testing.T) {
	type args struct {
		db      storage.PostsInterface
		errChan chan error
	}
	tests := []struct {
		name string
		args args
		want *API
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db, tt.args.errChan); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
