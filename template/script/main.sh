#!/usr/bin/env bash

# 执行各种构建、安装、分析等操作的脚本

flag=$1
project_dir=$2

# 默认以脚步执行目录为项目根目录
PROJECT_DIR=$(dirname $(pwd))
declare PROTO_OUTPUT_DIR
declare -A flagMap

FLAG_CMD_GEN_PROTO="gen"
FLAG_CMD_GOFMT="gofmt"
FLAG_CMD_GOVET="govet"

func_init_cmd_map() {
    # 在这里声明你要添加的脚本命令，包括FLAG, CONTENT
    flagMap=(
        # 关于proto命令，若后续在pb/proto/下增加文件夹，就需要适当添加相应目录到命令中，仅添加proto文件则无需修改命令
        [$FLAG_CMD_GEN_PROTO]="protoc
        -I=../pb/proto
        -I=../pb/proto/common_pb
        ../pb/proto/*.proto
        ../pb/proto/common_pb/*.proto
        --go_out=plugins=grpc:$PROTO_OUTPUT_DIR"
        [$FLAG_CMD_GOFMT]="gofmt -l -s -w $PROJECT_DIR" # 格式化代码(注意：会直接修改代码)
        [$FLAG_CMD_GOVET]="go vet $PROJECT_DIR/..." # 检测代码中可能的错误（目录必须是项目根目录，包含go.mod文件，否则报错）
    )
    readonly flagMap
}

#echo ${!flagMap[@]} # all keys
#echo ${#flagMap[*]} # map len

func_parse_flag() {
  	if [ -z "$flag" ] || [[ $flag = "-h" ]]; then
        func_usage
        return
    fi

    cmd=${flagMap[$flag]}

    if [ -z "$cmd" ]; then
        echo "Invalid flag: $flag"
        return
    fi

    echo -e "EXECUTE> $cmd \n"
    echo "******** output start *********"

    # 执行cmd对应指令
    ${cmd}
    result_code=$?

    echo -e "\n******** output end *********"
    # go vet可能会修改go.mod文件，执行tidy来恢复
    if [ $flag == $FLAG_CMD_GOVET ]; then
        go mod tidy
    fi

    if [[ result_code -eq 0 ]]; then
        echo -e "\n***EXECUTE successful.***"
    fi
}

func_usage() {
   	echo "You need to provide flag to continue, see also below:"
    output="[Flag] <> [Cmd]\n"
    for cmd in ${!flagMap[*]}; do
        output="$output $cmd <> ${flagMap[$cmd]}\n"
    done
    echo -e $output | column -t -s "<>"
}

func_init() {
    # 替换PROJECT_DIR
    if [ ! -z "$project_dir" ]; then
        PROJECT_DIR=$project_dir
    fi

    PROTO_OUTPUT_DIR=$(dirname "$PROJECT_DIR")
    readonly PROJECT_DIR
    readonly PROTO_OUTPUT_DIR

    echo "
    PROJECT_DIR: $PROJECT_DIR
    PROTO_OUTPUT_DIR: $PROTO_OUTPUT_DIR
    "
}

main() {
    echo '***main.sh started***'
    func_init
    func_init_cmd_map
    func_parse_flag
}

main

# cd /go/go-kit-examples/template/script
# Check usage:
# ./main.sh -h

# usage example:
# $./main.sh gen ../..  后面这个是项目目录, 如/root/go-kit-examples，注意proto文件中 "option go_package"需要设置为带有项目根目录路径的pkg名
# $./main.sh gofmt ../../   注意，为避免代码格式化的范围超出你的预期，可以用绝对路径指定，如/path/to/project_root
# $./main.sh govet ../../   这里也可以用绝对路径指定，比如/path/to/project_root

# shell script formatting
# cd ../tool
# ./shfmt.exe -s -w -i 4 -bn -ci -sr -kp ../script/main.sh

# 关于shfmt工具， https://github.com/mvdan/sh