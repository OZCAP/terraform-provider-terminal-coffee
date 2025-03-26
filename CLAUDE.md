# Terminal Coffee Terraform Provider - Development Guide

## Build & Test Commands
```bash
make build                   # Build the provider
make install                 # Install provider locally
make test                    # Run all tests
go test -v ./terminal        # Run all tests in terminal package
go test -v ./terminal -run TestCreateOrder  # Run a specific test
go fmt ./...                 # Format Go code
```

## Code Style Guidelines
- **Imports**: Standard library first, then external packages
- **Formatting**: Use `go fmt` for standard Go formatting with tabs
- **Types**: Use explicit types for all struct fields
- **Struct Tags**: JSON tags in camelCase, Terraform fields in snake_case
- **Naming**: CamelCase for exported items, camelCase for unexported
- **Error Handling**: Use `fmt.Errorf()` with context, `diag.Diagnostics` for Terraform errors
- **Comments**: Document all exported functions, types, and fields
- **Testing**: Use table-driven tests with descriptive messages
- **Resource Schema**: Include type, description, and validation for all fields

## Project Requirements
- Go >= 1.20
- Terraform >= 0.13.x
- Terraform Plugin SDK v2.29.0