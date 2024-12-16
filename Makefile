build:
	docker build -t zippygo .
	docker run --rm -p 23456:23456 zippygo /main
