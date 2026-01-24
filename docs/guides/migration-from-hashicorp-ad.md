# Migrating from hashicorp/ad Provider

This guide explains how to migrate from the archived [HashiCorp terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad) to this actively maintained provider.

## Overview

The `windowsad` provider is a maintained fork of `hashicorp/ad` with community bug fixes and ongoing development. Migration is straightforward because:

1. **Dual resource prefix support**: Both `ad_*` and `windowsad_*` resource names work
2. **Compatible provider configuration**: Same attributes, just a different provider name
3. **No state migration required**: Resources can be imported or state can be manipulated

## Quick Migration (Recommended)

The fastest migration path requires only two changes to your Terraform configuration:

### Step 1: Update required_providers

```hcl
# Before
terraform {
  required_providers {
    ad = {
      source  = "hashicorp/ad"
      version = "~> 0.4"
    }
  }
}

# After
terraform {
  required_providers {
    windowsad = {
      source  = "JohanVanosmael/windowsad"
      version = "~> 0.1"
    }
  }
}
```

### Step 2: Update provider block

```hcl
# Before
provider "ad" {
  winrm_hostname = "dc01.example.com"
  winrm_username = var.ad_username
  winrm_password = var.ad_password
}

# After
provider "windowsad" {
  winrm_hostname = "dc01.example.com"
  winrm_username = var.ad_username
  winrm_password = var.ad_password
}
```

### Step 3: Keep existing resources (no changes needed!)

Your existing `ad_user`, `ad_group`, etc. resources continue to work:

```hcl
# This still works with the windowsad provider!
resource "ad_user" "example" {
  display_name     = "John Doe"
  sam_account_name = "jdoe"
  principal_name   = "jdoe@example.com"
  container        = "OU=Users,DC=example,DC=com"
}

resource "ad_group" "admins" {
  name             = "App Admins"
  sam_account_name = "app-admins"
  container        = "OU=Groups,DC=example,DC=com"
  scope            = "global"
  category         = "security"
}
```

### Step 4: Handle Terraform state

After updating your configuration, you need to migrate the Terraform state. Choose one of these approaches:

#### Option A: Re-import resources (safest)

```bash
# Remove old resource from state
terraform state rm ad_user.example

# Import into same resource name (ad_* still works)
terraform import ad_user.example <GUID-or-DN>

# Or import into new name if you want to migrate names too
terraform import windowsad_user.example <GUID-or-DN>
```

#### Option B: Manipulate state directly

```bash
# Export current state
terraform state pull > terraform.tfstate.backup

# Edit the state file to change provider references
# Change: "provider": "provider[\"registry.terraform.io/hashicorp/ad\"]"
# To:     "provider": "provider[\"registry.terraform.io/JohanVanosmael/windowsad\"]"

# Push modified state
terraform state push terraform.tfstate
```

#### Option C: Use terraform state mv (if renaming resources)

If you're also renaming resources from `ad_*` to `windowsad_*`:

```bash
terraform state mv ad_user.example windowsad_user.example
terraform state mv ad_group.admins windowsad_group.admins
```

### Step 5: Verify migration

```bash
terraform init -upgrade
terraform plan
```

A successful migration shows **no changes** (resources already exist in AD).

## Environment Variables

If you use environment variables for provider configuration, update them:

| Old (hashicorp/ad) | New (windowsad) |
|-------------------|-----------------|
| `AD_HOSTNAME` | `WINDOWSAD_HOSTNAME` |
| `AD_USER` | `WINDOWSAD_USER` |
| `AD_PASSWORD` | `WINDOWSAD_PASSWORD` |
| `AD_PORT` | `WINDOWSAD_PORT` |
| `AD_PROTO` | `WINDOWSAD_PROTO` |
| `AD_KRB_REALM` | `WINDOWSAD_KRB_REALM` |

## Resource Name Mapping

Both prefixes work identically. Use `windowsad_*` for new configurations:

| Legacy (deprecated) | Recommended |
|--------------------|-------------|
| `ad_user` | `windowsad_user` |
| `ad_group` | `windowsad_group` |
| `ad_group_membership` | `windowsad_group_membership` |
| `ad_computer` | `windowsad_computer` |
| `ad_ou` | `windowsad_ou` |
| `ad_gpo` | `windowsad_gpo` |
| `ad_gpo_security` | `windowsad_gpo_security` |
| `ad_gplink` | `windowsad_gplink` |

Data sources follow the same pattern:

| Legacy (deprecated) | Recommended |
|--------------------|-------------|
| `data.ad_user` | `data.windowsad_user` |
| `data.ad_group` | `data.windowsad_group` |
| `data.ad_computer` | `data.windowsad_computer` |
| `data.ad_ou` | `data.windowsad_ou` |
| `data.ad_gpo` | `data.windowsad_gpo` |

## Gradual Migration

You can migrate resources gradually:

```hcl
# Old resources (still work)
resource "ad_user" "legacy_user" {
  # ...
}

# New resources (recommended for new additions)
resource "windowsad_user" "new_user" {
  # ...
}
```

## Bug Fixes Included

By migrating, you automatically get these community bug fixes:

- **#197**: Password special character escaping
- **#173**: Custom attributes with hyphens/numbers
- **#166**: Empty group membership handling
- **#159**: Recursive delete for computers with child objects
- **#156**: Slash delimiter instead of underscore
- **#128**: `cannot_change_password` state detection
- **#124**: Multiple AD user creation race condition

## Deprecation Notice

The `ad_*` prefix is provided for migration convenience and will be deprecated in a future major version. We recommend migrating to `windowsad_*` names when convenient.

## Troubleshooting

### "Provider not found" error

Ensure you've run `terraform init -upgrade` after changing the provider source.

### State shows resources will be destroyed

The state still references the old provider. Use one of the state migration options above.

### Kerberos authentication issues

The `windowsad` provider has improved Kerberos support. Ensure:
- `WINDOWSAD_USER` contains just the username (not `user@REALM`)
- `WINDOWSAD_KRB_REALM` is set to the uppercase realm
- For HTTPS/credential passing, set `WINDOWSAD_PROTO=https` and `WINDOWSAD_WINRM_PASS_CREDENTIALS=true`

## Getting Help

- [GitHub Issues](https://github.com/JohanVanosmael/terraform-provider-windowsad/issues)
- [Provider Documentation](../index.md)
