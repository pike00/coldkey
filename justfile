version := `git describe --tags --always --dirty 2>/dev/null || echo dev`
binary := "coldkey"
image := "ghcr.io/pike00/coldkey"

goflags := "-trimpath -ldflags=\"-s -w -X main.version=" + version + "\""
security_flags := "--network none --read-only --cap-drop ALL --security-opt no-new-privileges:true --tmpfs /tmp:rw,noexec,nosuid,size=10m"

# Build the binary
build:
    CGO_ENABLED=0 go build {{goflags}} -o {{binary}} ./cmd/coldkey

# Run tests
test:
    go test -race -count=1 ./...

# Run golangci-lint
lint:
    golangci-lint run ./...

# Run go vet
vet:
    go vet ./...

# Check formatting
fmt:
    gofumpt -l -d .

# Build Docker image
docker:
    docker build --build-arg VERSION={{version}} -t {{image}}:{{version}} -t {{image}}:latest .

# Run interactively in Docker (secure defaults)
docker-run:
    @mkdir -p output
    docker run --rm -it {{security_flags}} \
        -u $(id -u):$(id -g) \
        -v $(pwd)/output:/out \
        {{image}}:latest

# Create backup from existing key via Docker
docker-backup KEYFILE:
    @mkdir -p output
    docker run --rm {{security_flags}} \
        -u $(id -u):$(id -g) \
        -v {{KEYFILE}}:/keys/keys.txt:ro \
        -v $(pwd)/output:/out \
        {{image}}:latest backup -o /out/backup.html /keys/keys.txt

# Push Docker image to ghcr.io
push: docker
    docker push {{image}}:{{version}}
    docker push {{image}}:latest

# Tag a release and push (triggers CI release workflow)
release TAG:
    git tag {{TAG}}
    git push origin {{TAG}}

# Record asciinema demo
demo:
    ./demo/record.sh

# Remove build artifacts
clean:
    rm -f {{binary}}

# Full CI pipeline
ci: vet test build docker
