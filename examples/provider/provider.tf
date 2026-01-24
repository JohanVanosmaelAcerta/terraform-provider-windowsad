variable "hostname" { default = "ad.yourdomain.com" }
variable "username" { default = "user" }
variable "password" { default = "password" }

// Recommended: Kerberos authentication with HTTPS
provider "windowsad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  winrm_proto    = "https"
  winrm_port     = 5986
  krb_realm      = "YOURDOMAIN.COM"
}

// Kerberos authentication with krb5.conf file
provider "windowsad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  winrm_proto    = "https"
  winrm_port     = 5986
  krb_realm      = "YOURDOMAIN.COM"
  krb_conf       = "/etc/krb5.conf"
}

// Kerberos with self-signed certificate (insecure for testing)
provider "windowsad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  winrm_proto    = "https"
  winrm_port     = 5986
  winrm_insecure = true
  krb_realm      = "YOURDOMAIN.COM"
}

// Local execution (Windows only)
provider "windowsad" {
  winrm_hostname = ""
  winrm_username = ""
  winrm_password = ""
}
