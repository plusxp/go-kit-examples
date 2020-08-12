# go-kit-examples

[gokit.io](https://gokit.io)

## 官方介绍
Go-kit是使用Go语言构建微服务的一个工具箱，它可以解决分布式系统架构中的常见问题；
能够让我们专注于业务代码。

-   Github stars 17.7k

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
