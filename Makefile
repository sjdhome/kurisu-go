.PHONY: release debug clean npm-task
.DEFAULT_GOAL := debug

clean:
	rm -r build/* || true

npm-task:
	npm run compile:node
	npm run compile:web
	npm run blog:generate

debug: clean npm-task
	go build -o build/kurisu main.go

release: clean npm-task
	go build -o build/kurisu -ldflags "-s" main.go
