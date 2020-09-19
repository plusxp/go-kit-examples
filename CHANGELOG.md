
## 2020年9月19日12:32 Update 
- 基本完成 `/demo_project/hello`  
    增加了两个不同风格的接口示例
    - `client/grpc` :white_check_mark:
    - `pb` dir :white_check_mark:
    - `pkg` dir :white_check_mark:
- 基本完成 `/demo_project/gateway`，使用mux作为路由器
- 完善 `/gokit_foundation`，增加完善的logger机制

- [Kit](#1) 仓库 update:
    -   :tada: 增加[examples/hello](#2)，仅包含grpc作为transport, 并包含[文档](#3)
    -   生成的代码中，endpoint层的参数req&rsp更新为指针类型
    -   大幅优化对proto文件的支持，现支持根据当前shell解释器环境生成对应的proto文件编译脚本（以前是根据操作系统类型），修复部分bug
    -   若在指定的路径下已包含proto编译脚本，则自动执行脚本（当然也可自己执行）
 
 
[1]:https://github.com/chaseSpace/kit
[2]:https://github.com/chaseSpace/kit/tree/master/examples
[3]:https://github.com/chaseSpace/kit/blob/master/examples/hellosvc_doc.md