# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help


help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

go-security-scan:
	semgrep --config=p/trailofbits --config=p/jwt --config=p/command-injection --config=p/xss --config=p/insecure-transport --config=p/golang --config=p/owasp-top-ten --config=p/r2c-ci --config=p/r2c-security-audit --config=p/r2c-bug-scan --config=p/r2c-best-practices --config=p/secrets --config=p/sql-injection --config=p/r2c-best-practices --config=p/xss --config=p/ci --config=p/mobsfscan

security-scan:
	semgrep --config=auto


build_docker:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 linux/386 .
