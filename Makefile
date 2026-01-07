.PHONY: build clean run

build:
	go build -o rogue ./cmd/rogue

clean:
	rm -f rogue

run: build
	./rogue
