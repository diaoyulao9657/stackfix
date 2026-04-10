build:
	go build -o stackfix .

install:
	go install .

test:
	go test ./... -v

clean:
	rm -f stackfix

.PHONY: build install test clean
