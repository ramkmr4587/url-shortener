package service

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type URLServiceTestSuite struct {
	suite.Suite
	Service *URLService
	Store   *mockStore
}

func TestURLServiceTestSuite(t *testing.T) {
	suite.Run(t, new(URLServiceTestSuite))
}

func (suite *URLServiceTestSuite) SetupTest() {
	suite.Store = &mockStore{
		URLToShort: make(map[string]string),
		ShortToURL: make(map[string]string),
		DomainHits: make(map[string]int),
	}
	suite.Service = NewURLService((*storage.Store)(suite.Store))
}

func (suite *URLServiceTestSuite) TestShortenURLSuccessNewURL() {
	original := "https://example.com"
	hash := md5.Sum([]byte(original))
	expectedShort := hex.EncodeToString(hash[:])[:6]
	expectedDomain := "example.com"

	result := suite.Service.ShortenURL(original)

	assert.Equal(suite.T(), expectedShort, result)
	assert.Equal(suite.T(), expectedShort, suite.Store.URLToShort[original])
	assert.Equal(suite.T(), original, suite.Store.ShortToURL[expectedShort])
	assert.Equal(suite.T(), 1, suite.Store.DomainHits[expectedDomain])
}

func (suite *URLServiceTestSuite) TestShortenURLExistingURL() {
	original := "https://example.com"
	hash := md5.Sum([]byte(original))
	expectedShort := hex.EncodeToString(hash[:])[:6]
	expectedDomain := "example.com"
	suite.Store.URLToShort[original] = expectedShort
	suite.Store.ShortToURL[expectedShort] = original
	suite.Store.DomainHits[expectedDomain] = 1

	result := suite.Service.ShortenURL(original)

	assert.Equal(suite.T(), expectedShort, result)
	assert.Equal(suite.T(), 2, suite.Store.DomainHits[expectedDomain])
}

func (suite *URLServiceTestSuite) TestShortenURLEmptyURL() {
	result := suite.Service.ShortenURL("")

	assert.Equal(suite.T(), "", result)
	assert.Empty(suite.T(), suite.Store.URLToShort)
	assert.Empty(suite.T(), suite.Store.ShortToURL)
	assert.Empty(suite.T(), suite.Store.DomainHits)
}

func (suite *URLServiceTestSuite) TestGetOriginalURLSuccess() {
	short := "abc123"
	original := "https://example.com"
	suite.Store.ShortToURL[short] = original

	result, exists := suite.Service.GetOriginalURL(short)

	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), original, result)
}

func (suite *URLServiceTestSuite) TestGetOriginalURLNotFound() {
	result, exists := suite.Service.GetOriginalURL("invalid")

	assert.False(suite.T(), exists)
	assert.Empty(suite.T(), result)
}

func (suite *URLServiceTestSuite) TestGetTopDomainsSuccess() {
	suite.Store.DomainHits = map[string]int{
		"example.com": 10,
		"test.com":    5,
		"other.com":   3,
	}

	result := suite.Service.GetTopDomains(3)

	assert.Equal(suite.T(), suite.Store.DomainHits, result)
	assert.Len(suite.T(), result, 3)
	assert.Equal(suite.T(), 10, result["example.com"])
	assert.Equal(suite.T(), 5, result["test.com"])
	assert.Equal(suite.T(), 3, result["other.com"])
}

func (suite *URLServiceTestSuite) TestGetTopDomainsEmpty() {
	result := suite.Service.GetTopDomains(3)

	assert.Empty(suite.T(), result)
}

// mockStore mimics storage.Store
type mockStore storage.Store

func (m *mockStore) Lock() {
	// No-op for testing
}

func (m *mockStore) Unlock() {
	// No-op for testing
}

func (m *mockStore) RLock() {
	// No-op for testing
}

func (m *mockStore) RUnlock() {
	// No-op for testing
}
