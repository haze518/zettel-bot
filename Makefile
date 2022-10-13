.PHONY: build
build:
	docker build -t zettel .

.PHONY: run
run:
	docker run zettel