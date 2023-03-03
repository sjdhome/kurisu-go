.PHONY: release
debug: build/debug/kurisu
release: build/release/kurisu
clean:
	rm -r build || true
	rm -r web/static/js/* || true

build/debug/kurisu: clean
	npm run compile
	go build -o $@ main.go
build/release/kurisu: clean
	npm run compile
	go build -ldflags "-s" -o $@ main.go
