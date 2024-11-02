PATH := $(PATH):$(GOPATH)/bin
BOIL_VER := v4.16.2
BOIL_EXT_VER := v0.8.0

prepare-sqlboiler:
	@go get -u github.com/tiendc/sqlboiler-extensions@$(BOIL_EXT_VER)
	@go install github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@$(BOIL_VER)

gen-models:
	@sqlboiler mysql \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates $(GOPATH)/pkg/mod/github.com/tiendc/sqlboiler-extensions@$(BOIL_EXT_VER)/templates/boilv4/mysql

gen-gql:
	@go get github.com/99designs/gqlgen/codegen/config@v0.17.55
	@go get github.com/99designs/gqlgen/internal/imports@v0.17.55
	@go get github.com/99designs/gqlgen@v0.17.55
	@go run github.com/99designs/gqlgen

run-server:
	@godotenv -f /app/.env go run server.go

run-test:
	@godotenv -f /app/.env.test.local go test -v ./...

run-test-ci:
	@godotenv -f /app/.env.test go test -v ./...
