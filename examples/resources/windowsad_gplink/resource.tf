resource "windowsad_ou" "o" {
  name        = "gplinktestOU"
  path        = "dc=yourdomain,dc=com"
  description = "OU for gplink tests"
  protected   = false
}

resource "windowsad_gpo" "g" {
  name        = "gplinktestGPO"
  domain      = "yourdomain.com"
  description = "gpo for gplink tests"
  status      = "AllSettingsEnabled"
}

resource "windowsad_gplink" "og" {
  gpo_guid  = windowsad_gpo.g.id
  target_dn = windowsad_ou.o.dn
  enforced  = true
  enabled   = true
}
