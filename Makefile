default:
	go install github.com/yanfali/gensend-go
get:
	go get github.com/yanfali/gensend-go
initdb:
	gensend-go db init
	gensend-go db test
test:
	go test
