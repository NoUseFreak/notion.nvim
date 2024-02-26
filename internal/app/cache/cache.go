package cache

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type Cache struct {
	db *sql.DB
}

func GetDBPath(key string) string {
	return filepath.Join(os.TempDir(), key+".db")
}

func New(key string) *Cache {
	path := GetDBPath(key)
	isNew := false
	if _, err := os.Stat(path); err != nil {
		isNew = true
	}
	logrus.Debugf("Cache path: %s", path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	if isNew {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS cache (key TEXT PRIMARY KEY, value TEXT)")
		if err != nil {
			panic(err)
		}
	}

	return &Cache{
		db: db,
	}
}

func (c *Cache) Get(key string) (string, error) {
	var value string
	err := c.db.QueryRow("SELECT value FROM cache WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *Cache) Set(key, value string) error {
	_, err := c.db.Exec("INSERT OR REPLACE INTO cache (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		logrus.Errorf("Error setting cache: %s", err)
	}
	return err
}
