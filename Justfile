ci:
    go test ./...


# Run all tests
test:
    "go test -timeout 180s ./...";

# Build a binary 
build:
    go build .

# Run the full quality suite: format → build → lint → test
quality:
    @echo "▶ format..."
    just format
    @echo "▶ build.."
    just build 
    @echo "▶ lint"
    just lint 
    @echo "▶ test"
    just test 
    @echo "✅ quality suite passed"

# Run linting 
# TODO we need to configure the lint file yet 
lint:
    golangci-lint run ./...

# Auto‑format & fix lint issues for the entire project
format:
    golangci-lint run --fix ./... && golangci-lint fmt ./...

