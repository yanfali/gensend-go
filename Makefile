RM := /bin/rm
FILES := gensend-go gensend-go.test
default: clean
	go install github.com/yanfali/gensend-go
get:
	go get github.com/yanfali/gensend-go
initdb:
	gensend-go db init
	gensend-go db test
test:
	go test
clean:
	go clean
	$(RM) -f $(FILES)
