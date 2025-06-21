VERSION := $(shell cat VERSION | head -1)

.PHONY: build
build:
	docker compose up

backend-build:
	docker build -f build/backend-dockerfile -t babble-backend .
	docker run --rm \
	-e BABBLE_DB_HOST=10.100.0.6 \
	-e BABBLE_DB_USER=myuser \
	-e BABBLE_DB_PASSWORD=mypassword \
	-e BABBLE_DB_DBNAME=babble \
	-e BABBLE_DB_PORT=3306 \
	-e BABBLE_API_PRIVATE_KEY=taisthebest \
	-p 23456:80 \
	--add-host=host.docker.internal:host-gateway \
	babble-backend -vvv

cli-build:
	@docker build -f build/backend-cli-dockerfile -t babble-cli-backend .
	@docker run --rm \
	-e BABBLE_DB_HOST=10.100.0.6 \
	-e BABBLE_DB_USER=myuser \
	-e BABBLE_DB_PASSWORD=mypassword \
	-e BABBLE_DB_DBNAME=babble \
	-e BABBLE_DB_PORT=3306 \
	-e BABBLE_API_PRIVATE_KEY=taisthebest \
	-p 23456:80 \
	--add-host=host.docker.internal:host-gateway \
	babble-cli-backend $(CMD)

frontend-build:
	docker build -f build/frontend-dockerfile -t babble-frontend .
	docker run --rm \
	-p 3000:80 \
	--add-host=host.docker.internal:host-gateway \
	babble-frontend

.PHONY: test
backend-test:
	go test -C ./backend -v ./... -count=1

imagesize:
	@echo "Checking image sizes..."
	@docker images | grep babble

format:
	@echo "Formatting..."
	@pre-commit run --all-files

image:
	@echo $(VERSION)
	docker build -t shsung/babble:$(VERSION) .
	# docker push shsung/babble:$(VERSION)

.PHONY: frontend
frontend:
	npm run dev --prefix frontend

changelog:
	git cliff --unreleased --tag $(VERSION) --prepend changelog.md

revision:
	cd schema && alembic revision --autogenerate -m "$(REV)"
