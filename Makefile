.PHONY: test

test:
	@richgo test ./... | grep --color=always --invert-match --fixed-strings '[no test files]'
