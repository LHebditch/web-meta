package models

import (
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
)

type HtmlExtractor func(i int, s *goquery.Selection)

type HtmlMeta struct {
	Meta  map[string]string `json:"meta_tags"`
	OG    map[string]string `json:"og_tags"`
	Title string            `json:"title"`
}

func (h *HtmlMeta) AddMeta(prop, content string) {
	if h.Meta == nil {
		h.Meta = make(map[string]string)
	}
	h.Meta[prop] = content
}

func (h *HtmlMeta) AddOG(prop, content string) {
	if h.OG == nil {
		h.OG = make(map[string]string)
	}
	h.OG[prop] = content
}

func (h HtmlMeta) ToString() (string, error) {
	str, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	return string(str), nil
}