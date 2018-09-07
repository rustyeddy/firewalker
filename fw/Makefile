cmd = fw

build:
	make -C .. build
	go build -o $(cmd) main.go

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
