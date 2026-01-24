---
page_title: "Kerberos Authentication Guide"
subcategory: ""
description: |-
  Guide to configuring Kerberos authentication for secure WinRM connections.
---

# Kerberos Authentication Guide

This guide explains how to configure Kerberos authentication for the Windows AD provider, which is the recommended and most secure method for authenticating to Active Directory.

## Why Kerberos?

Kerberos is the recommended authentication method because:

- **No password over the wire**: Kerberos uses tickets, not passwords
- **Mutual authentication**: Both client and server verify each other's identity
- **Replay attack protection**: Each ticket is time-limited and single-use
- **Delegation support**: Safely delegate credentials when needed (credential passing)
- **Industry standard**: Native to Active Directory environments

### Deprecated Methods

| Method | Status | Security Concerns |
|--------|--------|-------------------|
| **Kerberos** | ✅ Recommended | None |
| **NTLM** | ⚠️ Deprecated | Vulnerable to relay attacks, weak hashing |
| **Basic** | ⚠️ Deprecated | Credentials sent in cleartext |

NTLM authentication (`winrm_use_ntlm = true`) will be removed in v0.2.0.

## Prerequisites

### WinRM HTTPS Listener

Your Windows server must have an HTTPS WinRM listener configured:

```powershell
# On the Windows server (Domain Controller or management server)
# Check existing listeners
winrm enumerate winrm/config/listener

# Create HTTPS listener with a certificate
$cert = Get-ChildItem -Path Cert:\LocalMachine\My | Where-Object { $_.Subject -like "*$env:COMPUTERNAME*" }
New-WSManInstance -ResourceURI winrm/config/Listener -SelectorSet @{Address="*";Transport="HTTPS"} -ValueSet @{CertificateThumbprint=$cert.Thumbprint}

# Or use a self-signed certificate for testing
$cert = New-SelfSignedCertificate -DnsName $env:COMPUTERNAME -CertStoreLocation Cert:\LocalMachine\My
New-WSManInstance -ResourceURI winrm/config/Listener -SelectorSet @{Address="*";Transport="HTTPS"} -ValueSet @{CertificateThumbprint=$cert.Thumbprint}
```

### Firewall Rules

Ensure port 5986 (HTTPS) is open:

```powershell
New-NetFirewallRule -Name "WinRM-HTTPS" -DisplayName "WinRM HTTPS" -Enabled True -Direction Inbound -Protocol TCP -LocalPort 5986 -Action Allow
```

## Linux/macOS Configuration

### krb5.conf

Create `/etc/krb5.conf` (or a custom path):

```ini
[libdefaults]
    default_realm = YOURDOMAIN.COM
    dns_lookup_realm = false
    dns_lookup_kdc = false
    ticket_lifetime = 24h
    renew_lifetime = 7d
    forwardable = true
    rdns = false
    default_ccache_name = FILE:/tmp/krb5cc_%{uid}

[realms]
    YOURDOMAIN.COM = {
        kdc = dc01.yourdomain.com:88
        admin_server = dc01.yourdomain.com:749
    }

[domain_realm]
    .yourdomain.com = YOURDOMAIN.COM
    yourdomain.com = YOURDOMAIN.COM
```

Replace:
- `YOURDOMAIN.COM` with your AD domain (uppercase)
- `dc01.yourdomain.com` with your domain controller FQDN

### Provider Configuration

```hcl
provider "windowsad" {
  winrm_hostname = "dc01.yourdomain.com"
  winrm_username = "admin@YOURDOMAIN.COM"  # UPN format recommended
  winrm_password = var.ad_password
  winrm_proto    = "https"                  # Required for security
  winrm_port     = 5986                     # HTTPS port
  winrm_insecure = true                     # Set false with valid certs

  krb_realm = "YOURDOMAIN.COM"              # Uppercase domain name
  krb_conf  = "/etc/krb5.conf"              # Path to krb5.conf

  # Optional: for double-hop scenarios
  winrm_pass_credentials = true
  domain_controller      = "dc01.yourdomain.com"
}
```

## Using Keytab Files

For service accounts or automated pipelines, use a keytab instead of passwords:

### Generate Keytab

On a Windows machine with RSAT tools:

```powershell
# Create keytab for service account
ktpass /out terraform.keytab /princ terraform@YOURDOMAIN.COM /mapuser terraform@yourdomain.com /pass * /crypto AES256-SHA1 /ptype KRB5_NT_PRINCIPAL
```

Or on Linux with MIT Kerberos:

```bash
# Using ktutil
ktutil
addent -password -p terraform@YOURDOMAIN.COM -k 1 -e aes256-cts-hmac-sha1-96
wkt terraform.keytab
quit
```

### Provider Configuration with Keytab

```hcl
provider "windowsad" {
  winrm_hostname = "dc01.yourdomain.com"
  winrm_username = "terraform"              # Just the username
  winrm_proto    = "https"
  winrm_port     = 5986

  krb_realm  = "YOURDOMAIN.COM"
  krb_conf   = "/etc/krb5.conf"
  krb_keytab = "/path/to/terraform.keytab"  # Path to keytab file
}
```

## Service Principal Names (SPN)

By default, the provider constructs the SPN as `HTTP/<hostname>`. For custom SPNs:

```hcl
provider "windowsad" {
  # ... other settings ...

  krb_spn = "HTTP/winrm.yourdomain.com@YOURDOMAIN.COM"
}
```

Check registered SPNs on the server:

```powershell
setspn -L <computername>
```

## Troubleshooting

### Common Errors

#### "KDC unreachable" or timeout

- Verify KDC hostname resolves: `nslookup dc01.yourdomain.com`
- Check port 88 is reachable: `nc -zv dc01.yourdomain.com 88`
- Verify krb5.conf has correct KDC

#### "Clock skew too great"

Kerberos requires synchronized time (within 5 minutes):

```bash
# Check time difference
ntpdate -q dc01.yourdomain.com

# Sync time
sudo ntpdate -s dc01.yourdomain.com
```

#### "Cannot find KDC for realm"

- Verify krb5.conf realm is uppercase
- Check krb5.conf path is correct in provider config
- Enable debug: `export KRB5_TRACE=/dev/stderr`

#### "Pre-authentication failed"

- Verify username format: `user@DOMAIN.COM` (domain uppercase)
- Check password is correct
- Verify account is not locked/disabled

### Debug Logging

Enable Kerberos debug logging:

```bash
export KRB5_TRACE=/dev/stderr
export TF_LOG=DEBUG
terraform plan
```

### Verify Kerberos Manually

Test authentication before Terraform:

```bash
# Get ticket
kinit admin@YOURDOMAIN.COM

# List tickets
klist

# Test WinRM (with curl)
curl --negotiate -u : https://dc01.yourdomain.com:5986/wsman
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `WINDOWSAD_HOSTNAME` | WinRM server hostname |
| `WINDOWSAD_USER` | Username (UPN format recommended) |
| `WINDOWSAD_PASSWORD` | Password |
| `WINDOWSAD_PORT` | WinRM port (default: 5986) |
| `WINDOWSAD_PROTO` | Protocol: https (default) |
| `WINDOWSAD_KRB_REALM` | Kerberos realm (uppercase domain) |
| `WINDOWSAD_KRB_CONF` | Path to krb5.conf |
| `WINDOWSAD_KRB_KEYTAB` | Path to keytab file |
| `WINDOWSAD_KRB_SPN` | Custom service principal name |

## Security Best Practices

1. **Always use HTTPS** (`winrm_proto = "https"`)
2. **Use Kerberos** (`krb_realm` configured)
3. **Never commit credentials** - use environment variables or vault
4. **Use keytabs** for CI/CD pipelines instead of passwords
5. **Rotate credentials** regularly
6. **Limit service account permissions** to only required AD operations

## Migration from NTLM

If you're currently using NTLM (`winrm_use_ntlm = true`):

1. Set up WinRM HTTPS listener on server
2. Create krb5.conf on your Terraform runner
3. Update provider configuration:

```hcl
# Before (insecure)
provider "windowsad" {
  winrm_hostname = "server.domain.com"
  winrm_username = "admin"
  winrm_password = "secret"
  winrm_use_ntlm = true  # ⚠️ Deprecated!
}

# After (secure)
provider "windowsad" {
  winrm_hostname = "server.domain.com"
  winrm_username = "admin@DOMAIN.COM"
  winrm_password = "secret"
  winrm_proto    = "https"
  winrm_port     = 5986
  krb_realm      = "DOMAIN.COM"
  krb_conf       = "/etc/krb5.conf"
}
```

4. Test with `terraform plan`
5. Remove `winrm_use_ntlm` from configuration

## See Also

- [Provider Configuration](../index.md)
- [Migration from hashicorp/ad](./migration-from-hashicorp-ad.md)
- [Microsoft: WinRM Kerberos Configuration](https://docs.microsoft.com/en-us/windows/win32/winrm/authentication-for-remote-connections)
