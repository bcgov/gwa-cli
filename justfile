set dotenv-load
set positional-arguments

install:
  go install -ldflags="-X 'main.ApiHost=${GWA_API_HOST}' -X 'main.ClientId=${GWA_CLIENT_ID}' -X 'main.ApiVersion=${GWA_VERSION}' -X 'main.Version=${CLI_VERSION}'"
  echo "gwa-cli installed in bin"

build:
  go build -o gwa -ldflags="-X 'main.ApiHost=${GWA_API_HOST}' -X 'main.ClientId=${GWA_CLIENT_ID}' -X 'main.ApiVersion=${GWA_VERSION}' -X 'main.Version=${CLI_VERSION}'" main.go

release:
  GWA_API_HOST=${GWA_API_HOST} GWA_CLIENT_ID=${GWA_CLIENT_ID} GWA_VERSION=${GWA_VERSION} CLI_VERSION=${CLI_VERSION} goreleaser release --snapshot --clean

run *args:
  go run -ldflags="-X 'main.ApiHost=$GWA_API_HOST' -X 'main.ClientId=$GWA_CLIENT_ID' -X 'main.ApiVersion=$GWA_VERSION' -X 'main.Version=${CLI_VERSION}'" main.go {{args}}

docs:
  go run build/gen-docs.go > docs/gwa-commands.md

test:
  go test ./...
