package webmeta

import (
	"bytes"
	"errors"
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

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewReader(m.Response)),
	}, m.Error
}

const mockHtml string = `
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
`

func TestHandlerGetsMeta(t *testing.T) {
	client := &MockClient{
		Response: []byte(mockHtml),
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

func TestHandlerReturnsErrorFromURL(t *testing.T) {
	expectedError := errors.New("Mock Error")
	client := &MockClient{
		Response: []byte(mockHtml),
		Error: expectedError,
	}
	log, _ := zap.NewDevelopment()
	meta, err := GetWebMeta(log, "https://test.com", client)
	assert.Equal(t, err, expectedError)
	assert.Equal(t, meta.Title, "")
}

func TestHandlerTrimsTabsandNewLines(t *testing.T) {
	client := &MockClient{
		Response: []byte(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Test Title</title>
				<meta property="og:foo" content="bar" />
				<meta property="og:title" content="test og title
				" />
				<meta name="foo" content="bar" />
				<meta name="title" content="test meta title	" />
			</head>
		</html>
		`),
		Error: nil,
	}
	log, _ := zap.NewDevelopment()
	meta, _ := GetWebMeta(log, "https://test.com", client)
	assert.Equal(t, meta.Meta["title"], "test meta title")
	assert.Equal(t, meta.OG["title"], "test og title")
}

func TestTrimAll(t *testing.T) {
	withTab := "some text\t"
	withNewline := "some text\n"
	withReturn := "some text\r"

	withAll := "some\n text\t\r"

	expected := "some text"
	assert.Equal(t, TrimAll(withTab), expected)
	assert.Equal(t, TrimAll(withNewline), expected)
	assert.Equal(t, TrimAll(withReturn), expected)
	assert.Equal(t, TrimAll(withAll), expected)
}
