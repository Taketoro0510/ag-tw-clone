.PHONY: test test-backend test-frontend test-cov

test: test-backend test-frontend

test-backend:
	$(MAKE) -C backend test

test-cov:
	$(MAKE) -C backend test-cov

test-frontend:
	docker compose up -d
	$(MAKE) -C frontend test
