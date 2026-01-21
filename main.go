package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: windowsad.Provider})
}
