all:
	go build -buildmode=c-shared -o saferm.so main.go
clean:
	rm -f saferm.h saferm.so
