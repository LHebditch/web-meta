package webmeta

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockClient struct {
	Response []byte
	Error error
}

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewReader(m.Response)),
	}, m.Error
}

func TestHandlerGetsMeta(t *testing.T) {
	client := &MockClient{
		Response: []byte(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Test Title</title>
					<meta property="og:foo" content="bar" />
					<meta property="og:title" content="test og title" />
					<meta name="foo" content="bar" />
					<meta name="title" content="test meta title" />
				</head>
			</html>
		`),
		Error: nil,
	}
	log, _ := zap.NewDevelopment()
	meta, err := GetWebMeta(log, "https://test.com", client)

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "bar", meta.OG["foo"])
	assert.Equal(t, "test og title", meta.OG["title"])
	assert.Equal(t, "bar", meta.Meta["foo"])
	assert.Equal(t, "test meta title", meta.Meta["title"])
	assert.Equal(t, "Test Title", meta.Title)
}

func TestHandlerReturnsErrorFromURL(t *testing.T) {}

func TestHandlerTrimsTabsandNewLines(t *testing.T) {}

func TestExtractMetaInto(t *testing.T) {}

func TestTrimAll(t *testing.T) {}

func TestCreateRequest(t *testing.T) {}
