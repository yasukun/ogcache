# ogcache-server

RPC Server to get OpenGraph metadata
used Apache Thrift

## build

```bash
$ dep ensure
$ go build
```

## test

```bash
$ ./ogcache-server
```

```bash
$ pip install thriftpy
$ python client-test.py
```