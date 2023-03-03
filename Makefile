.PHONY: release
debug: build/debug/kurisu
release: build/release/kurisu
clean:
	rm -r build || true

build/debug/kurisu: clean
	go build -o $@ main.go
build/release/kurisu: clean
	go build -ldflags "-s" -o $@ main.go
