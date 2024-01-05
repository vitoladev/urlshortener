package database

import "fmt"

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
