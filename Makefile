run: ## Run API in watch mode
	CONFIG_PATH=env.yaml wgo run cmd/aircraft_api/main.go

build: ## Build API
	CONFIG_PATH=env.yaml go build cmd/aircraft_api/main.go

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
