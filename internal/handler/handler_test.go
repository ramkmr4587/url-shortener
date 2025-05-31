package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/service"
	"url-shortener/model"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	urlShortenFail       = false
	urlGetOriginalFail   = false
	urlGetTopDomainsFail = false
)

type HandlerTestSuite struct {
	suite.Suite
	Handler          *Handler
	Webservice       *restful.WebService
	Container        *restful.Container
	ResponseRecorder *httptest.ResponseRecorder
	RestResponse     *restful.Response
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	mockService := &urlServiceMock{}
	suite.Handler = NewHandler((*service.URLService)(mockService))
	suite.Container = restful.NewContainer()
	suite.Webservice = new(restful.WebService).
		Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	suite.ResponseRecorder = httptest.NewRecorder()
	suite.RestResponse = &restful.Response{ResponseWriter: suite.ResponseRecorder}
	urlShortenFail = false
	urlGetOriginalFail = false
	urlGetTopDomainsFail = false

	suite.Webservice.Route(suite.Webservice.POST("/shorten").To(suite.Handler.Shorten))
	suite.Webservice.Route(suite.Webservice.GET("/r/{short}").To(suite.Handler.Redirect))
	suite.Webservice.Route(suite.Webservice.GET("/metrics").To(suite.Handler.Metrics))
	suite.Container.Add(suite.Webservice)
}

func (suite *HandlerTestSuite) TestShortenSuccess() {
	reqBody := model.URLRequest{OriginalURL: "https://example.com"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", restful.MIME_JSON)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	response := convertToURLResponse(suite.ResponseRecorder.Body.String())
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusOK, suite.ResponseRecorder.Result().StatusCode, "Expected 200 OK, got %d", suite.ResponseRecorder.Result().StatusCode)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "abc123", response.ShortURL)
}

func (suite *HandlerTestSuite) TestShortenParseError() {
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", restful.MIME_JSON)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, suite.ResponseRecorder.Result().StatusCode)
	assert.Contains(suite.T(), suite.ResponseRecorder.Body.String(), "invalid character")
}

func (suite *HandlerTestSuite) TestShortenError() {
	urlShortenFail = true
	reqBody := model.URLRequest{OriginalURL: "https://fail.com"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", restful.MIME_JSON)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusInternalServerError, suite.ResponseRecorder.Result().StatusCode)
	assert.Contains(suite.T(), suite.ResponseRecorder.Body.String(), "expected shorten to fail")
}

func (suite *HandlerTestSuite) TestRedirectSuccess() {
	req := httptest.NewRequest("GET", "/r/abc123", nil)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusMovedPermanently, suite.ResponseRecorder.Result().StatusCode)
	assert.Equal(suite.T(), "https://example.com", suite.ResponseRecorder.Header().Get("Location"))
}

func (suite *HandlerTestSuite) TestRedirectNotFound() {
	req := httptest.NewRequest("GET", "/r/invalid", nil)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusNotFound, suite.ResponseRecorder.Result().StatusCode)
	assert.Contains(suite.T(), suite.ResponseRecorder.Body.String(), "short URL not found")
}

func (suite *HandlerTestSuite) TestRedirectError() {
	urlGetOriginalFail = true
	req := httptest.NewRequest("GET", "/r/fail", nil)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusInternalServerError, suite.ResponseRecorder.Result().StatusCode)
	assert.Contains(suite.T(), suite.ResponseRecorder.Body.String(), "expected get original to fail")
}

func (suite *HandlerTestSuite) TestMetricsSuccess() {
	req := httptest.NewRequest("GET", "/metrics", nil)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	var response []struct {
		Domain string `json:"domain"`
		Count  int    `json:"count"`
	}
	err := json.Unmarshal(suite.ResponseRecorder.Body.Bytes(), &response)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, suite.ResponseRecorder.Result().StatusCode)
	assert.Len(suite.T(), response, 3)
	assert.Equal(suite.T(), "example.com", response[0].Domain)
	assert.Equal(suite.T(), 10, response[0].Count)
	assert.Equal(suite.T(), "test.com", response[1].Domain)
	assert.Equal(suite.T(), 5, response[1].Count)
	assert.Equal(suite.T(), "other.com", response[2].Domain)
	assert.Equal(suite.T(), 3, response[2].Count)
}

func (suite *HandlerTestSuite) TestMetricsError() {
	urlGetTopDomainsFail = true
	req := httptest.NewRequest("GET", "/metrics", nil)

	suite.Container.ServeHTTP(suite.ResponseRecorder, req)
	suite.T().Logf("Response body: %s", suite.ResponseRecorder.Body.String())
	assert.Equal(suite.T(), http.StatusInternalServerError, suite.ResponseRecorder.Result().StatusCode)
	assert.Contains(suite.T(), suite.ResponseRecorder.Body.String(), "expected get top domains to fail")
}

func convertToURLResponse(str string) *model.URLResponse {
	response := &model.URLResponse{}
	err := json.Unmarshal([]byte(str), response)
	if err != nil {
		return nil
	}
	return response
}

type urlServiceMock service.URLService

func (mock *urlServiceMock) ShortenURL(original string) string {
	if urlShortenFail {
		panic(errors.New("expected shorten to fail"))
	}
	if original == "" {
		return ""
	}
	return "abc123"
}

func (mock *urlServiceMock) GetOriginalURL(short string) (string, bool) {
	if urlGetOriginalFail {
		panic(errors.New("expected get original to fail"))
	}
	if short == "invalid" {
		return "", false
	}
	return "https://example.com", true
}

func (mock *urlServiceMock) GetTopDomains(limit int) map[string]int {
	if urlGetTopDomainsFail {
		panic(errors.New("expected get top domains to fail"))
	}
	return map[string]int{
		"example.com": 10,
		"test.com":    5,
		"other.com":   3,
		"more.com":    1,
	}
}
