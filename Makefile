.PHONY: generate
generate:
	buf generate proto

.PHONY: run
run:
	docker compose up --build