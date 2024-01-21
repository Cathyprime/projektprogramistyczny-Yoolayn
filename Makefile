server:
	go build -race -o bin/server cmd/main.go

full:
	go build -race -x -v -a -o bin/server cmd/main.go

release:
	go build -ldflags "-s -w" -race -o bin/server cmd/main.go

clean:
	rm -rf bin/*

.PHONY: server clean full
