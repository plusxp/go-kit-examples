
# 前言(Preface)
本文档将试图以最简明的方式讲述如何在实际环境中使用go-kit快速编写一个新的微服务项目，
如有逻辑纰漏，感谢指出。

<br/>

在开始前，建议读者先查看demo_project目录下的new_addsvc，了解go-kit的设计思想，以及
成型后的代码布局，继续阅读视为读者已经了解，不再过多解释相关。

---

了解`new_addsvc`的项目布局之后，读者应该知道如果所有的代码都需要自己编写，那使用go-kit开发只会降低
我们的开发效率，这显然无法接受；但是在仔细观察该服务的具体代码后，我们可以清楚的了解到，有几个层次的代码
都是根据service层的实现而编写的，它们的代码**模式**非常的有迹可循，经常使用go开发的读者很快可以想到，
为什么这些层次的代码不能使用工具生成呢？是的，早已有人想到这一点，它们开发出了go-kit的CLI工具，有了这个
工具，我们使用go-kit开发那就是如虎添翼，下面一起看看。

我调研了[go-kit](https://github.com/go-kit/kit)官方文档中列出的几个代码生成工具，找到了一个比较适合
使用的repo，那就是[GrantZheng/kit](https://github.com/GrantZheng/kit)，其他repo要么比较简陋，要么
正在开发，或者不再维护；这个repo也是fork一个不在维护的400+star的项目而来，该作者声称自己所在团队已深度使用
go-kit，无法接受没有可靠的go-kit辅助工具，所以自己fork来继续维护了（点赞！）

# 目录(Content)
-   [关于GrantZheng/kit](#关于GrantZheng/kit)
-   [开始(Getting Started)](#开始(Getting Started))


# 关于GrantZheng/kit
我们需要知道它的一些功能特点

-   可以生成指定名称的service模板代码，包含endpoint、transport(http/grpc)层
-   可以生成client代码
-   生成的以`_gen.go`结尾的文件不可修改，在重新生成时会覆盖
-   不会删除/覆盖已有的任何代码(除了`*_gen.go`)

> (以下部分描述译自GrantZheng/kit README.md, 少部分改动以适配本仓库，并带有额外的说明)  

# 开始(Getting Started)

首先安装go-kit CLI工具
```bash
# 默认使用modules包管理
# 注意：这个仓库是fork而来，go.mod的module name是原仓库名，就不能通过go get方式下载
mkdir $GOPATH/pkg/mod/git_repo -p
cd $GOPATH/pkg/mod/git_repo
git clone https://github.com/GrantZheng/kit.git
cd kit
go install 

# check usage
kit help
```

# 创建Project

kit文档说的是创建service，但这里替换为project或许更有助于理解

<br/>

```bash
kit new service hello
kit n s hello # 缩写
# 若要指定生成的go.mod内的模块名，指令后加上 --module module_hello，缩写-m，默认使用项目名作为模块名
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

# 生成Service模板

```bash
kit g s hello
kit g s hello --dmw # 创建默认middleware
kit g s hello -t grpc # 指定 transport (default is http)
```

这些命令会执行以下操作：

- 创建service层代码模板: `hello/pkg/service/service.go`
- 创建service层middleware: `hello/pkg/service/middleware.go`
- 创建endpoint: `hello/pkg/endpoint/endpoint.go` and `hello/pkg/endpoint/endpoint_gen.go`
- `--dmw` 创建endpoint middleware: `hello/pkg/endpoint/middleware.go`
- 创建transport files e.x http: `service-name/pkg/http/handler.go`
- 若使用`-t grpc`，则创建grpc transport files: `service-name/pkg/grpc/handler.go
- 创建service main file  
`hello/cmd/service/service.go`  
`hello/cmd/service/service_gen.go`  
`hello/cmd/main.go`

由于grpc作为微服务架构中常用的rpc选择，所以在这里我们直接执行 `kit g s hello --dmw -t grpc`，
在执行之前，首先需要在service.go中添加我们的api定义，示例：

```go
// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	SayHi(ctx context.Context, name, say string) (reply string, err error)
}
```
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
    │          compile.bat
    │          hello.pb.go
    │          hello.proto
    │
    └─service
            middleware.go
            service.go
```

注意，kit工具在/pkg目录生成了grpc目录，并且将pb目录也放在其中，根据`go_project_template`项目布局，/pb目录
应该放在项目根目录，这样方便快速找到一个服务的pb文件。

所以，我们需要