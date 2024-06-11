build: bin/blog_system

bin/%: $(shell find . -type f -name '*.go') # make sure to rebuild if any go file changed.
	go build -o bin/$(@F) cmd/$(@F)/main.go