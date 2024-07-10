package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	SURAH_LAST  = 114
	SURAH_FIRST = 1
)

type QG struct {
	url       string
	token     string
	lgCodes   map[string]string
	cacheSize int
	cache     []CacheItem
}

type CacheItem struct {
	key  CacheKey
	json map[string]interface{}
}

type CacheKey struct {
	surah int
	ayah  int
	lang  string
}

func NewQG(url, token string, lgCodes map[string]string, cacheSize int) *QG {
	return &QG{
		url:       url,
		token:     token,
		lgCodes:   lgCodes,
		cacheSize: cacheSize,
		cache:     make([]CacheItem, 0),
	}
}

func (q *QG) updateCache() {
	if q.cacheSize == 0 {
		return
	}
	if len(q.cache) > q.cacheSize {
		q.cache = q.cache[:q.cacheSize]
	}
}

func (q *QG) getAyah(surah, ayah int, lang string) (map[string]interface{}, error) {
	if len(lang) == 2 {
		fullLang, ok := q.lgCodes[lang]
		if !ok {
			return nil, errors.New(lang + " is not supported using 2 letter codes")
		}
		lang = fullLang
	}
	if surah < SURAH_FIRST || surah > SURAH_LAST {
		return nil, errors.New("surah(chapter) must be between " + strconv.Itoa(SURAH_FIRST) + " and " + strconv.Itoa(SURAH_LAST))
	}

	for i, item := range q.cache {
		if item.key.surah == surah && item.key.ayah == ayah && item.key.lang == lang {
			q.cache = append(q.cache[:i], q.cache[i+1:]...)
			q.cache = append([]CacheItem{item}, q.cache...)
			return item.json, nil
		}
	}

	reqURL := fmt.Sprintf("%s%d/%s", q.url, surah, lang)
	resp, err := http.Get(reqURL + "?key=" + q.token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ayahJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&ayahJSON)
	if err != nil {
		return nil, err
	}

	q.cache = append([]CacheItem{{key: CacheKey{surah: surah, ayah: ayah, lang: lang}, json: ayahJSON}}, q.cache...)
	q.updateCache()

	return ayahJSON, nil
}
