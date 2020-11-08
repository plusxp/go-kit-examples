:: Install proto3.
:: https://github.com/google/protobuf/releases
:: Update protoc Go bindings via
::  go get -u github.com/golang/protobuf/proto
::  go get -u github.com/golang/protobuf/protoc-gen-go
::
:: See also
::  https://github.com/grpc/grpc-go/tree/master/examples

:: On windows cmd, the wildcard like `*.proto` is not supported,
:: so there would be only generate it with for loop.

@echo off
set matchExp="*.proto"
set matchExp2="pbcommon/*.proto"

for %%i in ("%matchExp%") do (
    protoc %%i --go_out=plugins=grpc:../../../
)

for %%i in ("%matchExp2%") do (
    protoc pbcommon/%%i --go_out=plugins=grpc:../../../
)

echo protoc exec successful!
