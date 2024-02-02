package repository

import "fmt"

type ShortenUrlPayload struct {
	OriginalUrl string
	ShortUrl    string
}

type UrlRepository interface {
	ShortenUrl(url ShortenUrlPayload) error
	GetOriginalUrl(shortUrl string) (string, error)
	GetShortUrl(originalUrl string) (string, error)
}

func NewUrlRepository(r *Repository) UrlRepository {
	return &urlRepository{
		Repository: r,
	}
}

type urlRepository struct {
	*Repository
}

func (r *urlRepository) ShortenUrl(url ShortenUrlPayload) error {
	query := "INSERT INTO url (original_url, short_url) VALUES ($1, $2)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(url.OriginalUrl, url.ShortUrl)

	return err
}

func (r *urlRepository) GetOriginalUrl(shortUrl string) (string, error) {
	query := "SELECT original_url FROM url WHERE short_url = $1"
	stmt, err := r.db.Prepare(query)
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

func (r *urlRepository) GetShortUrl(originalUrl string) (string, error) {
	query := "SELECT short_url FROM url WHERE original_url = $1"
	stmt, err := r.db.Prepare(query)
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
