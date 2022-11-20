package webmeta

import (
	"net/http"
	"strings"

	"github.com.LHebditch.htmlmeta/models"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func ExtractMetaInto(meta *models.HtmlMeta) (extractor models.HtmlExtractor) {
	return func(i int, s *goquery.Selection) {
		prop := s.AttrOr("property", "")
		name := s.AttrOr("name", "")
		content := s.AttrOr("content", "")

		if prop != "" && strings.Contains(prop, "og:") {
			meta.AddOG(prop, content)
		} else if name != "" {
			meta.AddMeta(name, content)
		}
	}
}

func CreateRequest(log *zap.Logger, url string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", zap.Error(err))
		return
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Set("accept", "text/html")
	req.Header.Set("cache-control", "max-ag=0")
	req.Header.Set("charset", "utf-8")

	return
}

func GetWebMeta(log *zap.Logger, url string) (meta models.HtmlMeta, err error) {
	log.With(zap.String("url", url))
	log.Info("Processing request for web meta")
	client := &http.Client{}
	req, err := CreateRequest(log, url)
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
	doc.Find("meta").Each(ExtractMetaInto(&meta))
	meta.Title = doc.Find("head > title").First().Text()
	log.Info("Processed web meta request")
	return
}
