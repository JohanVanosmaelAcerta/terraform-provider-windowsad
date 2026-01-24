variable name { default = "TestOU" }
variable path { default = "dc=yourdomain,dc=com" }
variable description { default = "some description" }
variable protected { default = false }
variable container { default = "CN=Users,dc=yourdomain,dc=com" }

variable name { default = "test group" }
variable sam_account_name { default = "TESTGROUP" }
variable scope { default = "global" }
variable category { default = "security" }

resource "windowsad_group" "g" {
  name             = var.name
  sam_account_name = var.sam_account_name
  scope            = var.scope
  category         = var.category
  container        = var.container
}

resource windowsad_group "g2" {
    name             = "${var.name}-2"
    sam_account_name = "${var.sam_account_name}-2"
    container        = var.container
}


resource windowsad_user "u" {
    display_name     = "test user"
    principal_name   = "testUser"
    sam_account_name = "testUser"
    initial_password = "SuperSecure1234!!"
    container        = var.container
}

resource windowsad_group_membership "gm" {
    group_id = windowsad_group.g.id
    group_members  = [ windowsad_group.g2.id, windowsad_user.u.id ]
}
