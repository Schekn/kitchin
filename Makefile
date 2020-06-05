.PHONY: build statik run doc test lint help

.DEFAULT_GOAL := help

APP := delivery

build: statik ## Build project
	go build -o build/$(APP)

run: statik ## Run application
	go run main.go -o orders_test.json -c config.yml

test: ## Run tests
	GO_ENV="testing" go test -v ./...

doc: ## Show documentation
	@echo "open browser at http://localhost:6060/pkg/delivery/"
	@godoc -http=:6060

lint: ## Lint code
	revive -config revive.toml -formatter stylish -exclude vendor/... ./...

dev: ## Install dev tools
	go install github.com/mgechev/revive

help: ## Display callable targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'