all: deps build

export GOPROXY=https://goproxy.cn
deps:
	go mod download

build:
	GOOS=windows go build -o ../bin/chat.exe
	GOOS=linux go build -o ../bin/chat
	GOOS=darwin go build -o ../bin/chat_mac
