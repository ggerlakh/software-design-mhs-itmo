.PHONY: build run-server run-client clean

build:
	go build -o chat .

run-server:
	./chat -mode=server -port=8080

run-client:
	./chat -mode=client -host=158.160.218.84 -port=8080

clean:
	rm -f chat
