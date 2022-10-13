.PHONY: build
build:
	docker build -t zettel .

.PHONY: run
run:
	docker-compose up --remove-orphans $(options)

.PHONY: stop
stop:
	docker-compose down --remove-orphans $(options)
