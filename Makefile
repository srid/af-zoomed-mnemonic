NAME := af-zoomed-mnemonic
URL := afmnemonic.happyandrmless.com
PORT := 8080

all:	fmt image run
	@true

image:
	docker build -t ${NAME} .

run:
	@echo "Visit --> http://127.0.0.1:${PORT} (or boot2docker IP)"
	docker run --rm --name ${NAME}-run -p ${PORT}:${PORT} ${NAME}

runprod:
	docker rm -vf ${NAME}-run || true
	docker run --name=${NAME}-run -d -e VIRTUAL_HOST=${URL} ${NAME}
	docker logs -f ${NAME}-run

fmt:
	gofmt -w .
