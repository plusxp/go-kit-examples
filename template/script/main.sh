#!/usr/bin/env bash

# 执行各种构建、安装、分析等操作的脚本

flag=$1
project_dir=$2

# 默认以脚步执行目录为项目根目录
PROJECT_DIR=$(dirname $(pwd))

# 替换PROJECT_DIR
if [ ! -z "$project_dir" ]; then
    PROJECT_DIR=$project_dir
fi

_proto_output_dir=$(dirname "$PROJECT_DIR")
echo "
PROJECT_DIR: $PROJECT_DIR
_proto_output_dir: $_proto_output_dir
"

FLAG_CMD_GEN_PROTO="gen"
FLAG_CMD_GOFMT="gofmt"
FLAG_CMD_GOVET="govet"

# 在这里声明你要添加的脚步命令，包括FLAG, CONTENT
declare -A flagMap=(
  # 关于proto命令，若后续在pb/proto/下增加文件夹，就需要适当添加相应目录到命令中，仅添加proto文件则无需修改命令
  [$FLAG_CMD_GEN_PROTO]="protoc
    -I=../pb/proto
    -I=../pb/proto/common_pb
    ../pb/proto/*.proto
    ../pb/proto/common_pb/*.proto
    --go_out=plugins=grpc:$_proto_output_dir"

  [$FLAG_CMD_GOFMT]="gofmt -l -s -w $PROJECT_DIR" # 格式化代码(注意：会修改代码)
  [$FLAG_CMD_GOVET]="go vet $PROJECT_DIR..." # 检测代码中可能的错误（目录必须是项目根目录，包含go.mod文件，否则报错）
)


#echo ${!flagMap[@]} # all keys
#echo ${#flagMap[*]} # map len

parse_flag(){
  if [ -z "$flag" ]; then
      echo "You need to provide flag to continue, see also below:"
      usage
      return
  fi

  validFlag=false

  for cmd in ${!flagMap[*]}; do
    if [ $flag == $cmd ]; then
        validFlag=true

        echo -e "   EXECUTE> ${flagMap[$cmd]}"
        ${flagMap[$cmd]}

        # go vet可能会修改go.mod文件，执行tidy来恢复
        if [ $flag == $FLAG_CMD_GOVET ];then
            go mod tidy
        fi

        if [ $? -eq 0 ]; then
            echo "EXECUTE successful."
        fi

        break
    fi
  done

  if [ $validFlag == false ]; then
    echo "Invalid flag: $flag"
    fi
}

usage(){
  output="[Cmd] <> [Detail]\n"
  for cmd in ${!flagMap[*]}; do
    output="$output $cmd <> ${flagMap[$cmd]}\n"
  done
  echo -e $output | column -t -s "<>"
}

main() {
  echo '***main.sh started***'
  parse_flag
}

main

# example:
# $./main.sh gen ../../../..  后面这个是项目目录
# $./main.sh gofmt ../../
# $./main.sh govet ../../