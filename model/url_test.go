package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

func (suite *ModelTestSuite) TestURLRequestMarshalJSON() {
	req := URLRequest{OriginalURL: "https://example.com"}
	expectedJSON := `{"original_url":"https://example.com"}`

	data, err := json.Marshal(req)

	assert.NoError(suite.T(), err, "Marshalling URLRequest should not error")
	assert.JSONEq(suite.T(), expectedJSON, string(data), "JSON output should match expected")
}

func (suite *ModelTestSuite) TestURLRequestUnmarshalJSON() {
	jsonStr := `{"original_url":"https://example.com"}`
	expected := URLRequest{OriginalURL: "https://example.com"}

	var req URLRequest
	err := json.Unmarshal([]byte(jsonStr), &req)

	assert.NoError(suite.T(), err, "Unmarshalling valid JSON should not error")
	assert.Equal(suite.T(), expected, req, "Unmarshalled URLRequest should match expected")
}

func (suite *ModelTestSuite) TestURLRequestUnmarshalInvalidJSON() {
	jsonStr := `{invalid json}`

	var req URLRequest
	err := json.Unmarshal([]byte(jsonStr), &req)

	assert.Error(suite.T(), err, "Unmarshalling invalid JSON should error")
	assert.Contains(suite.T(), err.Error(), "invalid character", "Error should indicate invalid JSON")
}

func (suite *ModelTestSuite) TestURLRequestUnmarshalEmptyURL() {
	jsonStr := `{"original_url":""}`
	expected := URLRequest{OriginalURL: ""}

	var req URLRequest
	err := json.Unmarshal([]byte(jsonStr), &req)

	assert.NoError(suite.T(), err, "Unmarshalling empty URL should not error")
	assert.Equal(suite.T(), expected, req, "Unmarshalled URLRequest should have empty OriginalURL")
}

func (suite *ModelTestSuite) TestURLResponseMarshalJSON() {
	resp := URLResponse{ShortURL: "abc123"}
	expectedJSON := `{"short_url":"abc123"}`

	data, err := json.Marshal(resp)

	assert.NoError(suite.T(), err, "Marshalling URLResponse should not error")
	assert.JSONEq(suite.T(), expectedJSON, string(data), "JSON output should match expected")
}

func (suite *ModelTestSuite) TestURLResponseUnmarshalJSON() {
	jsonStr := `{"short_url":"abc123"}`
	expected := URLResponse{ShortURL: "abc123"}

	var resp URLResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)

	assert.NoError(suite.T(), err, "Unmarshalling valid JSON should not error")
	assert.Equal(suite.T(), expected, resp, "Unmarshalled URLResponse should match expected")
}

func (suite *ModelTestSuite) TestURLResponseUnmarshalInvalidJSON() {
	jsonStr := `{invalid json}`

	var resp URLResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)

	assert.Error(suite.T(), err, "Unmarshalling invalid JSON should error")
	assert.Contains(suite.T(), err.Error(), "invalid character", "Error should indicate invalid JSON")
}

func (suite *ModelTestSuite) TestURLResponseUnmarshalEmptyShortURL() {
	jsonStr := `{"short_url":""}`
	expected := URLResponse{ShortURL: ""}

	var resp URLResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)

	assert.NoError(suite.T(), err, "Unmarshalling empty short URL should not error")
	assert.Equal(suite.T(), expected, resp, "Unmarshalled URLResponse should have empty ShortURL")
}
