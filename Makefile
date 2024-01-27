server:
	go build -race -o bin/server cmd/server/main.go

full:
	go build -race -x -v -a -o bin/server cmd/server/main.go

release:
	go build -ldflags "-s -w" -race -o bin/server cmd/server/main.go

generator:
	go build -o bin/generator cmd/generator/gen.go

clean:
	rm -rf bin/*

.PHONY: server clean full
