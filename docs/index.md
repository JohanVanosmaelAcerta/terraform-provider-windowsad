---
layout: ""
page_title: "Provider: Windows AD (Active Directory)"
description: |-
  The Windows AD provider provides resources to interact with an AD domain controller.
---

# Windows AD (Active Directory) Provider

The Windows AD provider provides resources to interact with an AD domain controller.

## About This Provider

This provider is a **community continuation** of the archived [HashiCorp terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad). When HashiCorp archived the original provider, this fork was created to:

- **Maintain the WinRM/PowerShell approach** — Continue using remote PowerShell execution via WinRM
- **Implement missing features** — Add AD and GPO PowerShell module capabilities not in the original
- **Merge community contributions** — Incorporate bug fixes from unmerged upstream PRs
- **Improve security** — Enforce Kerberos authentication over HTTPS

See the [GitHub repository](https://github.com/JohanVanosmaelAcerta/terraform-provider-windowsad) for the full roadmap and heritage from the original provider.

## Requirements

- Windows Server 2012R2 or greater.
- WinRM enabled with HTTPS listener (recommended).
- Kerberos authentication configured (required for non-Windows clients).
- Appropriate AD permissions (see below).

## Required Permissions

The service account used by the provider needs specific Active Directory permissions depending on which resources you manage:

| Resource | Required Permission |
|----------|---------------------|
| `windowsad_user` | Create/modify users in target OU |
| `windowsad_group` | Create/modify groups in target OU |
| `windowsad_group_membership` | Modify group membership |
| `windowsad_ou` | Create/modify OUs in target container |
| `windowsad_computer` | Create/modify computer accounts in target OU |
| `windowsad_gpo` | **Group Policy Creator Owners** membership |
| `windowsad_gplink` | Link GPOs to OUs (requires GPO and OU permissions) |
| `windowsad_gpo_security` | Modify GPO security settings |

### GPO Permissions

Creating and managing Group Policy Objects requires membership in the **Group Policy Creator Owners** group or equivalent permissions:

```powershell
# Add service account to Group Policy Creator Owners
Add-ADGroupMember -Identity "Group Policy Creator Owners" -Members "svc-terraform"
```

Without this permission, GPO operations will fail with:
```
An operations error occurred. (Exception from HRESULT: 0x80072020)
```

### Recommended: Delegated Administration

For production use, create a dedicated service account with delegated permissions rather than using Domain Admins:

```powershell
# Example: Delegate OU management to service account
$ou = "OU=Terraform,DC=example,DC=com"
$account = "EXAMPLE\svc-terraform"

# Grant full control over the OU and child objects
dsacls $ou /G "${account}:GA" /I:T
```

## Security Defaults (v0.1.0+)

This provider defaults to secure settings:

| Setting | Default | Notes |
|---------|---------|-------|
| `winrm_proto` | `https` | Secure by default |
| `winrm_port` | `5986` | HTTPS port |
| Authentication | Kerberos | Required for non-Windows clients |

### Deprecated Features

| Feature | Status | Removal |
|---------|--------|---------|
| `winrm_proto = "http"` | ⚠️ Deprecated | v0.2.0 |

See the [Kerberos Authentication Guide](guides/kerberos-authentication.md) for secure configuration.

## Migration from hashicorp/ad

This provider supports both `windowsad_*` (recommended) and legacy `ad_*` resource names for easy migration. See the [Migration Guide](guides/migration-from-hashicorp-ad.md) for details.

## Kerberos Authentication (Recommended)

Kerberos is the recommended authentication method. Set `krb_realm` to enable it:

```hcl
provider "windowsad" {
  winrm_hostname = "dc01.yourdomain.com"
  winrm_username = "admin@YOURDOMAIN.COM"
  winrm_password = var.password
  krb_realm      = "YOURDOMAIN.COM"
  krb_conf       = "/etc/krb5.conf"  # Optional: custom krb5.conf
}
```

If no `krb_conf` is supplied, the provider generates a minimal configuration using `krb_realm` and `winrm_hostname`.

For detailed setup instructions, see the [Kerberos Authentication Guide](guides/kerberos-authentication.md).

## Double-Hop Authentication

Starting with version 0.4.3 it is possible to point the provider to a host other than a Domain Controller and perform
all the management tasks through that host. Here is an example of the provider config:
```
provider "windowsad" {
  winrm_hostname         = "10.0.0.1"
  winrm_username         = var.username
  winrm_password         = var.password
  krb_realm              = "YOURDOMAIN.COM"
  krb_conf               = "${path.module}/krb5.conf"
  krb_spn                = "winserver1"
  winrm_port             = 5986
  winrm_proto            = "https"
  winrm_pass_credentials = true
}
```

In this case krb5.conf would look like this:
```
[libdefaults]
   default_realm = YOURDOMAIN.COM
   dns_lookup_realm = false
   dns_lookup_kdc = false


[realms]
	YOURDOMAIN.COM = {
		kdc 	= 	172.16.12.109
        admin_server = 172.16.12.109
		default_domain = YOURDOMAIN.COM
	}

[domain_realm]
    .kerberos.server = YOURDOMAIN.COM
	.yourdomain.com = YOURDOMAIN.COM
	yourdomain.com = YOURDOMAIN.COM
	yourdomain = YOURDOMAIN.COM
```


 A few things to note:
    - Double Hop Authentication is only enabled when using https
    - Authentication between management host and DC is done via Kerberos
    - The [AD Powershell module](https://docs.microsoft.com/en-us/powershell/module/activedirectory/?view=winserver2012r2-ps) as well as the [Group Policy Powershell Module](https://docs.microsoft.com/en-us/powershell/module/grouppolicy/?view=windowsserver2019-ps) is expected to be installed
      on the server before running the provider.


## Note about Local execution (Windows only)

It is possible to execute commands locally if the OS on which terraform is running is Windows.
In such case, your need to put the following settings in the provider configuration :

- Set winrm_username to null
- Set winrm_password to null
- Set winrm_hostname to null

Note: it will set to local only `if all 3 parameters are set to null`

### Example
```terraform
provider "windowsad" {
  winrm_hostname = ""
  winrm_username = ""
  winrm_password = ""
}
```

 ## Example Usage

```terraform
variable "hostname" { default = "dc01.yourdomain.com" }
variable "username" { default = "admin@YOURDOMAIN.COM" }
variable "password" { default = "password" }

// Recommended: Kerberos authentication with HTTPS (default)
provider "windowsad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  krb_realm      = "YOURDOMAIN.COM"
  krb_conf       = "/etc/krb5.conf"
  winrm_insecure = true  # Set false with valid certificates
}

// Kerberos with keytab (for CI/CD pipelines)
provider "windowsad" {
  winrm_hostname = var.hostname
  winrm_username = "terraform"
  krb_realm      = "YOURDOMAIN.COM"
  krb_conf       = "/etc/krb5.conf"
  krb_keytab     = "/path/to/terraform.keytab"
}

// Double-hop: WinRM to management server, Kerberos to DC
provider "windowsad" {
  winrm_hostname         = "mgmt.yourdomain.com"
  winrm_username         = var.username
  winrm_password         = var.password
  krb_realm              = "YOURDOMAIN.COM"
  krb_conf               = "/etc/krb5.conf"
  krb_spn                = "mgmt"
  winrm_pass_credentials = true
  domain_controller      = "dc01.yourdomain.com"
}

// Local execution (Windows only)
provider "windowsad" {
  winrm_hostname = ""
  winrm_username = ""
  winrm_password = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `winrm_hostname` (String) The hostname of the server we will use to run powershell scripts over WinRM. (Environment variable: WINDOWSAD_HOSTNAME)
- `winrm_password` (String) The password used to authenticate to the server's WinRM service. (Environment variable: WINDOWSAD_PASSWORD)
- `winrm_username` (String) The username used to authenticate to the server's WinRM service. (Environment variable: WINDOWSAD_USER)

### Optional

- `domain_controller` (String) Use a specific domain controller. (default: none, environment variable: WINDOWSAD_DC)
- `krb_conf` (String) Path to kerberos configuration file. (default: none, environment variable: WINDOWSAD_KRB_CONF)
- `krb_keytab` (String) Path to a keytab file to be used instead of a password
- `krb_realm` (String) The name of the kerberos realm (domain) we will use for authentication. (default: "", environment variable: WINDOWSAD_KRB_REALM)
- `krb_spn` (String) Alternative Service Principal Name. (default: none, environment variable: WINDOWSAD_KRB_SPN)
- `winrm_insecure` (Boolean) Trust unknown certificates. (default: false, environment variable: WINDOWSAD_WINRM_INSECURE)
- `winrm_pass_credentials` (Boolean) Pass credentials in WinRM session to create a System.Management.Automation.PSCredential. (default: false, environment variable: WINDOWSAD_WINRM_PASS_CREDENTIALS)
- `winrm_port` (Number) The port WinRM is listening for connections. (default: 5986, environment variable: WINDOWSAD_PORT)
- `winrm_proto` (String) The WinRM protocol we will use. (default: https, environment variable: WINDOWSAD_PROTO). Note: HTTP is deprecated.
