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
	"crypto/tls"
	"fmt"
	"ogcache"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/yasukun/ogcache-server/lib"
)

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, conf lib.Config, db *ledis.DB) error {
	var transport thrift.TServerTransport
	var err error
	if conf.Main.Secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("keys/server.crt", "keys/server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(conf.Main.Addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(conf.Main.Addr)
	}

	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	handler := NewOgcacheHandler(conf, db)
	processor := ogcache.NewOgServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the ogcache server... on ", conf.Main.Addr)
	return server.Serve()
}
