build:
	docker build -t babel .
	docker run --rm \
	-e BABEL_DB_HOST=host.docker.internal \
	-e BABEL_DB_USER=myuser \
	-e BABEL_DB_PASSWORD=mypassword \
	-e BABEL_DB_DBNAME=babel \
	-e BABEL_DB_PORT=3306 \
	-p 23456:23456 \
	--add-host=host.docker.internal:host-gateway \
	babel -vvv

.PHONY: test
test:
	go test -v ./... -count=1

imagecheck:
	@echo "Checking image sizes..."
	@docker images babel
