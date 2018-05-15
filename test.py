import thriftpy
ogcache_thrift = thriftpy.load("service.thrift", module_name="ogcache_thrift")

from thriftpy.rpc import make_client

client = make_client(ogcache_thrift.OgService, '127.0.0.1', 9090)
og = client.inquiry("")
print(og)
