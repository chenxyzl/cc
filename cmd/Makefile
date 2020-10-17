all: deps


deps:
	GOOS=windows go build -o ../bin/chat.exe
	GOOS=linux go build -o ../bin/chat
	GOOS=darwin go build -o ../bin/chat_mac