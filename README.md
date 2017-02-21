# GRPC TEST

This is a sample that shows a bit about grpc and docker swarm. The demo contains
two kind of services. A GRPC server that works as a counter and a HTTP server which
use the counter service to keep track of the number of times the endpoint /count
has been hit

## Requirements (OSX)

### Go lang

This example is written in Go

```shell
$ brew install go
```

### Glide

Go package manager

```shell
$ brew install glide
```

### GRPC

```shell
$ brew install protobuf
$ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

### Docker

```shell
$ brew cask install docker
```

You have to configure docker in swarm mode

```shell
$ docker swarm init
```

## Usage

To start the demo run:

```shell
$ make
```

This will run one replica of the server and three replicas of the client. To see
some action run this:

```shell
$ while true;  do curl "http://localhost:8080/count"; sleep 0.5; done
```

and in a separate terminal run:

```shell
$ make log
```

To stop the services run:

```shell
$ make stop
```


