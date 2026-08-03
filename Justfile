# Run all tests
test:
    go test ./...

# Compile-check by default; pass an output path to keep the binary (e.g. `just build bagre`)
build out="/dev/null":
    go build -o {{out}} .

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

# Auto‑format the entire project
format:
    golangci-lint fmt ./...

