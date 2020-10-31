# go-kit-examples
本仓库致力于为想要使用Go-kit构建微服务架构的中文开发者提供帮助。

仓库包含以下内容：
-   Go-kit框架中文介绍（README.md）
-   带中文注释的go-kit官方示例，包含循序渐进的多个[示例项目](https://github.com/chaseSpace/go-kit-examples/tree/master/official_examples)  
    (不过官方的4,5示例仅作功能演示，代码不够简明，新手看起来会有些晦涩，建议关注下面的演示项目)
-   [完善优化后的new_addsvc(强烈推荐)](#完整可用于实战参考的项目)
-   [业界常用的Go项目目录结构(推荐了解)](#业界常用的Go项目目录结构)

演示项目：  
> 1. 演示项目使用kit辅助开发go-kit项目，以提高开发效率！项目代码十分简明，并包含充分注释，强烈推荐~  
> 2. hello和gateway项目目前涉及到的技术关键词: 
`router`, `grpc`, `trace`, `metrics`, `ratelimit`, `breaker`, 
`jwt`, `redis`, `timer task`, `distrib lock`
-   [使用kit代码生成工具快速开发go-kit微服务-hellosvc](#使用代码生成工具快速开发go-kit微服务) 🐦
-   [go-kit API网关](#API网关) 🐤
___
-   [更新日志](#更新日志)
-   [Go-kit中文群组(推送仓库更新)](#go-kit中文群组)

## Go-kit官方介绍
[![](https://img.shields.io/static/v1?label=Github&message=go-kit&color=important)](https://github.com/go-kit/kit)
![](https://badgen.net/github/stars/go-kit/kit)
![](https://badgen.net/github/release/go-kit/kit)

Go-kit是使用Go语言构建微服务的一个工具箱，它可以解决分布式系统架构中的常见问题；
能够让我们专注于业务代码。

## 特点

工具集，通过不同pkg支持auth/metrics/log/service-discovery/tracing/transport, 
并且可以自己决定哪个功能/协议使用哪个库/后端，比如服务发现可以使用consul或者etcd，都有集成好的模块直接使用。

## 框架分层

 1. Transport layer (HTTP, gRPC, Thrift, and net/rpc, 自由定制)
 2. Endpoint layer
 3. Service layer (聚焦业务逻辑)

请求由第一层按序流向第三层。 

### Endpoints

   相当于controller内的一个handler，如果同时实现HTTP和gRPC协议，可以方便的将两种请求发往同一个endpoint。

### Services

  - 这一层主要聚焦业务逻辑，一个service可以被多个endpoint调用，一个endpoint可以调用多个service方法
  - 与gRPC一样，一个service在代码是作为一个interface，需要我们实现这个interface
  - service层只关注业务，与endpoint/transport无关，代码层分离
  
## Middlewares
 
 Go kit使用中间件来添加更多功能，包括日志、限速、负载均衡、链路追踪等等。
 
## 框架设计
典型的洋葱模型

![Design](https://gokit.io/faq/onion.png)

这些层完全可以划分为三个域：

-   最靠内的--Service域

    `服务定义以及具体业务逻辑的地方`
-   中间的--Endpoint域

    `接口的抽象，称为端点，一个端点对外映射为一个接口，一个端点内可以调用Service域定义的一个或多个方法；
    在这里会实施保证接口安全和防脆弱的措施
    `
-  最外层的Transport域

    `endpoint绑定具体传输协议的地方`
    
## Error处理

service可能会返回err，有两种方式在endpoint中来封装err：

-  第一种，在service response struct中包含业务err
-  第二种，在将service err传递到endpoint的err

注意一点：endpoint返回的err会被中间件捕获到，比如断路器；所以这里区分网络导致的err和service返回的业务err。


## 服务发现

go-kit有相应的pkg支持Consul, etcd, ZooKeeper, and DNS SRV

## 监控

go-kit有组件支持现代化的监控系统Prometheus, 同时官方也推荐使用它来建立go-kit服务

## 业界常用的Go项目目录结构

规范的Go项目目录结构，参考[golang-standards/project-layout](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)

在此基础上我添加了proto文件的目录划分，请参考`go_project_template`目录，
关于相关目录的作用说明详见各子目录中的readme.md或代码；在`/script`目录存放了一个`main.sh`，可供参考，也可直接使用

> 按照Go项目规范组织代码结构，可以极大地减少沟通成本，提高团队开发效率

## 完整可用于实战参考的项目

[new_addsvc](https://github.com/chaseSpace/go-kit-examples/tree/master/demo_project/new_addsvc)

- `/pkg`目录包含了service、endpoint、transport三层的代码，前两者都有中间件，也可以在transport层添加中间件以实现完整的链路追踪
- `/pb`目录包含了存放proto文件的`proto/`目录和存放`*.pb.go`文件的`gen-go/`目录
- `/internal`目录包含了这个app私有的方法

这个项目会持续更新，包括项目目录结构，代码优化，不过基本骨架已搭成，后续要做的是提取可以提取的代码到foundation中，以及必要的结构调整，
较大更新会以日志形式贴出。

建议先拉取到本地研究/学习，跟随作者持续优化~

欢迎提出优化意见！ 

## 使用代码生成工具快速开发go-kit微服务

[GettingStart](https://github.com/chaseSpace/go-kit-examples/blob/master/GettingStart.md)

## API网关

- 使用具有强大路由和参数匹配功能的[mux](https://github.com/gorilla/mux) 库作为路由器（当然也可以使用你喜欢的库替换）
- 包含了grpc接口调用，并简单演示了如何使用mux的参数匹配功能
- 极为简洁实用的代码

[Gateway](https://github.com/chaseSpace/go-kit-examples/tree/master/demo_project/gateway) 

## 更新日志

[CHANGELOG][CHANGELOG] (上次更新于2020年10月25日)

[CHANGELOG]:https://github.com/chaseSpace/go-kit-examples/blob/master/CHANGELOG.md

## TODO

-   增加定时任务的实战示例（gateway+new_addsvc）
-   gateway增加service discovery(consul)

## Go-kit中文群组

![](https://github.com/chaseSpace/go-kit-examples/blob/master/wxgroup.jpg)