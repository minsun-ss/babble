VERSION := $(shell cat VERSION | head -1)

.PHONY: build
build:
	docker compose up

backend-build:
	docker build -f build/backend-dockerfile -t babel-backend .
	docker run --rm \
	-e BABEL_DB_HOST=10.100.0.6 \
	-e BABEL_DB_USER=myuser \
	-e BABEL_DB_PASSWORD=mypassword \
	-e BABEL_DB_DBNAME=babel \
	-e BABEL_DB_PORT=3306 \
	-e BABEL_API_PRIVATE_KEY=taisthebest \
	-p 23456:80 \
	--add-host=host.docker.internal:host-gateway \
	babel-backend -vvv

cli-build:
	@docker build -f build/backend-cli-dockerfile -t babel-cli-backend .
	@docker run --rm \
	-e BABEL_DB_HOST=10.100.0.6 \
	-e BABEL_DB_USER=myuser \
	-e BABEL_DB_PASSWORD=mypassword \
	-e BABEL_DB_DBNAME=babel \
	-e BABEL_DB_PORT=3306 \
	-e BABEL_API_PRIVATE_KEY=taisthebest \
	-p 23456:80 \
	--add-host=host.docker.internal:host-gateway \
	babel-cli-backend $(CMD)

frontend-build:
	docker build -f build/frontend-dockerfile -t babel-frontend .
	docker run --rm \
	-p 3000:80 \
	--add-host=host.docker.internal:host-gateway \
	babel-frontend

.PHONY: test
backend-test:
	go test -C ./backend -v ./... -count=1

imagesize:
	@echo "Checking image sizes..."
	@docker images | grep babel

format:
	@echo "Formatting..."
	@pre-commit run --all-files

image:
	@echo $(VERSION)
	docker build -t shsung/babel:$(VERSION) .
	# docker push shsung/babel:$(VERSION)

.PHONY: frontend
frontend:
	npm run dev --prefix frontend

changelog:
	git cliff --unreleased --tag $(VERSION) --prepend changelog.md

revision:
	cd schema && alembic revision --autogenerate -m "$(REV)"
