default: all

all: \
	clean \
	buildb \
	buildi \
	run \

buildb:
	@echo "Building binary..."
	@protoc counter/counter.proto --go_out=plugins=grpc:.
	@GOARCH=amd64 GOOS=linux go build -o grpctestapp main/main.go

buildi:
	@docker build -t grpctestapp -f Dockerfile .

clean:
	@rm -f grpctestapp

run:
	@docker deploy --compose-file docker-compose.yaml grpctest

stop:
	@docker stack rm grpctest

log:
	@docker service logs -f grpctest_counter

.PHONY: test buildb buildi clean run stop