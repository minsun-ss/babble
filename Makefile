build:
	docker build -t babel .
	docker run --rm -p 23456:23456 babel /main
