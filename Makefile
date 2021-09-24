build:
	go build -o bin/needle-pairwise cmd/needle-pairwise/main.go
	go build -o bin/needle-group-seq cmd/needle-group-seq/main.go
	go build -o bin/needle-all-vs-all cmd/needle-all-vs-all/main.go

test:
	go test -v

install:
	cp bin/needle-pairwise /usr/local/bin/.
	cp bin/needle-group-seq /usr/local/bin/.
	cp bin/needle-all-vs-all /usr/local/bin/.

uninstall:
	rm /usr/local/bin/needle-pairwise
	rm /usr/local/bin/needle-group-seq
	rm /usr/local/bin/needle-all-vs-all