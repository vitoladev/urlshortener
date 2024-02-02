package main

import (
	"fmt"
	"urlshortener/internal/handler"
	"urlshortener/internal/repository"
	"urlshortener/internal/server"
)

func main() {
	db := repository.NewDB()
	repo := repository.NewRepository(db)
	urlRepo := repository.NewUrlRepository(repo)
	urlHandler := handler.NewUrlHandler(urlRepo)

	server := server.NewServer(urlHandler)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
