#!/usr/bin/env bash
source _import_echo_color.sh
source _os_util.sh

# 执行各种构建、安装、分析等操作的脚本
# 函数命名习惯：main func内调用的func命名以 fn_ 开头, 其他func内调用的func命名以 _fn_开头
#
flag=$1
project_dir=$2

# 默认以脚本执行目录为项目根目录
PROJECT_DIR=$(dirname $(pwd))
declare PROTO_OUTPUT_DIR
declare CMD_ARRAY=()

fn_init_cmd() {
	# ------------------- 所有的CMD选项 ----------------------
	readonly CMD_ARRAY=("gen" "gofmt" "govet")
	# ...CMD_on_ok后缀的指令 表示 CMD指令执行成功后要继续执行的指令，类似的还有_on_fail,  _on_any
	# 新增命令，只需在这里定义即可，无需其他操作
	# 注意：这里定义的变量当做全局变量使用，请不要在此函数外定义xxx_cmd这样的变量，会干扰

	# gen, 关于protoc命令，若后续在pb/proto/下增加目录，就需要适当添加相应目录到命令中(-I=../pb/proto/sub_folder)，仅添加proto文件则无需修改命令
	readonly    gen_cmd="protoc -I=../pb/proto ../pb/proto/*.proto --go_out=plugins=grpc:$PROTO_OUTPUT_DIR"
	readonly    gen_cmd_on_ok="echo gen proto ok"
	readonly    gen_cmd_on_fail="echo gen proto fail"
	# gofmt
	readonly    gofmt_cmd="gofmt -l -s -w $PROJECT_DIR"
	# govet
	readonly    govet_cmd="go vet $PROJECT_DIR/..." # go vet可能会修改go.mod文件，执行tidy来恢复
	readonly    govet_cmd_on_ok="go mod tidy"

}

_fn_execute() {
	local    cmd=$1
	local    cmd_on_ok=$2
	local    cmd_on_fail=$3
	local    cmd_on_any=$4

	echo    -e "EXECUTE> $cmd \n"
	echo    "******** output start *********"

	# 执行cmd对应指令
	$cmd
	local result_code=$?

	# 下面执行可能需要执行的附加指令

	$cmd_on_any    # 任何时候都需要执行的cmd

	if    [[ ${result_code} -eq 0 ]]; then
		$cmd_on_ok
		echo       -e "\n --- EXECUTE successful"
	else
		$cmd_on_fail
	fi

	echo    -e "\n******** output end *********"
}

# DO NOT NEED EDIT THIS FUNC
_fn_concat_if_not_empty() {
	local old=$1
	local mid=$2
	local append=$3

	if [[ -z "$append" ]]; then
		echo $old
		return
	fi
	echo "$old $mid $append\n"
}

# DO NOT NEED EDIT THIS FUNC
_fn_usage() {
	echo    "You need to provide flag to continue, see also below:"
	output="[Flag] <> [Cmd]\n"

	for flag in    "${CMD_ARRAY[@]}"; do
		cmd=$(      eval echo '$'"${flag}_cmd")
		cmd_on_ok=$(      eval echo '$'"${flag}_cmd_on_ok")
		cmd_on_fail=$(   eval echo '$'"${flag}_cmd_on_fail")
		cmd_on_any=$(   eval echo '$'"${flag}_cmd_on_any")

		# 注意双引号传参
		output=$(_fn_concat_if_not_empty "$output" "${flag} <>" "$cmd")
		output=$(_fn_concat_if_not_empty "$output" "---${flag}_cmd_on_ok <>" "$cmd_on_ok")
		output=$(_fn_concat_if_not_empty "$output" "---${flag}_cmd_on_fail <>" "$cmd_on_fail")
		output=$(_fn_concat_if_not_empty "$output" "---${flag}_cmd_on_any <>" "$cmd_on_any")
	done

	echo    -e $output | column -t -s "<>"
}

fn_parse_flag() {
	if    [[ -z $flag  ]] || [[ $flag == "-h" ]]; then
		_fn_usage
		return
	fi

	local    will_do_cmd
	local    will_do_cmd_on_ok
	local    will_do_cmd_on_fail
	local    will_do_cmd_on_any

	will_do_cmd=$(   eval echo '$'"${flag}_cmd")

	if    [[ -z ${will_do_cmd} ]]; then
		echo       "Invalid flag:$flag"
		return
	fi

	# 获取变量的间接引用变量值
	will_do_cmd_on_ok=$(   eval echo '$'"${flag}_cmd_on_ok")
	will_do_cmd_on_fail=$(   eval echo '$'"${flag}_cmd_on_fail")
	will_do_cmd_on_any=$(   eval echo '$'"${flag}_cmd_on_any")
	#	echo $will_do_cmd 111
	#	echo $will_do_cmd_on_ok 222
	#	echo $will_do_cmd_on_fail 333
	#	echo $will_do_cmd_on_any 444

	# 调用执行方法(每个参数都含有空格，需要双引号包括)
	_fn_execute    "$will_do_cmd" "$will_do_cmd_on_ok" "$will_do_cmd_on_fail" "$will_do_cmd_on_any"
}

fn_init() {
	# 替换PROJECT_DIR
	if    [[ -n $project_dir     ]]; then
		PROJECT_DIR=$project_dir
	fi

	PROTO_OUTPUT_DIR=$(dirname "$PROJECT_DIR")
	# windows上运行需要转换路径为windows路径，d:\\xxx
	if [[ $(_fn_is_windows) == 'true' ]]; then
		PROTO_OUTPUT_DIR=$(_fn_convert_to_windows_path $PROTO_OUTPUT_DIR)
	fi
	readonly    PROJECT_DIR
	readonly    PROTO_OUTPUT_DIR

	imported_fn_echo_color_msg    "> initial vars"
	_echo="
    PROJECT_DIR: $PROJECT_DIR
    PROTO_OUTPUT_DIR: $PROTO_OUTPUT_DIR
    "
	imported_fn_echo_color_msg    "$_echo"
}

main() {
	#	msg="***main.sh started***" # 问题：传入的前三个*全部丢失
	msg="---------- main.sh started ----------"
	imported_fn_echo_color_msg    'textcolor_red' "$msg"

	fn_init
	fn_init_cmd
	fn_parse_flag
}

# start
main

# comment
<<EOF
cd /go/go-kit-examples/template/script
- Check usage:
./main.sh -h

- Usage examples:
./main.sh gen ../../../  后面这个是项目根目录所在路径, 如/path/to/new_addsvc，注意proto文件中 "option go_package"最好设置为以项目根目录名开头的pkg名
./main.sh gofmt ../../   注意，为避免代码格式化的范围超出你的预期，可以用绝对路径指定，如/path/to/project_root
./main.sh govet ../../   这里也可以用绝对路径指定，比如/path/to/project_root

# 以上命令末尾也可不加路径，默认是new_addsvc所在路径，若要移植脚本适配其他项目，只需修改 fn_init_cmd 函数中定义的cmd即可
EOF

# 脚本格式化
# cd go-util/tool
# ./shfmt.exe -s -w -i 4 -bn -ci -sr -kp ../script/main.sh

# 关于shfmt工具， https://github.com/mvdan/sh
