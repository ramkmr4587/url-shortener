package main

//https://github.com/ramkmr4587/url-shortener.git
import (
	"log"
	"net/http"

	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"

	restful "github.com/emicklei/go-restful/v3"
)

func main() {
	store := storage.NewStore()
	svc := service.NewURLService(store)
	api := handler.NewHandler(svc)

	container := restful.NewContainer()
	api.Register(container)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", container))
}
