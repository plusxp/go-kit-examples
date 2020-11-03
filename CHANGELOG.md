## 2020年10月25日 Update
- [kit][1] 现已支持 `go get`命令安装，无需clone后安装
- [kit][1] 增加 `-v` flag以支持查看version
- [kit][1] `cmd/service.go` 内的initXXX函数写法优化，解决端口冲突后运行不正常
- `demo_project/hello` 新增定时任务示例，包含分布式锁，redis相关

## 2020年10月21日 Update
- 修复[kit][1]在go1.15(可能包含其他go版本，但go1.12是正常的)下不能正常工作的bug
- 修正注释：`gateway/hellosvc.go`
- `GettingStart.md`（部分文字表述更新）,`README.md`（优化变量命名）
- 优化`go-util/_util/common_util.go`中的 `ListenSignalTask`函数

## 2020年10月10日 Update
- `demo_project/new_addsvc`优化`main.go`，添加更详细的注释
- 更新`README.md` (`go_project_template`部分)
- 抽离公共bash scripts到 `bash-util/`
- 更新`GettingStart.md`（`小结`部分）

## 2020年10月9日 Update
- `demo_project/hello`增加新接口`UpdateUserInfo`
- 网关增加`Prepare`方法，完成req鉴权(jwt)和参数反序列化操作
- 简单优化已有的网关接口逻辑
- [Kit][1] 仓库更新:
    -   生成service的方法名中的Reply更新为Response
    -   endpoint方法注释简单优化

## 2020年9月30日20:04 Update
-  main.sh支持在windows上运行

## 2020年9月26日12:23 Update
- `gokit_foundation\gateway` 优化，支持标准proto/gogo.proto Message的marshal，演示如何解决json序列化零值字段被忽略的问题
- `gokit_foundation\log.go` 优化，简化结构及方法
- 更新`GettingStart.md`，匹配最新kit

## 2020年9月19日12:32 Update 
- 基本完成 `/demo_project/hello`  
    增加了两个不同风格的接口示例
    - `client/grpc` :white_check_mark:
    - `pb` dir :white_check_mark:
    - `pkg` dir :white_check_mark:
- 基本完成 `/demo_project/gateway`，使用mux作为路由器
- 完善 `/gokit_foundation`，增加完善的logger机制

- [Kit][1] 仓库更新:
    -   :tada: 增加[examples/hello][2]，仅包含grpc作为transport, 并包含[文档][3]
    -   生成的代码中，endpoint层的参数req&rsp更新为指针类型
    -   大幅优化对proto文件的支持，现支持根据当前shell解释器环境生成对应的proto文件编译脚本（以前是根据操作系统类型），修复部分bug
    -   若在指定的路径下已包含proto编译脚本，则自动执行脚本（当然也可自己执行）
 
 
[1]:https://github.com/chaseSpace/kit
[2]:https://github.com/chaseSpace/kit/tree/master/examples
[3]:https://github.com/chaseSpace/kit/blob/master/examples/hellosvc_doc.md