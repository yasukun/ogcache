package main

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
	"context"
	"encoding/json"
	"ogcache"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/yasukun/ogcache-server/lib"
)

type OgcacheHandler struct {
	DB      *ledis.DB
	Headers map[string]string
}

// NewOgcacheHandler ...
func NewOgcacheHandler(conf lib.Config, db *ledis.DB) *OgcacheHandler {
	headers := map[string]string{}
	for _, v := range conf.Headers {
		headers[v.Name] = v.Value
	}
	return &OgcacheHandler{Headers: headers, DB: db}
}

// (o *OgcacheHandler) Inquiry ...
func (o *OgcacheHandler) Inquiry(ctx context.Context, url string) (*ogcache.OpenGraph, error) {
	var og *ogcache.OpenGraph
	if !lib.ExistsOpenGraph(o.DB, url) {
		byteHtml, err := lib.GetHtml(o.Headers, url)
		if err != nil {
			return og, err
		}
		openGraph, err := lib.GetOpenGraph(byteHtml)
		if err != nil {
			return og, err
		}
		if err = lib.SetOpenGraphToCache(o.DB, url, openGraph.String()); err != nil {
			return og, err
		}
		og = lib.ConvOpenGraph(openGraph)
	} else {
		jsonRaw, err := lib.GetOpenGraphFromCache(o.DB, url)
		if err != nil {
			return og, err
		}
		openGraph := new(opengraph.OpenGraph)
		if err = json.Unmarshal(jsonRaw, openGraph); err != nil {
			return og, err
		}

		og = lib.ConvOpenGraph(openGraph)
	}

	return og, nil
}
