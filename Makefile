build:
	go build

run:
	make -C cmd run

test:
	go test

testv:
	go test -v

clean:
	go clean
	rm *~
	make -C cmd clean
