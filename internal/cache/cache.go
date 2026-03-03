package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	dir     string
	expires time.Duration
}

type Option func(*Cache)

func WithTTL(ttl time.Duration) Option {
	return func(c *Cache) {
		c.expires = ttl
	}
}

func New(opts ...Option) (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	c := &Cache{
		dir:     filepath.Join(home, ".echo-cli", "cache"),
		expires: 24 * time.Hour,
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.setup(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cache) setup() error {
	dirs := []string{
		filepath.Join(c.dir, "api"),
		filepath.Join(c.dir, "files"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create cache dir %s: %w", d, err)
		}
	}
	return nil
}

func (c *Cache) hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func (c *Cache) apiPath(key string) string {
	hash := c.hashKey(key)
	return filepath.Join(c.dir, "api", hash)
}

func (c *Cache) filesPath(key string) string {
	hash := c.hashKey(key)
	return filepath.Join(c.dir, "files", hash)
}

func (c *Cache) isExpired(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	age := time.Since(info.ModTime())
	return age > c.expires
}

func (c *Cache) GetAPI(url string) ([]byte, error) {
	path := c.apiPath(url)
	if c.isExpired(path) {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(path)
}

func (c *Cache) SetAPI(url string, data []byte) error {
	path := c.apiPath(url)
	return os.WriteFile(path, data, 0644)
}

func (c *Cache) GetFile(url string) (string, error) {
	path := c.filesPath(url)
	if c.isExpired(path) {
		return "", os.ErrNotExist
	}
	return path, nil
}

func (c *Cache) SetFile(url string, data []byte) (string, error) {
	path := c.filesPath(url)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (c *Cache) FetchWithCache(url string) ([]byte, error) {
	if data, err := c.GetAPI(url); err == nil {
		return data, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_ = c.SetAPI(url, data)
	return data, nil
}
