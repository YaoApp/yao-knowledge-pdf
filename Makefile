all: clean linux macos

macos: 
	go build -o ./dist/pdf-macos.so .
	chmod +x ./dist/pdf-macos.so

linux: 
	GOOS=linux GOARCH=amd64  go build -o ./dist/pdf-linux.so .
	chmod +x ./dist/pdf-linux.so

clean:
	rm -rf ./dist/pdf-macos.so
	rm -rf ./dist/pdf-linux.so