.PHONY: release
debug: build/debug/kurisu
release: build/release/kurisu
clean:
	rm -r build

build/debug/kurisu: clean
	go build -o $@ app/main.go
build/release/kurisu: clean
	go build -ldflags "-s" -o $@ app/main.go
