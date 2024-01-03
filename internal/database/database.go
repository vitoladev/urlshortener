package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	ShortenUrl(urlData UrlData) error
	GetOriginalUrl(shortUrl string) (string, error)
	GetShortUrl(originalUrl string) (string, error)
}	

type service struct {
	db *sql.DB
}

type UrlData struct {
	OriginalUrl string
	ShortUrl    string
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) ShortenUrl(urlData UrlData) error {
	query := "INSERT INTO url (original_url, short_url) VALUES ($1, $2)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(urlData.OriginalUrl, urlData.ShortUrl)

	return err
}

func (s *service) GetOriginalUrl(shortUrl string) (string, error) {
	query := "SELECT original_url FROM url WHERE short_url = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(shortUrl)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var originalUrl string
		if err := rows.Scan(&originalUrl); err != nil {
			return "", err
		}
		return originalUrl, nil
	}

	// Handle the case where no rows were found (no match for the short URL)
	return "", fmt.Errorf("no original URL found for short URL: %s", shortUrl)
}

func (s *service) GetShortUrl(originalUrl string) (string, error) {
	query := "SELECT short_url FROM url WHERE original_url = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(originalUrl)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var shortUrl string
		if err := rows.Scan(&shortUrl); err != nil {
			return "", err
		} 
		return shortUrl, nil
	}

	return "", fmt.Errorf("no short URL found for original URL: %s", originalUrl)
}
