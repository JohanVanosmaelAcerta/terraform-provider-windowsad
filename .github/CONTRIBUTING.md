# Contributing to terraform-provider-windowsad

Thank you for your interest in contributing! This provider is a maintained fork of the archived HashiCorp terraform-provider-ad.

## Getting Started

1. Fork and clone the repository
2. Ensure you have Go 1.25+ installed
3. Run `make build` to verify the build works

## Development Workflow

### Building

```bash
make build      # Build the provider
go install      # Install to $GOPATH/bin
```

### Testing

```bash
# Unit tests (no AD required)
go test ./... -v

# Acceptance tests (requires AD environment)
TF_ACC=1 go test ./windowsad -v -timeout 120m
```

### Code Quality

```bash
make fmt        # Format code
make lint       # Run linter
make vet        # Run go vet
```

## Submitting Changes

1. Create a feature branch: `git checkout -b feat/my-feature`
2. Make your changes with tests
3. Ensure `make fmt && make lint && go test ./...` passes
4. Commit using [conventional commits](https://www.conventionalcommits.org/):
   - `feat: add new attribute`
   - `fix: handle edge case`
   - `docs: update examples`
5. Push and open a Pull Request against `main`

## Guidelines

- See [AGENTS.md](../AGENTS.md) for detailed coding standards and patterns
- Check [GitHub Issues](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/issues) for planned work
- Discuss significant changes in an issue before starting work

## Local Provider Installation

For testing with Terraform locally:

```bash
make build
terraform init -plugin-dir /path/to/terraform-provider-windowsad
```

Or create a symlink in your [plugins directory](https://developer.hashicorp.com/terraform/cli/config/config-file#provider-installation).
