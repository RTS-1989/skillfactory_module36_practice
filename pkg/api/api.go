package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// API Программный интерфейс сервера GoNews
type API struct {
	db      storage.PostsInterface
	router  *mux.Router
	errChan chan error
}

// New Конструктор объекта API
func New(db storage.PostsInterface, errChan chan error) *API {
	log.Println("Create new API instance")
	api := API{
		db:      db,
		errChan: errChan,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// endpoints Регистрация обработчиков API.
func (api *API) endpoints() {
	log.Println("Register endpoints")
	api.router.Use(api.HeadersMiddleware)
	api.router.HandleFunc("/news/{n}", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Router получаем объект роутера
func (api *API) Router() *mux.Router {
	return api.router
}

// postsHandler Получение опционального количества публикаций
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Start request for posts")
	countPosts, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/news/"))
	if err != nil {
		api.errChan <- err
		return
	}
	posts, err := api.db.Posts(countPosts)
	if err != nil {
		api.errChan <- err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		api.errChan <- err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		api.errChan <- err
	}
}
