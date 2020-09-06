
# cmd目录

存放结构：
```text
~/cmd/some_app/main.go
~/cmd/some_tool/main.go
...
```

`/cmd/some_app/`目录应该只包含main.go，这是app的启动入口，它可以import`/pkg`和`/internal`  
`/cmd/some_tool/`目录一般包含main.go在内的多个go文件，一般 **不会** import`/pkg`和`/internal`，可以import `/util`,
它的使用方式是先`go install`，然后运行可执行文件

参考[prometheus/cmd][1]

[1]: https://github.com/prometheus/prometheus/tree/master/cmd