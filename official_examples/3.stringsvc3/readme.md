stringsvc3 演示的go-kit服务如何调用另一个go-kit服务，官方推荐使用proxyMiddleware

- go-kit的中间件模式在这里不太好用，go-kit中所有的中间件都是需要实现Service的所有接口的，它与其他框架中的
中间件不太一样，在这里使用起来比较麻烦

- proxy在这里指的是A服务的接口a实际调用的是B服务的b接口，那么就得写一个proxy中间件代理a接口，同时你还需要实现包含a在内的所有接口，
你才能安装这个中间件


