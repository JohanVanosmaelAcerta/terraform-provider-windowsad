<a href="https://terraform.io">
    <img src="https://github.com/hashicorp/terraform-provider-azurerm/raw/main/.github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for Windows Active Directory

[![Releases](https://img.shields.io/github/release/JohanVanosmaelAcerta/terraform-provider-windowsad.svg)](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/releases)
[![LICENSE](https://img.shields.io/github/license/JohanVanosmaelAcerta/terraform-provider-windowsad.svg)](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/blob/main/LICENSE)

This Windows AD provider for Terraform allows you to manage users, groups and group policies in your AD installation.

This is a maintained fork of the archived [HashiCorp terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad), with community bug fixes applied and ongoing development.

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 1.0+
* [Windows Server](https://www.microsoft.com/en-us/windows-server) 2012R2 or greater 
* [Go](https://golang.org/doc/install) version 1.25+ (for development)

## Getting Started

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    windowsad = {
      source  = "JohanVanosmaelAcerta/windowsad"
      version = "~> 1.0"
    }
  }
}

provider "windowsad" {
  winrm_hostname = "dc01.example.com"
  winrm_username = "admin@example.com"
  winrm_password = var.ad_password
}
```

Review the [docs](docs/) folder to understand which configuration options are available. You can find examples in [our examples folder](examples/).

## Migrating from hashicorp/ad

This provider supports both `windowsad_*` and legacy `ad_*` resource names, making migration seamless.

**Quick migration** (existing `ad_*` resources continue to work):

1. Update provider source from `hashicorp/ad` to `JohanVanosmael/windowsad`
2. Rename provider block from `ad` to `windowsad`
3. Migrate Terraform state (re-import or edit state file)
4. Run `terraform plan` — should show no changes

```hcl
# Your existing ad_user, ad_group, etc. resources work without modification!
resource "ad_user" "example" {
  display_name     = "John Doe"
  sam_account_name = "jdoe"
  # ...
}
```

For detailed instructions, see the **[Migration Guide](docs/guides/migration-from-hashicorp-ad.md)**.

> **Note:** The `ad_*` prefix is deprecated and will be removed in a future major version. Use `windowsad_*` for new configurations.

## Community Bug Fixes Included

This fork includes fixes from the following upstream PRs:

- #173: Fix custom_attributes hyphen/number issues
- #166: Permit empty group membership
- #159: Remove leaf objects on computer delete (recursive delete)
- #156: Use slash as delimiter instead of underscore
- #128: Fix cannot_change_password state detection
- #124: Fix multiple AD user creation
- #197: Fix password special character escaping

## Development

### Running Acceptance Tests

Acceptance tests require a Windows runner with access to an Active Directory environment. The runner machine requires specific configuration:

#### Windows Runner Configuration

1. **Windows Defender Exclusions** - Newly compiled Go test binaries are blocked by real-time scanning:

```powershell
# Run as Administrator on the runner machine
Add-MpPreference -ExclusionPath "D:\actions-runner-terraform"
Add-MpPreference -ExclusionPath "C:\Users\s-gmsa-gha$\go"
Add-MpPreference -ExclusionPath "C:\Users\s-gmsa-gha$\AppData\Local\go-build"
Add-MpPreference -ExclusionPath "C:\Users\s-gmsa-gha$\AppData\Local\Temp"
```

2. **Windows Developer Mode** - Terraform plugin testing framework creates symlinks, which requires the `SeCreateSymbolicLinkPrivilege`:

```powershell
# Run as Administrator on the runner machine
reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" /t REG_DWORD /f /v "AllowDevelopmentWithoutDevLicense" /d "1"
```

Or enable via: Settings → For developers → Developer Mode

3. **Reboot** the runner after applying these changes.

#### Environment Variables

The acceptance tests require these environment variables:

| Variable | Description |
|----------|-------------|
| `WINDOWSAD_HOSTNAME` | Domain controller hostname for WinRM |
| `WINDOWSAD_USER` | AD admin username (**without** `@realm` suffix) |
| `WINDOWSAD_PASSWORD` | AD admin password |
| `WINDOWSAD_KRB_REALM` | Kerberos realm (uppercase, e.g., `EXAMPLE.COM`) |
| `TF_VAR_ad_domain_name` | AD domain name (e.g., `example.com`) |
| `TF_VAR_ad_user_container` | OU for test users |
| `TF_VAR_ad_group_container` | OU for test groups |
| `TF_VAR_ad_computer_container` | OU for test computers |

> **Important:** For Kerberos authentication, `WINDOWSAD_USER` must be just the username (e.g., `svc-terraform`), not `svc-terraform@EXAMPLE.COM`. The realm is passed separately via `WINDOWSAD_KRB_REALM`.

## Contributing

We welcome contributions! Please [create an issue](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/issues/new) to discuss changes before submitting a PR.

## License

[Mozilla Public License v2.0](LICENSE)
