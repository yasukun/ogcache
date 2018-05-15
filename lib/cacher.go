package lib

/*
 * MIT License
 *
 * Copyright (c) 2018 yasukun
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import (
	"io/ioutil"
	"net/http"
	"ogcache"
	"strings"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/siddontang/ledisdb/ledis"
)

// GetHtml ...
func GetHtml(headers map[string]string, url string) ([]byte, error) {
	var byteHtml []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return byteHtml, err
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return byteHtml, err
	}
	defer resp.Body.Close()
	byteHtml, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return byteHtml, err
	}
	return byteHtml, nil
}

// GetOpenGraph ...
func GetOpenGraph(html []byte) (*opengraph.OpenGraph, error) {
	og := opengraph.NewOpenGraph()
	err := og.ProcessHTML(strings.NewReader(string(html)))
	if err != nil {
		return og, err
	}
	return og, nil
}

// ExistsOpenGraph ...
func ExistsOpenGraph(db *ledis.DB, url string) bool {
	i64, err := db.Exists([]byte(url))
	if err != nil {
		return false
	}
	if i64 > 0 {
		return true
	}
	return false
}

// GetOpenGraphFromCache ...
func GetOpenGraphFromCache(db *ledis.DB, url string) ([]byte, error) {
	var og []byte
	og, err := db.Get([]byte(url))
	if err != nil {
		return og, err
	}
	return og, nil
}

// SetOpenGraphToCache ...
func SetOpenGraphToCache(db *ledis.DB, url, ogJson string) error {
	err := db.Set([]byte(url), []byte(ogJson))
	if err != nil {
		return err
	}
	return nil
}

// ConvOpenGraph ...
func ConvOpenGraph(og *opengraph.OpenGraph) *ogcache.OpenGraph {

	image_url := ""
	for _, v := range og.Images {
		image_url = v.URL
		break
	}
	audio_url := ""
	for _, v := range og.Audios {
		audio_url = v.URL
		break
	}
	video_url := ""
	for _, v := range og.Videos {
		video_url = v.URL
		break
	}

	localeAlternate := []string{}
	for _, v := range og.LocalesAlternate {
		localeAlternate = append(localeAlternate, v)
	}

	return &ogcache.OpenGraph{
		Title:           og.Title,
		Type:            og.Type,
		Image:           image_url,
		URL:             og.URL,
		Audio:           audio_url,
		Description:     og.Description,
		Determiner:      og.Determiner,
		Locale:          og.Locale,
		LocaleAlternate: localeAlternate,
		SiteName:        og.SiteName,
		Video:           video_url,
	}
}
