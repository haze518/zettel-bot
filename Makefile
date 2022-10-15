.PHONY: build
build:
	docker-compose build

.PHONY: run
run:
	docker-compose up --remove-orphans $(options)

.PHONY: stop
stop:
	docker-compose down --remove-orphans $(options)
