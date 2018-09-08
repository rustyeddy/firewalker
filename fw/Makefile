cmd = fw
files = *.go

build:
	make -C .. build
	go build -o $(cmd) $(files)

test:
	go test

vet:
	go vet

run:
	make -C .. build
	go run main.go

clean:
	go clean
	rm -rf ff *~ 
