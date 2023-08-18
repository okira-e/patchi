build:
	go build -o bin/patchi -v

test:
	go test -v ./...

clean:
	rm -rf bin/patchi