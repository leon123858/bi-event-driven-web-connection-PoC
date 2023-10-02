all:  
	

compile:
	# use `export PATH="$PATH:$(go env GOPATH)/bin"` in shell
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative todo.proto

init:
	# please use `brew install protobuf` install compiler first
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	