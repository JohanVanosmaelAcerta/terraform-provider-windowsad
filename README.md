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

If you're migrating from the archived HashiCorp provider:

1. Update your provider source from `hashicorp/ad` to `JohanVanosmaelAcerta/windowsad`
2. Rename all resources from `ad_*` to `windowsad_*`
3. Update environment variables from `AD_*` to `WINDOWSAD_*`

## Community Bug Fixes Included

This fork includes fixes from the following upstream PRs:

- #173: Fix custom_attributes hyphen/number issues
- #166: Permit empty group membership
- #159: Remove leaf objects on computer delete (recursive delete)
- #156: Use slash as delimiter instead of underscore
- #128: Fix cannot_change_password state detection
- #124: Fix multiple AD user creation
- #197: Fix password special character escaping

## Contributing

We welcome contributions! Please [create an issue](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/issues/new) to discuss changes before submitting a PR.

## License

[Mozilla Public License v2.0](LICENSE)
