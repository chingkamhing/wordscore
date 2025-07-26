.PHONY: build
build:
	go build -o wordscore .

.PHONY: clean
clean:
	rm -f wordscore

.PHONY: benchmark
benchmark:
	go test -benchmem -bench .

.PHONY: test
test:
	go test .
