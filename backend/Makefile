.PHONY: test test-cov

test:
	cd .. && docker compose run --rm api_test go test -v ./...

test-cov:
	cd .. && docker compose run --rm api_test go test -cover ./...
