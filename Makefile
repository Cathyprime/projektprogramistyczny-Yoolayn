server:
	go build -o bin/server cmd/main.go

clean:
	rm -rf bin/*

.PHONY: server clean
