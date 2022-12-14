package webmeta

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/LHebditch/htmlmeta/models"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func TrimAll(s string) string {
	reNewLine := regexp.MustCompile(`\r?\n`)
	reTb := regexp.MustCompile(`\t`)
	s = reNewLine.ReplaceAllString(s, "")
	s = reTb.ReplaceAllLiteralString(s, "")
	trimmed := strings.TrimSpace(s)
	return trimmed
}

func extractMetaInto(meta *models.HtmlMeta) (extractor models.HtmlExtractor) {
	return func(_ int, s *goquery.Selection) {
		prop := s.AttrOr("property", "")
		name := s.AttrOr("name", "")
		content := TrimAll(s.AttrOr("content", ""))

		if prop != "" && strings.Contains(prop, "og:") {
			meta.AddOG(strings.ReplaceAll(prop, "og:", ""), content)
		} else if name != "" {
			meta.AddMeta(name, content)
		}
	}
}

func createRequest(log *zap.Logger, url string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", zap.Error(err))
		return
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	req.Header.Set("accept", "text/html")
	req.Header.Set("cache-control", "max-ag=0")
	req.Header.Set("charset", "utf-8")

	return
}

func GetWebMeta(log *zap.Logger, url string, client HTTPClient) (meta models.HtmlMeta, err error) {
	log.With(zap.String("url", url))
	log.Info("Processing request for web meta")
	req, err := createRequest(log, url)
	if err != nil {
		log.Error("failed to create request", zap.Error(err))
		return
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send request", zap.Error(err))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Failed to conect to url, responded with " + res.Status, zap.Error(err))
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Failed to create go query doc from reponse", zap.Error(err))
		return
	}
	meta = models.HtmlMeta{}
	doc.Find("meta").Each(extractMetaInto(&meta))
	meta.Title = TrimAll(doc.Find("head > title").First().Text())
	log.Info("Processed web meta request")
	return
}
