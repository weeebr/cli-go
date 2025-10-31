# Go CLI Build System

# Default target: build all tools
default: bin

# Build all tools to specified directory
bin BIN=~/git-helpers/bin:
  mkdir -p {{BIN}}
  go run build.go -o {{BIN}}

# Build only affected tools (incremental build)
affected BIN=~/git-helpers/bin:
  mkdir -p {{BIN}}
  go run build.go -o {{BIN}}

# Install using Go's native installation
install:
  go install ./...

# Clean build artifacts
clean BIN=~/git-helpers/bin:
  rm -f {{BIN}}/*

# Build individual groups
build-git BIN=~/git-helpers/bin:
  mkdir -p {{BIN}}
  go run build.go -o {{BIN}} --groups git

build-ai BIN=~/git-helpers/bin:
  mkdir -p {{BIN}}
  go run build.go -o {{BIN}} --groups ai

build-tools BIN=~/git-helpers/bin:
  mkdir -p {{BIN}}
  go run build.go -o {{BIN}} --groups tools

# Test all commands
test:
  @echo "Testing git commands..."
  @echo '{"branch":"main"}' | ./bin/gco
  @./bin/gname
  @./bin/gmain
  @./bin/grt
  @./bin/gstats
  @echo "All tests passed!"

# Show help
help:
  @echo "Available targets:"
  @echo "  bin        - Build all tools (default)"
  @echo "  affected   - Build only affected tools (incremental)"
  @echo "  install    - Install using go install"
  @echo "  clean      - Remove build artifacts"
  @echo "  test       - Test all commands"
  @echo "  help       - Show this help"
