NAME := af-zoomed-mnemonic
PORT := 8080

all:	fmt	image run
	@true

image:
	docker build -t ${NAME} .

run:
	@echo "Visit --> http://127.0.0.1:${PORT} (or boot2docker IP)"
	docker run --rm --name ${NAME}-run -p ${PORT}:${PORT} ${NAME}

fmt:
	gofmt -w .
