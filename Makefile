server:
	go build -o bin/server server/src/*.go

clean:
	rm -rf bin/*

.PHONY: server clean
