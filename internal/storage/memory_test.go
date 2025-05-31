package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	Store *Store
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (suite *StoreTestSuite) SetupTest() {
	suite.Store = NewStore()
}

func (suite *StoreTestSuite) TestNewStoreInitialization() {
	assert.NotNil(suite.T(), suite.Store, "Store should not be nil")
	assert.NotNil(suite.T(), suite.Store.URLToShort, "URLToShort map should be initialized")
	assert.NotNil(suite.T(), suite.Store.ShortToURL, "ShortToURL map should be initialized")
	assert.NotNil(suite.T(), suite.Store.DomainHits, "DomainHits map should be initialized")
	assert.Empty(suite.T(), suite.Store.URLToShort, "URLToShort map should be empty")
	assert.Empty(suite.T(), suite.Store.ShortToURL, "ShortToURL map should be empty")
	assert.Empty(suite.T(), suite.Store.DomainHits, "DomainHits map should be empty")
}

func (suite *StoreTestSuite) TestConcurrentReadWriteURLToShort() {
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			original := "https://example.com/" + string(rune(i))
			short := "short" + string(rune(i))

			suite.Store.Mutex.Lock()
			suite.Store.URLToShort[original] = short
			suite.Store.Mutex.Unlock()

			suite.Store.Mutex.RLock()
			result, exists := suite.Store.URLToShort[original]
			suite.Store.Mutex.RUnlock()

			assert.True(suite.T(), exists, "URL should exist in URLToShort")
			assert.Equal(suite.T(), short, result, "Short URL should match")
		}(i)
	}

	wg.Wait()
	assert.Len(suite.T(), suite.Store.URLToShort, numGoroutines, "URLToShort should contain all entries")
}

func (suite *StoreTestSuite) TestConcurrentReadWriteShortToURL() {
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			original := "https://example.com/" + string(rune(i))
			short := "short" + string(rune(i))

			suite.Store.Mutex.Lock()
			suite.Store.ShortToURL[short] = original
			suite.Store.Mutex.Unlock()

			suite.Store.Mutex.RLock()
			result, exists := suite.Store.ShortToURL[short]
			suite.Store.Mutex.RUnlock()

			assert.True(suite.T(), exists, "Short URL should exist in ShortToURL")
			assert.Equal(suite.T(), original, result, "Original URL should match")
		}(i)
	}

	wg.Wait()
	assert.Len(suite.T(), suite.Store.ShortToURL, numGoroutines, "ShortToURL should contain all entries")
}

func (suite *StoreTestSuite) TestConcurrentReadWriteDomainHits() {
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			domain := "example.com"

			suite.Store.Mutex.Lock()
			suite.Store.DomainHits[domain]++
			suite.Store.Mutex.Unlock()

			suite.Store.Mutex.RLock()
			count := suite.Store.DomainHits[domain]
			suite.Store.Mutex.RUnlock()

			assert.GreaterOrEqual(suite.T(), count, 1, "Domain hit count should be at least 1")
		}(i)
	}

	wg.Wait()
	assert.Equal(suite.T(), numGoroutines, suite.Store.DomainHits["example.com"], "DomainHits should reflect all increments")
}

func (suite *StoreTestSuite) TestConcurrentReadEmptyMaps() {
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			suite.Store.Mutex.RLock()
			urlToShortLen := len(suite.Store.URLToShort)
			shortToURLLen := len(suite.Store.ShortToURL)
			domainHitsLen := len(suite.Store.DomainHits)
			suite.Store.Mutex.RUnlock()

			assert.Equal(suite.T(), 0, urlToShortLen, "URLToShort should remain empty")
			assert.Equal(suite.T(), 0, shortToURLLen, "ShortToURL should remain empty")
			assert.Equal(suite.T(), 0, domainHitsLen, "DomainHits should remain empty")
		}()
	}

	wg.Wait()
}
