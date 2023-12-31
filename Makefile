# .PHONY: migrate migrate_down migrate_up migrate_version docker prod docker_delve local swaggo test

# ==============================================================================
# Docker compose commands

develop:
	echo "Starting docker environment"
	docker-compose -f docker-compose.local.yml up --build

# ==============================================================================
# Tools commands

run-linter:
	echo "Starting linters"
	golangci-lint run ./... --timeout=20m0s

# ==============================================================================
# Main

run:
	go run ./cmd/main.go

build:
	go build ./cmd/main.go

test:
	go test -cover ./...

release:
	goreleaser release --clean


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache


# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)

# =============================================================================
# Profiling ( Requires pprof debug/middleware enabled and secondary debug server)

profile-web:
	go tool pprof -http=":8000" pprofbin http://localhost:5555