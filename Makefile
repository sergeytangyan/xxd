build:
	go build -o bin/ccxxd.exe cmd/xxd/main.go

build-lean:
	go build -o bin/ccxxd.exe -ldflags "-s -w" cmd/xxd/main.go

clean:
	rm -rf bin/