# Server parameters
UAT_SERVER_USER=root
UAT_SERVER_ADDRESS=tellme_uat

PROD_SERVER_USER=root
PROD_SERVER_ADDRESS=185.25.116.234

# Application parameters
SERVER_BIN_NAME=apiserver
MIGRATE_BIN_NAME=tellme-migrate

CONFIG_DIR=/etc/tellme/config
TELLME_STATIC_DIR=/usr/share/tellme

MIGRATIONS_DIR=/etc/tellme/migrations

# Build ldflags
VERSION = $(shell git tag --sort="v:refname" | tail -n1)
GITHASH = $(shell git rev-parse --short HEAD)

LDFLAGS += -X gitlab.com/tellmecomua/tellme.api/app/version.Version=$(VERSION)
LDFLAGS += -X gitlab.com/tellmecomua/tellme.api/app/version.GitHash=$(GITHASH)
LDFLAGS += -X gitlab.com/tellmecomua/tellme.api/app/version.Created=$(shell date +"%D-%T")

# message wrappers
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

deploy-uat-backend:	     ## deploy UAT backend files
	@echo -en "$(OK_COLOR)->> Deploying UAT cofnig files...$(NO_COLOR)"
	@scp ./app/config/.env.uat ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:${CONFIG_DIR}/.env

	@echo -en "$(OK_COLOR)->> Deploying migration files...$(NO_COLOR)"
	@scp -r ./app/persistence/_migrations/* ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:${MIGRATIONS_DIR}

	@echo -en "$(OK_COLOR)->> Deploying admin static files...$(NO_COLOR)"
	@scp -r ./static ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:${TELLME_STATIC_DIR}

	@echo -en "$(OK_COLOR)->> Deploying bin app...$(NO_COLOR)"
	@scp ./dist/apiserver ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:/usr/bin/tellme

	@echo -en "$(OK_COLOR)->> Deploying bin migrate...$(NO_COLOR)"
	@scp ./dist/tellme-migrate ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:/usr/bin/tellme-migrate

	@echo -en "$(OK_COLOR)->> Deploying nginx.conf file...$(NO_COLOR)"
	@scp ./dist/nginx.conf ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:/etc/nginx/conf.d/tellme.conf

	@echo -en "$(OK_COLOR)->> Deploying systemd tellme backend unit file...$(NO_COLOR)"
	@scp ./dist/tellme.service ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:/etc/systemd/system/

	@echo -en "$(OK_COLOR)->> Deploying systemd tellme migrate unit file...$(NO_COLOR)"
	@scp ./dist/tellme-migrate.service ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:/etc/systemd/system/

.PHONY: deploy-prod-backend

deploy-prod-backend:	## deploy PROD backend files
	@echo -en "$(OK_COLOR)->> Deploying PORD cofnig files...$(NO_COLOR)"
	@scp ./app/config/.env.prod ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:${CONFIG_DIR}/.env

	@echo -en "$(OK_COLOR)->> Deploying migration files...$(NO_COLOR)"
	@scp -r ./app/persistence/_migrations/* ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:${MIGRATIONS_DIR}

	@echo -en "$(OK_COLOR)->> Deploying admin static files...$(NO_COLOR)"
	@scp -r ./static ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:${TELLME_STATIC_DIR}

	@echo -en "$(OK_COLOR)->> Deploying bin app...$(NO_COLOR)"
	@scp ./dist/apiserver ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:/usr/bin/tellme

	@echo -en "$(OK_COLOR)->> Deploying bin migrate...$(NO_COLOR)"
	@scp ./dist/tellme-migrate ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:/usr/bin/tellme-migrate

	@echo -en "$(OK_COLOR)->> Deploying nginx.conf file...$(NO_COLOR)"
	@scp ./dist/nginx.conf ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:/etc/nginx/conf.d/tellme.conf

	@echo -en "$(OK_COLOR)->> Deploying systemd tellme backend unit file...$(NO_COLOR)"
	@scp ./dist/tellme.service ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:/etc/systemd/system/

	@echo -en "$(OK_COLOR)->> Deploying systemd tellme migrate unit file...$(NO_COLOR)"
	@scp ./dist/tellme-migrate.service ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:/etc/systemd/system/

.PHONY: deploy-prod-backend

deploy-uat-swagger:     ## deploy UAT swagger apidocs
	@echo -en "$(OK_COLOR)->> Deploying UAT swagger api docs...$(NO_COLOR)"
	@scp ./apidocs/swagger.json ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:${TELLME_STATIC_DIR}/apidocs
	@echo -e "$(OK_COLOR) Done! $(NO_COLOR)"

.PHONY: deploy-uat-swagger

deploy-prod-swagger:     ## deploy PROD swagger apidocs
	@echo -en "$(OK_COLOR)->> Deploying PROD swagger api docs...$(NO_COLOR)"
	@scp ./apidocs/swagger.json ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:${TELLME_STATIC_DIR}/apidocs
	@echo -e "$(OK_COLOR) Done! $(NO_COLOR)"

.PHONY: deploy-prod-swagger

deploy-uat-web:     ## deploy UAT web static
	@echo -en "$(OK_COLOR)->> Deploying web static files...$(NO_COLOR)"
	@scp -r ./dist/project/uat/* ${UAT_SERVER_USER}@${UAT_SERVER_ADDRESS}:/usr/share/nginx/html
	@echo -e "$(OK_COLOR) Done! $(NO_COLOR)"

.PHONY: deploy-uat-web

deploy-prod-web:     ## deploy PROD web static
	@echo -en "$(OK_COLOR)->> Deploying web static files...$(NO_COLOR)"
	@scp -r ./dist/project/prod/* ${PROD_SERVER_USER}@${PROD_SERVER_ADDRESS}:/usr/share/nginx/html
	@echo -e "$(OK_COLOR) Done! $(NO_COLOR)"

.PHONY: deploy-prod-web

build:          ## build backend app (API | Admin | Migrate)
	@echo -en "$(OK_COLOR)->> Building backend server...$(NO_COLOR)"
	@CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags "${LDFLAGS}" -a -installsuffix cgo -o ./dist/${SERVER_BIN_NAME} ./cmd/server/main.go

	@echo -en "$(OK_COLOR)->> Building migrate tool...$(NO_COLOR)"
	@CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ./dist/${MIGRATE_BIN_NAME} ./cmd/migrate/main.go

	@echo -e "$(OK_COLOR) Done! $(NO_COLOR)"

.PHONY: build

swagger:          ## generate new swagger .yaml spec
	@echo -en "$(OK_COLOR)->> Generating swagger spec...$(NO_COLOR)"
	@go run -mod=vendor cmd/swagger/main.go
	@echo -e "$(OK_COLOR) Done! $(NO_COLOR)"

.PHONY: swagger

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help