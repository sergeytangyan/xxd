build:
	go build -o bin/ccxxd.exe cmd/ccxxd/main.go

build-lean:
	go build -o bin/ccxxd.exe -ldflags "-s -w" cmd/ccxxd/main.go

clean:
	rm -rf bin/