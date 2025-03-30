VERSION := $(shell cat VERSION | head -1)

build:
	docker build -t babel .
	docker run --rm \
	-e BABEL_DB_HOST=10.100.0.6 \
	-e BABEL_DB_USER=myuser \
	-e BABEL_DB_PASSWORD=mypassword \
	-e BABEL_DB_DBNAME=babel \
	-e BABEL_DB_PORT=3306 \
	-p 23456:80 \
	--add-host=host.docker.internal:host-gateway \
	babel -vvv

.PHONY: test
test:
	go test -C ./backend -v ./... -count=1

imagecheck:
	@echo "Checking image sizes..."
	@docker images babel

format:
	@echo "Formatting..."
	@pre-commit run --all-files

image:
	@echo $(VERSION)
	docker build -t shsung/babel:$(VERSION) .
	docker push shsung/babel:$(VERSION)
