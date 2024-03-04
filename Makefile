# Load env from file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="mysql://${DSN}" up

## migrations/up-staging: apply all up database migrations to staging
.PHONY: migrations/up-staging
migrations/up-staging:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="mysql://${StagingDSN}" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="mysql://${DSN}" down

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="mysql://${DSN}" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="mysql://${DSN}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="mysql://${DSN}" version
