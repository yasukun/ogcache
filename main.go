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
	"flag"
	"fmt"
	"log"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/yasukun/ogcache-server/lib"
)

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func main() {
	flag.Usage = Usage
	confname := flag.String("c", "ogcache-server.toml", "path to config")
	flag.Parse()
	_, err := os.Stat(*confname)
	if err != nil {
		log.Fatalln(err)
	}

	conf, err := lib.DecodeConfigToml(*confname)
	if err != nil {
		log.Fatalln(err)
	}

	var protocolFactory thrift.TProtocolFactory
	switch conf.Main.Protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", conf.Main.Protocol, "\n")
		Usage()
		os.Exit(1)
	}
	var transportFactory thrift.TTransportFactory
	if conf.Main.Bufferd {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if conf.Main.Framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	// always run server here
	if err := runServer(transportFactory, protocolFactory, conf); err != nil {
		fmt.Println("error running server:", err)
	}
}
