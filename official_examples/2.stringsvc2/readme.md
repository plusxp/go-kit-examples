stringsvc2 演示的go-kit服务对各层做了简单的文件分离，包含了

- Service 定义
- Endpoint 定义
- transport 封装（HTTP）

- metrics 采集 (Prometheus)
- logging 记录

tips: 阅读代码时请关注go-kit中middleware的使用