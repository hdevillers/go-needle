build:
	go build -o bin/needle-pairwise cmd/needle-pairwise/main.go

test:
	go test -v