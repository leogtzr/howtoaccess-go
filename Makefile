.DEFAULT_GOAL := install

# INSTALL_SCRIPT=./install.sh
BIN_FILE=howtoaccess

install:
	go test
	go build -o "${BIN_FILE}"

clean:
	go clean

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

run:
	./howtoaccess -input "${HOME}/inputs/HowToAccess.csv"
