cmd = fw

# make -C above will cause this directory to be built
all:
	make -C $(cmd) build 

build:
	go build 

buildv:
	go build -v

run:
	make -C $(cmd) run

test:
	go test

testv:
	go test -v

install:
	make -C $(cmd) install -i

clean:
	go clean
	rm *~
	make -C $(cmd) clean
