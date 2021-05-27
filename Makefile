.PHONY: test

test:
	@richgo test ./... -timeout 10s
