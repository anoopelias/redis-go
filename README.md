# Redis clone in Go
This is an experimental clone of a small subset of features of redis in go.

This repo started as an attempt towards solving CodeCrafter's ["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).

In this clone, we built a toy Redis clone that's capable of handling
basic commands like `PING`, `SET` and `GET`.

**Note**: Highly recommend
[codecrafters.io](https://codecrafters.io) challenges.

## Run

```
$ go run app/server.go
```

No `Makefile` yet.

## Test

```
$ go test app/redis_go/*.go
$ go test e2e/*.go
```

## Generate Mocks

From root directory:
```
mockgen -package=mocks redis-go/app/ev SysCall,SysCallError >app/mocks/mock_syscall.go
mockgen -package=mocks redis-go/app/redis_go RespReader >app/mocks/mock_redis_go.go
```
