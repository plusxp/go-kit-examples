# addsvc

addsvc is an example microservice which takes full advantage of most of Go
kit's features, including both service- and transport-level middlewares,
speaking multiple transports simultaneously, distributed tracing, and rich
error definitions. The server binary is available in cmd/addsvc. The client
binary is available in cmd/addcli.

Finally, the addtransport package provides both server and clients for each
supported transport. The client structs bake-in certain middlewares, in order to
demonstrate the _client library pattern_. But beware: client libraries are
generally a bad idea, because they easily lead to the
 [distributed monolith antipattern](https://www.microservices.com/talks/dont-build-a-distributed-monolith/).
If you don't _know_ you need to use one in your organization, it's probably best
avoided: prefer moving that logic to consumers, and relying on 
 [contract testing](https://docs.pact.io/best_practices/contract_tests_not_functional_tests.html)
to detect incompatibilities.


## 关于mw(middleware)的使用说明

addsvc这个示例中，service和endpoint都添加了mw，但你注意到，service的mw的添加尤其麻烦，安装一个log mw的写法是：
```go
svc = NewBasicService()
svc = LoggingMiddleware(logger)(svc)
```

可以看到，LoggingMiddleware() 接受和传入的都是svc，通过查看代码可以发现，示例又创建了一个type：loggingMiddleware，这个type完整的实现
了svc的每个接口，mw的功能在每个新实现的接口中得以体现；

当我们要添加一个新的监控mw时，如果要和log mw硬隔离，就只能再创建一个type：InstrumentingMiddleware，然后完整实现svc的每个接口，才可以安装到svc。

#### endpoint的mw

通过代码但我们发现，endpoint的mw编写方式比较简单，与server的mw截然不同，它只需要写一个mw方法即可安装。

然后我们可以思考一下两个问题：

-   service层是否有必要添加mw？
-   与endpoint层的mw有和不同？
-   如果有必要，是否需要将不同作用的mw代码硬隔离？

第一个问题，service层是否有必要添加mw？  

**我的回答是**：看需求，我们先看一下service的log mw的多个实现接口的其中一个的写法：

```go
// Sum是svc的一个接口，作用是return sum(a,b)
func (mw loggingMiddleware) Sum(ctx context.Context, a, b int) (v int, err error) {
	defer func() {
		mw.logger.Log("method", "Sum", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}
```

注意观察，因为是完整的svc的接口的另一实现，所以在这一层可以拿到接口request的每个参数和最终响应（struct），所以我们就可以选择想要操作的参数进行log，或者监控。

那可不可以在endpoint层实现呢？ 看一下endpoint层的mw实现（只需要写一个方法）：
```go
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}
```
显然，在这一层我们只能拿到一个interface{}的request和response，如果要log，我们只能log整个对象，同样的，监控也是；
这里体现了有得必有失。。

第三个问题，是否需要将不同作用的mw代码硬隔离？

我的想法，没有必要，这是在开发效率、代码耦合两方面权衡后的一个结果，即使将log和监控的代码都放在一个mw对象中，我们也不会丢失多少代码简洁性，
看看放在一个mw对象中的代码：
```go
// 统一mw
func UnifyMiddleware(loggermw log.Logger, ints, chars metrics.Counter) Middleware {
	return func(next Service) Service {
		instrumw := instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
		return unifyMiddleware{loggermw, instrumw,next}
	}
}

// 监控mw
type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  Service
}

// 统一mw实体
type unifyMiddleware struct {
	// 通过命名规范代码， 不同功能的对象通过嵌套struct添加
	loggermw log.Logger
	instrumw instrumentingMiddleware
	next     Service
}

func (mw unifyMiddleware) Sum(ctx context.Context, a, b int) (v int, err error) {
	defer func() {
		mw.loggermw.Log("method", "Sum", "a", a, "b", b, "v", v, "err", err)
	}()
	v, err = mw.next.Sum(ctx, a, b)
	mw.instrumw.ints.Add(float64(v))
	return v, err
}
```

是不是也足够简洁？ 虽然，这需要的更多是编码人员良好的编码习惯。