:: If compile.bat exists, kit tool will never update/delete it again.
:: NOTE: THIS SCRIPT CAN NOT CONTAINS CHINESE, WHICH CAUSE ERROR CODE.

protoc hello.proto --go_out=plugins=grpc:.