#!/usr/bin/env sh

# Install proto3
# sudo apt-get install -y git autoconf automake libtool curl make g++ unzip
# git clone https://github.com/google/protobuf.git
# cd protobuf/
# ./autogen.sh
# ./configure
# make
# make check
# sudo make install
# sudo ldconfig # refresh shared library cache.
#
# Update protoc Go bindings via
#  go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
#
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

# 注意：common下的proto文件的go_package与当前路径下的proto文件不一致，两条命令貌似无法合并，会报：inconsistent package import paths:
protoc *.proto --go_out=plugins=grpc:../../../
protoc pbcommon/*.proto --go_out=plugins=grpc:../../../
