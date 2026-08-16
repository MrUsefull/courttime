bin:
	mkdir -p bin

lint:
	golangci-lint run

test: bin
	go test -count=1 ./...

clean:
	rm -rf bin
