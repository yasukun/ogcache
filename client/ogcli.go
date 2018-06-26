package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"ogcache"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// newProtocolFactory ...
func newProtocolFactory(proto string) (protocolFactory thrift.TProtocolFactory, err error) {
	switch proto {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		err = fmt.Errorf("Invalid protocol specified %v", proto)
	}
	return
}

// newTransportFactory ...
func newTransportFactory(buffered, framed bool) (transportFactory thrift.TTransportFactory) {
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}
	return
}

// RunClient ...
func RunClient(addr, proto string, buffered, framed, secure bool, url string) (og *ogcache.OpenGraph, err error) {
	var transport thrift.TTransport
	protocolFactory, err := newProtocolFactory(proto)
	if err != nil {
		return
	}
	transportFactory := newTransportFactory(buffered, framed)
	if secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTSocket(addr)
	}
	if err != nil {
		err = fmt.Errorf("Error opening socket: %v", err)
		return
	}
	if transport == nil {
		err = fmt.Errorf("Error opening socket, got nil transport. Is server available?")
		return
	}
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return
	}

	err = transport.Open()
	if err != nil {
		return
	}
	defer transport.Close()
	return handleClient(ogcache.NewOgServiceClientFactory(transport, protocolFactory), url)
}

func handleClient(client *ogcache.OgServiceClient, url string) (og *ogcache.OpenGraph, err error) {
	ctx := context.Background()
	og, err = client.Inquiry(ctx, url)
	return
}
