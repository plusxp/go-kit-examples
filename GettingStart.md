
# :zap: 使用代码生成工具快速开发Go-kit微服务！


## 前言
本文档将试图以最简明的方式讲述如何在实际环境中使用go-kit快速编写一个新的微服务项目，
如有逻辑纰漏或表述错误，请不吝赐教，谢谢。

<br/>


请注意，阅读此文档默认读者具备以下条件：

- 一年以上的go语言开发经验
- 一定的微服务理论基础
- 了解gRPC及proto协议
- 了解go-kit框架分层思想(transport, endpoint, service)
- 最好阅读过`demo_project/new_addsvc`项目，以便了解成型后的项目布局

---

读完官方示例或者`demo_project/new_addsvc`的项目之后，读者应该知道如果所有的代码都需要自己编写，那使用go-kit开发
微服务的效率将是难以接受的；但是在仔细观察该服务的具体代码后，我们可以清楚了解到，有几个层次的代码
都是根据service层的实现而编写的，它们的代码**模式**非常的有迹可循，有经验的go开发者很快可以想到，
为什么这些层次的代码不能使用工具生成呢？是的，早已有人想到这一点，它们开发出了go-kit的CLI工具，有了这个
工具，我们使用go-kit开发那就是如虎添翼，下面一起来看看。

我调研了[go-kit](https://github.com/go-kit/kit) 官方文档中列出的几个代码生成工具，找到了一个比较适合
使用的repo，那就是[kit](https://github.com/GrantZheng/kit) ，其他repo要么比较简陋，要么
正在开发，或者不再维护；这个repo其实也是fork一个不再维护的400+star的项目而来，该作者声称自己所在团队已深度使用
go-kit，无法接受没有可靠的go-kit辅助工具，所以自己fork来继续维护了（点赞！）  
> 注：在使用过程我发现该工具仍不够灵活以及缺乏一些功能，目前我已fork此项目，并增加了诸多功能，请查看 [chaseSpace/kit][kit] ，
下文也是基于此仓库编写。

## 目录
-   [1. 关于kit](#1-关于kit)
-   [2. 开始](#2-开始)
-   [3. 创建Project](#3-创建Project)
-   [4. 生成Service模板](#4-生成Service模板)
-   [5. 编辑proto文件](#5-编辑proto文件)
-   [6. 实现Service接口](#6-实现Service接口)
-   [7. 需要完善的工作](#7-需要完善的工作)
-   [8. 启动server](#8-启动server)

___
-   [9. 生成Client-side代码](#9-生成Client-side代码)
-   [10. 塑造适合你(的团队)的Client](#10-塑造适合你(的团队)的Client)
-   [11. Let's test it now](#11-Let's-test-it-now)

___
-   [12. 自由尚在](#12-自由尚在)
-   [13. 结束，新的开始](#13-结束新的开始)


## 1. 关于kit
我们需要知道它的一些功能、特点

-   可以生成指定名称的service模板代码，包含endpoint、transport(http/grpc)层
-   可以生成client代码
-   生成的以`_gen.go`结尾的文件不可修改，在重新生成时会覆盖
-   不会删除/覆盖已有的任何代码(除了`*_gen.go`)

> (以下部分描述译自GrantZheng/kit README.md, 少部分改动以适配本仓库，并带有额外的说明)  

## 2. 开始

首先安装go-kit CLI工具
```bash
# 默认使用modules包管理
# 注意：这个仓库是fork而来，go.mod的module name是原仓库名，不能通过go get方式下载
$ mkdir $GOPATH/pkg/mod/git_repo -p
$ cd $GOPATH/pkg/mod/git_repo
$ git clone https://github.com/chaseSpace/kit.git
$ cd kit
$ go install 

# check usage
$ kit help
```

## 3. 创建Service

```bash
$ kit new service hello       # 缩写： kit n s hello
# 若要指定生成的go.mod内的模块名，指令后加上 --module module_hello，缩写-m，默认使用service名作为模块名
```

这条命令将会在当前目录下创建一个hello目录，结构如下

```DOS
c:\users\...\go-kit-examples\demo_project\hello
│  go.mod
│
└─pkg
    └─service
            service.go
```

## 4. 生成Service模板

在执行之前，首先需要在service.go中添加我们的api定义，示例：

```go
// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	SayHi(ctx context.Context, name string) (reply string, err error)
}
```

这里选择grpc作为RPC技术栈

```bash
$ kit g s hello
$ kit g s hello --dmw -t grpc  # -dmw生成endpoint中间件(可选)，-t grpc指定 transport (default http)
```

这些命令会执行以下操作：

- 创建service层代码模板: `hello/pkg/service/service.go`
- 创建service层middleware: `hello/pkg/service/middleware.go`
- 创建endpoint: `hello/pkg/endpoint/endpoint.go` and `hello/pkg/endpoint/endpoint_gen.go`
- `--dmw` 创建endpoint middleware: `hello/pkg/endpoint/middleware.go`
- 创建transport files e.g. http: `service-name/pkg/http/handler.go` 以及 `service-name/pkg/http/handler_gen.go`
- 若使用`-t grpc`，则创建grpc transport files: `service-name/pkg/grpc/handler.go` 以及 `service-name/pkg/grpc/handler_gen.go`
- 若使用`-t grpc`，则创建`pkg/grpc/pb/`目录存放proto文件，以及它的编译脚本
- 创建service main file  
`hello/cmd/service/service.go`  
`hello/cmd/service/service_gen.go`  
`hello/cmd/main.go`

最后得到的目录结构如下：

```DOS
c:\users\...\go-kit-examples\demo_project\hello
│  go.mod
│
├─cmd
│  │  main.go
│  │
│  └─service
│          service.go
│          service_gen.go
│
└─pkg
    ├─endpoint
    │      endpoint.go
    │      endpoint_gen.go
    │      middleware.go
    │
    ├─grpc
    │  │  handler.go
    │  │  handler_gen.go
    │  │
    │  └─pb
    │          compile.sh  <-------proto文件编译脚本，根据当前shell环境生成对应脚本
    │          hello.pb.go
    │          hello.proto
    │
    └─service
            middleware.go
            service.go
```

注意，kit工具在/pkg目录生成了grpc目录，并且将pb目录也放在其中，根据`go_project_template`项目布局，`/pb`目录
最好放在项目根目录，这样方便快速找到一个服务的pb文件，当然你可以有自己的布局；

现在把上面生成的grpc、endpoint目录都删除，重新生成它们：
```go
// 查看-t grpc后面可跟的选项
$ kit g s hello -t grpc --help
// -p 指定proto文件要放的位置，这里我放到hello/pb/proto/下；-i 指定代码中pb文件的import路径，gen-go/pb目录会自动创建
// 注意：这里我们使用hello/pb/gen-go/来存放所有的*.pb.go文件
$ cd demo_project/
$ mkdir -p hello/pb/proto  // 提前创建此目录
$ kit g s hello --dmw -t grpc -p hello/pb/proto -i hello/pb/gen-go/pb
```

现在的hello目录结构如下：
```go
c:\users\...\go-kit-examples\demo_project\hello
│  go.mod
│  go.sum
│
├─cmd
│  │  main.go
│  │
│  └─service
│          service.go
│          service_gen.go
│
├─pb
│  │
│  └─proto
│        compile.sh
|        hello.pb.go
│        hello.proto
│
└─pkg
    ├─endpoint
    │      endpoint.go
    │      endpoint_gen.go
    │
    ├─grpc
    │      handler.go
    │      handler_gen.go
    │
    └─service
            middleware.go
            service.go
```

## 5. 编辑脚本和proto文件
因为我们指定了`-i`修改了pb文件import path为`hello/pb/proto`, 所以还需要手动修改`compile.sh`：
```bash
// old
protoc hello.proto --go_out=plugins=grpc:.

// new
protoc hello.proto --go_out=plugins=grpc:../../../
```
`../../../`表示从脚本所在位置向外【三层】走到`hello/pb/proto`这个路径根目录`hello/`的父目录位置，
如果你熟悉pb协议，那应该不会有什么问题，如果你完全没接触过grpc，那请先按照文档操作


然后编辑pb/hello.proto文件，按如下修改：
```proto
message SayHiRequest {
 string name = 1;
}

message SayHiReply {
 string reply = 1;
}
```
再次运行刚才的kit命令生成pb代码：
```bash
$ kit g s hello --dmw -t grpc -p hello/pb/proto -i hello/pb/gen-go/pb
time="2020-09-26T08:59:19+08:00" level=info msg="exec>[sh -c /C/Users/.../go-kit-examples/hello/pb/proto/compile.sh]"
time="2020-09-26T08:59:19+08:00" level=warning msg="==============================================================="
time="2020-09-26T08:59:19+08:00" level=warning msg="The GRPC implementation is not finished you need to update your"
time="2020-09-26T08:59:19+08:00" level=warning msg=" service proto buffer and run the compile script."
time="2020-09-26T08:59:19+08:00" level=warning msg=---------------------------------------------------------------
time="2020-09-26T08:59:19+08:00" level=warning msg="You also need to implement the Encoders and Decoders!"
time="2020-09-26T08:59:19+08:00" level=warning msg="==============================================================="
```
如你所想，kit会执行已存在的`compile.sh`

看看生成的pb代码位置:
```bash
c:\users\...\go-kit-examples\demo_project\hello
│  go.mod
│
├─pb
│  ├─gen-go
│  │  └─pb
│  │          hello.pb.go
| ...忽略了一些文件
```

## 6. 实现Service接口
修改/pkg/service/service.go, 实现SayHi接口逻辑：
```go
func (b *basicHelloService) SayHi(ctx context.Context, name string) (reply string, err error) {
	return "Hi," + name, err
}
```

## 7. 需要完善的工作
打开`/pkg/grpc/handler.go`, 你看到`encode...`和`decode...`这样的函数了吗？
这里我们还需要完成两项工作：
- gRPC-layer的Req ---decode---> Endpoint-layer的Req
- gRPC-layer的Rsp <---encode--- Endpoint-layer的Rsp

像下面这样：
```go
func decodeSayHiRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SayHiRequest)
	return &endpoint.SayHiRequest{Name: req.Name}, nil
}

func encodeSayHiResponse(_ context.Context, r interface{}) (interface{}, error) {
	rsp := r.(*endpoint.SayHiResponse)
	return &pb.SayHiReply{Reply: rsp.Reply}, nil
}
```
:rotating_light: 注意：这是容易出错的地方，因为编码时没有任何提示帮助我们填写正确的类型，同时我们也不应该
使用_,ok的方式来避免panic，出现类型错误一定是编码bug，不应该hide it。  
>当然，为避免程序退出，我们应该有额外的panic捕获措施，如grpc的recovery中间件

## 8. 启动server
OK，现在可以启动这个服务了
```go
cd hello/cmd
go run .
/* OUTPUT
ts=2020-09-12T12:36:00.620891Z caller=service.go:85 tracer=none
ts=2020-09-12T12:36:00.6258776Z caller=service.go:143 transport=debug/HTTP addr=:8080
ts=2020-09-12T12:36:00.6258776Z caller=service.go:107 transport=gRPC addr=:8082
*/
```

然后来简单看一下cmd目录下的代码，main.go就不用看了，它调用了`cmd/service/service.go`的Run()，
所以我们直接看后者代码, 下面是部分代码片段:
```go
var fs = flag.NewFlagSet("hello", flag.ExitOnError)
var debugAddr = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
var httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
var grpcAddr = fs.String("grpc-addr", ":8082", "gRPC listen address")
var thriftAddr = fs.String("thrift-addr", ":8083", "Thrift listen address")
```

这里提供了http,grpc,thrift三个RPC地址变量作为启动参数，但其实我们只用到了grpc，如果你用goland你会看到除了grpc
其他几个变量都是有提示‘unused variable’，所以这里可以直接删掉这几行代码，以免扰乱视线。

然后再提一下Run()方法的最后几行代码：
```go
svc := service.New(getServiceMiddleware(logger))
eps := endpoint.New(svc, getEndpointMiddleware(logger))
g := createService(eps)
initMetricsEndpoint(g)
initCancelInterrupt(g)
logger.Log("exit", g.Run())
```

这里的代码非常简洁明了，前三行就是以洋葱模式一层层的封装svc对象（各层都可应用中间件），`initMetricsEndpoint`就是
启动指标http服务供prometheus来拉svc的监控数据，`initCancelInterrupt`就是监听信号，优雅退出服务。
> 关于`cmd/service/service.go`中的部分编码模式，我们可以修改为自己惯用的方式，kit命令不会覆盖已有的`service.go`

### github.com/oklog/oklog/pkg/group

生成的代码使用了这个库来完成了服务启动时需要启动多个后台goroutine的任务，每个人或者团队在这方面也许
都有自己的实践，可以进行service.go二次塑形，kit下一次执行不会再改动此文件（因为存在），当然你也可以直接用这个库，
只是你需要搞清楚它的用法。

## 9. 生成Client side代码

接下来我们使用kit生成grpc的client side代码：
```go
cd demo_project/
# -i 指 --pb_import_path，如果你前面没有指定pb/的位置，这里就不需要加上-i
kit g c hello -t grpc -i hello/pb
```
看一下client目录结构：
```go
c:\users\...\go-kit-examples\demo_project\hello
│  go.mod
│  go.sum
│
├─client
│  └─grpc
│          grpc.go
```

同样的，这里也需要完成req&rsp的转换工作, 打开client/grpc/grpc.go，修改如下：
```go
// Client-side:  endpoint-Req --encode-->> gRPC-Req
func encodeSayHiRequest(_ context.Context, request interface{}) (interface{}, error) {
	r := request.(*endpoint1.SayHiRequest)
	return &pb.SayHiRequest{Name: r.Name}, nil
}

// Client-side:  endpoint-Rsp <<--decode-- gRPC-Rsp
func decodeSayHiResponse(_ context.Context, reply interface{}) (interface{}, error) {
	r := reply.(*pb.SayHiReply)
	return &endpoint1.SayHiResponse{Reply: r.Reply},nil
}
```

## 10. 塑造适合你(的团队)的Client

先来看看`client/grpc/grpc.go`的New()方法：
```go
func New(conn *grpc.ClientConn, options map[string][]grpc1.ClientOption) (service.HelloService, error) {

	var sayHiEndpoint endpoint.Endpoint
	{
		sayHiEndpoint = grpc1.NewClient(conn, "pb.Hello", "SayHi", encodeSayHiRequest, decodeSayHiResponse, pb.SayHiReply{}, options["SayHi"]...).Endpoint()
		sayHiEndpoint = opentracing.TraceClient()
	}

	return endpoint1.Endpoints{
		SayHiEndpoint: sayHiEndpoint,
	}, nil
}
```

这是获取Service client的主要方法，如果你有一些微服务经验，你应该可以想到在RPC接口的client处需要添加几种
机制来使其更可靠、安全，更具有可观测性，可选的措施/机制列出如下：
- 服务发现(consul/etcd...)
- 负载均衡
- 重试
- 链路追踪
- 限速
- 断路器

同时还需注意这些措施的添加顺序，上面的措施应该更靠近底层conn对象。

另外，如果你有更多措施建议想添加到此文档中，可以通过issue告知我

为方便快速启动client，我就不添加服务发现/负载均衡到示例代码中了，如有需要可参考
`demo_project/new_addsvc/client/client.go`

下面是添加了部分措施之后的`client/grpc/grpc.go`:
```go
func New(conn *grpc.ClientConn) (service.HelloService, error) {
	/*
		Create some security measures
	*/
	var otTracer stdopentracing.Tracer
	otTracer = stdopentracing.GlobalTracer()
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	breaker := circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "SayHi",
		Timeout: 30 * time.Second,
	}))

	// Create go-kit grpc hooks, e.g.
	//      - grpctransport.ClientAfter(),
	//      - grpctransport.ClientFinalizer()
	// Injecting tracing ctx to grpc metadata, optionally.
	grpcBefore := grpc1.ClientBefore(opentracing.ContextToGRPC(otTracer, log.NewNopLogger()))
	/*
		Install into endpoints with above measures
	*/
	var sayHiEndpoint endpoint.Endpoint
	{
		sayHiEndpoint = grpc1.NewClient(conn, "pb.Hello", "SayHi",
			encodeSayHiRequest, decodeSayHiResponse, pb.SayHiReply{}, grpcBefore).Endpoint()
		sayHiEndpoint = opentracing.TraceClient(otTracer, "sayHi")(sayHiEndpoint)
		sayHiEndpoint = limiter(sayHiEndpoint)
		sayHiEndpoint = breaker(sayHiEndpoint)
	}

	return endpoint1.Endpoints{
		SayHiEndpoint: sayHiEndpoint,
	}, nil
}
```

断路器使用： https://github.com/sony/gobreaker （Sony公司写的，不过理论还是参考的[微软](https://docs.microsoft.com/en-us/previous-versions/msp-n-p/dn589784(v=pandp.10)?redirectedfrom=MSDN) ）

代码非常简洁，相关插件的使用方法可参考对应仓库，这里不再表述。

然后创建`hello/client/grpc/client.go`完成相关对象初始化工作以便caller service调用：
```go
package grpc

import (
	"context"
	"go-util/_util"
	"gokit_foundation"
	"google.golang.org/grpc"
	"hello/pkg/service"
	"time"
)

type Client struct {
	service.HelloService
	conn *grpc.ClientConn
}

var svcClient *Client

func newHelloClient(logger *gokit_foundation.Logger) *Client {
	var grpcOpts = []grpc.DialOption{
		grpc.WithInsecure(), // 因为没有使用tls，必须加上这个，否则连接失败
		grpc.WithBlock(),    // 若不加这项，远程服务断开再恢复时，网关调用会继续失败
	}
	var err error
	var conn *grpc.ClientConn
	var sc service.HelloService

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err = grpc.DialContext(ctx, "localhost:8082", grpcOpts...)
	// 出错时直接在这一层panic，外面就不需要处理
	logger.Must(err == nil, "HelloClient is nil")

	sc, err = NewSvc(conn)
	_util.PanicIfErr(err, nil)

	return &Client{
		HelloService: sc,
		conn:         conn,
	}
}

func MustNew(logger *gokit_foundation.Logger) *Client {
	if svcClient == nil {
		svcClient = newHelloClient(logger)
	}
	return svcClient
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

```

client的工作已经完成，现在我们来写一个test方法，创建`hello/client/grpc/grpc_test.go`：
```go
package grpc

import (
	"context"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	lgr := gokit_foundation.NewLogger(nil)
	c := MustNew(lgr)
	defer c.Close()
	reply, err := c.SayHi(context.Background(), "Jack Ma")
	if err != pbcommon.R_OK {
		t.Error(err)
	}
	lgr.Log("rsp", reply)
}
```

## 11. Let's test it now
```go
cd hello/client/grpc/
$ go test -run=TestNew
ts="2020-09-26 12:20:32" caller=client/grpc/grpc_in_use_test.go:18 rsp="Hi,Jack Ma"
PASS
ok      hello/client/grpc       1.118s
```

## 12. 自由尚在
你应该注意到，不管是kit，还是本仓库下的`go_project_template` （参考[golang-standards/project-layout][1]），
都没有涉及到数据访问层的目录规划，我想这是因为不同开发语言背景的开发团队/个人对这一层目录命名以及代码结构都有着不同的习惯，
比如Java背景的开发者习惯创建一个`dao`目录，代表的是Data access object, 当然，DAO并不和Java绑定，它是针对数据访问层的
一种设计模式，只是在Java Web开发历程中应用较深，并且卓有成效！

我这里就不详细介绍DAO模式了，你可以直接在项目根目录创建一个`dao`目录，这足够清晰明了。

但是其他语言背景的开发者也许就不太熟悉DAO这种模式了，在web开发领域，MVC是应用较广的一种软件设计架构（也可称框架模式），
这种架构下，Model（模型）是应用程序中用于处理应用程序数据逻辑的部分，所以你也可以直接创建一个`model`目录，它包含的文件结构
可以是这样：
```go
├─model
│      table_define.go  // 表结构定义
│      user.go          // user表的操作方法集合
│      user_wallet.go   // user_wallet表的操作方法集合
```
或许你还有自己的一套规划方案，这都不重要，重要的是解耦数据访问层与service层的代码，让项目的层次足够清晰，以便我们能够
持续的保持愉快心情来维护它，而不是在往后的某一天，一边敲代码一边操着一口流利的"what the f**k" 😤

在model层定义直接操作底层数据的方法，然后愉快的在service层引用，解耦你的代码！

## 13. 结束，新的开始

### 小结
本文档较为全面的介绍了如何使用[GrantZheng/kit](https://github.com/GrantZheng/kit) 作为go-kit框架的代码生成工具来辅助开发微服务，
在使用过程中，我发现了该工具的一些可以优化的问题以及bug，目前已经fork该项目进行优化以及bug修改。

### 改进
文档是花费我的业余时间编写的，因时间紧凑，无法避免一些小的瑕疵，以及更全面的介绍go-kit其他功能，例如
- 如何接入Etcd（服务发现）
- 如何接入OpenTelemetry（tracing）
- 如何深度接入Prometheus（monitor）

等等。。。如果您愿意一起来改进，欢迎提交您的PR！

### 微服务旅程建议
相比于快速开发出一个微服务，读者应该首先丰富自己对微服务的理论了解，没有稳固的理论基础，一切实践都是盲目的；
在开发过程中遇到一些微服务架构相关问题时，可能会丈二和尚摸不着头脑，对于团队来说，这是不可靠的。

这里分享一些微服务的理论、实践文章：

- [码农周报-微服务篇](https://github.com/chaseSpace/MNWeeklyCategory/blob/master/docs/MicroServiceLinks.md)
- [阿里云-正确入门Service Mesh](https://mp.weixin.qq.com/s/KHsxiOOHjTosQcd61rPsgg)
- [一文详解微服务架构知识](https://mp.weixin.qq.com/s/lpXkFsm01M9-27qeuo5JzA)
- [煎鱼博客: go与微服务相关讲解](https://eddycjy.com/posts/)

[kit]: https://github.com/chaseSpace/kit.git
[1]: https://github.com/golang-standards/project-layout/blob/master/README_zh.md