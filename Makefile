build:
	go build -o bin/server server/src/*.go

clean:
	rm -rf bin/*

.PHONY: build clean
