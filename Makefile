VERSION := $(shell cat VERSION | head -1)

backend-build:
	docker build -f build/backend-dockerfile -t babel-backend .
	docker run --rm \
	-e BABEL_DB_HOST=10.100.0.6 \
	-e BABEL_DB_USER=myuser \
	-e BABEL_DB_PASSWORD=mypassword \
	-e BABEL_DB_DBNAME=babel \
	-e BABEL_DB_PORT=3306 \
	-p 23456:80 \
	--add-host=host.docker.internal:host-gateway \
	babel-backend -vvv

frontend-build:
	docker build -f build/frontend-dockerfile -t babel-frontend .
	docker run --rm \
	-p 3000:3000 \
	--add-host=host.docker.internal:host-gateway \
	babel-frontend

.PHONY: test
backend-test:
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

frontend:
	npm run dev --prefix frontend

build2:
	npm run build

run2:
	npx serve out

changelog:
	git cliff --unreleased --tag $(VERSION) --prepend changelog.md
