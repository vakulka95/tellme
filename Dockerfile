FROM golang:1.14.4-alpine3.12 as builder

RUN echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen

RUN apk update \
 && apk upgrade \
 && apk add ca-certificates gcc make libc-dev musl-dev binutils git curl jq \
 && rm -rf /var/cache/apk/*

# golang
ENV GOPATH=/go
ENV GOROOT=/usr/local/go

ENV APP_NAME=tellme.api
ENV APP_VERSION=v1.0.0
ENV APP_BIN=/usr/local/bin/${APP_NAME}
ENV APP_PKG=gitlab.com/tellmecomua/${APP_NAME}
ENV APP_DIR=${GOPATH}/src/${APP_PKG}

WORKDIR ${APP_DIR}

COPY .git/ .git/
COPY app/ app/
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/

RUN go version \
 && export GO111MODULE=off \
 && export BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
 && export BUILD_COMMIT=$(git rev-parse HEAD) \
 && go build -o ${APP_BIN} -i -v \
     -ldflags="-linkmode external -extldflags '-static' \
     -X ${APP_PKG}/app/version.Created=${BUILD_DATE} \
     -X ${APP_PKG}/app/version.GitHash=${BUILD_COMMIT} \
     -X ${APP_PKG}/app/version.Version=${APP_VERSION}" \
     ./cmd/server/main.go

FROM alpine:3.12 as application

RUN echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen

RUN apk update \
 && apk upgrade \
 && apk add tzdata ca-certificates \
 && rm -rf /var/cache/apk/*

ENV APP_NAME=tellme.api
ENV APP_BIN=/usr/local/bin/${APP_NAME}
ENV MIGRATION_FILES_DIR=/etc/${APP_NAME}/migrations

COPY static/ /usr/share/${APP_NAME}/static
COPY app/persistence/_migrations/* ${MIGRATION_FILES_DIR}/
COPY --from=builder ${APP_BIN} ${APP_BIN}

WORKDIR /

CMD ${APP_BIN}