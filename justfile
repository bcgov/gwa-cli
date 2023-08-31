set dotenv-load
set positional-arguments

install:
  go install -ldflags="-X 'main.ApiHost=${GWA_API_HOST}' -X 'main.ClientId=${GWA_CLIENT_ID}'"
  echo "gwa-cli installed in bin"

build:
  go build -o gwa -ldflags="-X 'main.ApiHost=${GWA_API_HOST}' -X 'main.ClientId=${GWA_CLIENT_ID}'" main.go
  echo "gwa executable compiled"

run *args:
  go run -ldflags="-X 'main.ApiHost=$GWA_API_HOST' -X 'main.ClientId=$GWA_CLIENT_ID'" main.go {{args}}

docs:
  go run build/gen-docs.go

test:
  go test ./...
