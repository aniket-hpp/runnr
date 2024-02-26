all: runnr

runnr: main.go
	go build -o runnr .

install:
	go build -o runnr .
	sudo cp ./runnr /usr/local/bin
	mkdir -p ~/.runnr
	cp -r ./docs ~/.runnr
	cp -r ./templates ~/.runnr

uninstall:
	sudo rm /usr/local/bin/runnr
