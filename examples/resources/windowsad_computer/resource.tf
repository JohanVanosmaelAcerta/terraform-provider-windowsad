variable "name" { default = "test" }
variable "pre2kname" { default = "TEST" }

resource "windowsad_computer" "c" {
  name      = var.name
  pre2kname = var.pre2kname
}
