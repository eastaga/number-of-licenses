all : build test

test:
	go test ./helpers  -cover

build: fmt clean
	go build -o main main.go

fmt:
	gofmt -w ./helpers

clean:
	rm -f main