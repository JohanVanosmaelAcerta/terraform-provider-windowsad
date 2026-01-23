# AGENTS.md

Instructions for AI agents working on this Terraform provider for Windows Active Directory.

## Project Overview

This is a **Go-based Terraform provider** that manages Windows Active Directory resources (users, groups, OUs, GPOs, computers) via WinRM/PowerShell. It's a maintained fork of the archived HashiCorp `terraform-provider-ad`.

## Architecture

```
├── main.go                    # Provider entrypoint
├── windowsad/                 # Provider implementation
│   ├── provider.go            # Provider schema and configuration
│   ├── resource_ad_*.go       # Resource implementations (CRUD operations)
│   ├── data_source_ad_*.go    # Data source implementations (read-only)
│   └── internal/
│       ├── config/            # Provider configuration
│       ├── winrmhelper/       # WinRM/PowerShell execution layer
│       ├── adschema/          # AD schema definitions
│       └── gposec/            # GPO security policy helpers
├── docs/                      # Terraform registry documentation
├── examples/                  # Example Terraform configurations
└── vendor/                    # Vendored Go dependencies
```

### Key Patterns

- **Resources**: Each AD object type has a `resource_ad_<type>.go` file implementing CRUD via `resourceAD<Type>Create/Read/Update/Delete` functions
- **Data Sources**: Read-only lookups in `data_source_ad_<type>.go` files
- **WinRM Layer**: All AD operations execute PowerShell scripts remotely via `internal/winrmhelper/`
- **Schema**: Uses `terraform-plugin-sdk/v2` with `schema.Resource` definitions

## Boundaries

### ✅ Always
- Run `make fmtcheck` before committing
- Run `go test ./...` to verify changes don't break tests
- Use vendored dependencies (`go mod vendor` after dependency changes)
- Assign new issues to a milestone

### ⚠️ Ask First
- Before modifying `provider.go` schema (affects all users)
- Before changing WinRM authentication methods
- Before adding new provider configuration options
- Before modifying `go.mod` dependencies

### 🚫 Never
- Commit credentials or secrets
- Modify files in `vendor/` directly (use `go mod vendor`)
- Push directly to `main` branch
- Remove existing resource attributes (breaking change)
- Use `-Properties *` in new AD queries (performance issue #27)

## Git Workflow

### Branch Naming
```
feat/INFRA-123-short-description
fix/INFRA-456-bug-summary
chore/update-dependencies
```

### Commit Messages
Use [Conventional Commits](https://www.conventionalcommits.org/):
```
feat(user): add name property support (#2)
fix(computer): handle transient AD errors (#33)
docs: update AGENTS.md with boundaries
chore: update to Go 1.25
```

### Pull Requests
- Always target `main` branch
- Use squash merge (default)
- Include detailed description with Summary, Changes, Testing sections
- Link to GitHub issue(s) being addressed

## Development Commands

```bash
# Build the provider
make build

# Run unit tests (no AD environment required)
go test ./... -v

# Run acceptance tests (requires AD environment)
TF_ACC=1 go test ./windowsad -v -timeout 120m

# Format code
make fmt          # gofmt
make fumpt        # gofumpt (stricter)

# Lint
make lint         # golangci-lint
make tflint       # terraform provider linter
make vet          # go vet

# Check dependencies
make depscheck    # go mod tidy + go mod vendor
```

## Coding Standards

### Go Style

- Go 1.25+ with vendored dependencies (`go mod vendor`)
- Format with `gofmt -s` (CI enforces this)
- Follow standard Go naming conventions (camelCase for unexported, PascalCase for exported)
- Use `log.Printf("[DEBUG]...")` for debug logging with appropriate level prefixes

### Terraform Provider Conventions

- Resource functions follow pattern: `resourceAD<Type>` returns `*schema.Resource`
- CRUD handlers: `resourceAD<Type>Create/Read/Update/Delete`
- Use `schema.TypeString`, `schema.TypeBool`, `schema.TypeList`, etc. for attribute types
- Always implement `Importer` with `StateContext: schema.ImportStatePassthroughContext`
- Use `DiffSuppressFunc` for case-insensitive DN comparisons

### Resource Implementation Pattern

```go
func resourceADExample() *schema.Resource {
    return &schema.Resource{
        Description: "`ad_example` manages Example objects in Active Directory.",
        Create:      resourceADExampleCreate,
        Read:        resourceADExampleRead,
        Update:      resourceADExampleUpdate,
        Delete:      resourceADExampleDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: map[string]*schema.Schema{
            // ... attribute definitions
        },
    }
}
```

### WinRM Helper Pattern

PowerShell commands are built and executed in `internal/winrmhelper/`:

```go
// Define struct with JSON tags matching PowerShell output
type Example struct {
    GUID string `json:"ObjectGUID"`
    Name string `json:"Name"`
}

// Execute PowerShell and unmarshal JSON response
func GetExample(client *config.ProviderConf, guid string) (*Example, error) {
    cmd := fmt.Sprintf("Get-ADExample -Identity %q | ConvertTo-Json", guid)
    result, err := RunWinRMCommand(client, cmd)
    // ... unmarshal and return
}
```

## Testing

### Unit Tests

- Located alongside source files (`*_test.go`)
- Run without AD environment: `go test ./...`
- Test helper functions and schema validation

### Acceptance Tests

- Prefix: `TestAcc*`
- Require `TF_ACC=1` environment variable
- Require live AD environment with these env vars:
  - `WINDOWSAD_HOSTNAME` - Domain controller hostname
  - `WINDOWSAD_USER` - AD admin username
  - `WINDOWSAD_PASSWORD` - AD admin password
- Use `resource.Test()` with `TestCase` and `TestStep` structs
- Always include import test steps

### Test Pattern

```go
func TestAccResourceADExample_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:  func() { testAccPreCheck(t, requiredEnvVars) },
        Providers: testAccProviders,
        CheckDestroy: testAccResourceADExampleDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccResourceADExampleConfig(),
                Check: resource.ComposeTestCheckFunc(
                    testAccResourceADExampleExists("ad_example.test"),
                ),
            },
            {
                ResourceName:      "ad_example.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}
```

## Documentation

- **Provider docs**: `docs/index.md`
- **Resource docs**: `docs/resources/<resource>.md`
- **Data source docs**: `docs/data-sources/<datasource>.md`
- **Examples**: `examples/resources/<resource>/` and `examples/data-sources/<datasource>/`

When adding/modifying resources:
1. Update schema descriptions (they appear in generated docs)
2. Add/update example in `examples/`
3. Validate examples: `make validate-examples`

## Common Pitfalls

1. **DN Case Sensitivity**: Active Directory DNs are case-insensitive; use `suppressCaseDiff` in schema
2. **Password Special Characters**: Escape special chars in PowerShell strings (see PR #197 fix)
3. **Empty Collections**: Handle nil vs empty slices carefully (see PR #166 fix)
4. **Recursive Delete**: Use `-Recursive` flag when deleting containers with children (see PR #159 fix)
5. **WinRM Timeouts**: Long-running operations may need timeout adjustments in provider config

## CI/CD

GitHub Actions workflows in `.github/workflows/`:

- **ci.yml**: Build and test on push/PR to main
- **release.yml**: GoReleaser for tagged releases
- **unit_tests.yaml**: Unit test matrix

Checks run on every PR:
- `go build ./...`
- `go test ./...`
- `gofmt` formatting check
- `go vet`

## Environment Variables

Provider configuration (for testing):

| Variable | Description |
|----------|-------------|
| `WINDOWSAD_HOSTNAME` | DC hostname for WinRM connection |
| `WINDOWSAD_USER` | Admin username |
| `WINDOWSAD_PASSWORD` | Admin password |
| `WINDOWSAD_PORT` | WinRM port (default: 5985) |
| `WINDOWSAD_PROTO` | Protocol: http/https (default: http) |
| `WINDOWSAD_KRB_REALM` | Kerberos realm for auth |

## Adding New Resources

1. Create `windowsad/resource_ad_<type>.go` with CRUD functions
2. Create `windowsad/internal/winrmhelper/winrm_<type>.go` for PowerShell operations
3. Register in `provider.go` ResourcesMap
4. Add tests in `windowsad/resource_ad_<type>_test.go`
5. Add documentation in `docs/resources/<type>.md`
6. Add example in `examples/resources/<type>/`
7. Run `make lint && make test` before committing

## Roadmap & Planned Work

See [GitHub Issues](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/issues) for the full backlog.

### Milestones

| Milestone | Focus |
|-----------|-------|
| **v0.1.0 - Initial Release** | Registry setup, security hardening (Kerberos-only) |
| **v0.2.0 - Bug Fixes** | Reliability, transient error handling, import fixes |
| **v0.3.0 - Enhanced Features** | New attributes, group improvements, OU protection |
| **v0.4.0 - Quality** | Test coverage, documentation, examples |
| **v0.5.0 - Advanced Features** | gMSA, bulk data sources, performance optimizations |

When creating or triaging issues, always assign to an appropriate milestone.

### Key Priorities

### P1 - Release Blockers
| Issue | Description |
|-------|-------------|
| #23 | Complete Terraform Registry setup and tag v0.1.0 |
| #32 | Security: Enforce Kerberos-only auth, deprecate NTLM/Basic |

### P2 - High Priority
| Issue | Description |
|-------|-------------|
| #31 | Set up self-hosted runner with AD DS for acceptance tests |
| #26 | Enable parallel resource reads in Terraform SDK |
| #25 | Implement batch queries for AD resources |
| #21 | Add OU `protected_from_accidental_deletion` support |
| #20 | Fix timeout issues on large AD operations |
| #18 | Add acceptance tests for all resources |
| #13 | Fix: user import fails to populate all attributes |
| #12 | Fix: user password change with special characters |
| #8 | Fix: OU with children cannot be deleted |
| #4 | Update WinRM library for Unicode encoding |
| #3 | Add group `managed_by` and `custom_attributes` |
| #2 | Add `name` property for AD user objects |
| #1 | Fix group membership permadiff |

### P3 - Enhancements
| Issue | Description |
|-------|-------------|
| #30 | Add unit tests for internal packages |
| #28 | Add bulk data sources for read-only queries |
| #27 | Use selective AD properties instead of `-Properties *` |
| #19 | Add HTTPS/TLS support for WinRM |
| #17 | Add comprehensive examples directory |
| #16 | Fix: group scope/category change requires recreation |
| #15 | Support multi-valued custom attributes |
| #14 | Add computer `description` attribute |
| #11 | Add `windowsad_service_account` for gMSA/sMSA |
| #10 | Add `windowsad_domain` data source |
| #9 | Fix: GPO import permissions error |
| #7 | Publish provider to Terraform Registry |

### Implementation Notes for Agents

When working on issues:
1. **Check related upstream PRs** - Many issues reference HashiCorp PRs that can be ported
2. **Run existing tests first** - Establish baseline before changes
3. **Security issues (#32)** - Coordinate deprecation warnings across versions
4. **Performance issues (#25-27)** - Profile before/after with large AD environments
5. **WinRM changes (#4, #19)** - Test with Unicode characters and TLS certificates
