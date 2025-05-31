package service

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"

	"url-shortener/internal/storage"
)

type URLService struct {
	store *storage.Store
}

func NewURLService(s *storage.Store) *URLService {
	return &URLService{store: s}
}

func (s *URLService) ShortenURL(original string) string {
	s.store.Mutex.Lock()
	defer s.store.Mutex.Unlock()

	if short, exists := s.store.URLToShort[original]; exists {
		return short
	}

	hash := md5.Sum([]byte(original))
	short := hex.EncodeToString(hash[:])[:6]

	s.store.URLToShort[original] = short
	s.store.ShortToURL[short] = original

	u, _ := url.Parse(original)
	domain := strings.TrimPrefix(u.Hostname(), "www.")
	s.store.DomainHits[domain]++

	return short
}

func (s *URLService) GetOriginalURL(short string) (string, bool) {
	s.store.Mutex.RLock()
	defer s.store.Mutex.RUnlock()
	url, exists := s.store.ShortToURL[short]
	return url, exists
}

func (s *URLService) GetTopDomains(limit int) map[string]int {
	s.store.Mutex.RLock()
	defer s.store.Mutex.RUnlock()

	result := make(map[string]int)
	for k, v := range s.store.DomainHits {
		result[k] = v
	}
	return result
}
