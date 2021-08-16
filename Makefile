build:
	go build -o bin/needle-pairwise cmd/needle-pairwise/main.go
	go build -o bin/needle-group-seq cmd/needle-group-seq/main.go

test:
	go test -v