package windowsad

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider exports the provider schema
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"winrm_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_USER", nil),
				Description: "The username used to authenticate to the server's WinRM service. (Environment variable: WINDOWSAD_USER)",
				//lintignore: V013
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					os := runtime.GOOS
					if v == "" && os != "windows" {
						errs = append(errs, fmt.Errorf("%q is allowed to be empty only if terraform runs on windows, (current os: %q) ", key, os))
					}
					return
				},
			},
			"winrm_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_PASSWORD", nil),
				Description: "The password used to authenticate to the server's WinRM service. (Environment variable: WINDOWSAD_PASSWORD)",
			},
			"winrm_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_HOSTNAME", nil),
				Description: "The hostname of the server we will use to run powershell scripts over WinRM. (Environment variable: WINDOWSAD_HOSTNAME)",
				//lintignore: V013
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					os := runtime.GOOS
					if v == "" && os != "windows" {
						errs = append(errs, fmt.Errorf("%q is allowed to be empty only if terraform runs on windows, (current os: %q) ", key, os))
					}
					return
				},
			},
			"winrm_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_PORT", 5986),
				Description: "The port WinRM is listening for connections. (default: 5986 for HTTPS, environment variable: WINDOWSAD_PORT)",
			},
			"winrm_proto": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_PROTO", "https"),
				Description: "The WinRM protocol we will use. (default: https, environment variable: WINDOWSAD_PROTO). Note: HTTP is deprecated for security reasons.",
			},
			"winrm_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_WINRM_INSECURE", false),
				Description: "Trust unknown certificates. (default: false, environment variable: WINDOWSAD_WINRM_INSECURE)",
			},
			"krb_realm": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_KRB_REALM", ""),
				Description: "The name of the kerberos realm (domain) we will use for authentication. (default: \"\", environment variable: WINDOWSAD_KRB_REALM)",
			},
			"krb_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_KRB_CONF", ""),
				Description: "Path to kerberos configuration file. (default: none, environment variable: WINDOWSAD_KRB_CONF)",
			},
			"krb_spn": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_KRB_SPN", ""),
				Description: "Alternative Service Principal Name. (default: none, environment variable: WINDOWSAD_KRB_SPN)",
			},
			"krb_keytab": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_KRB_KEYTAB", ""),
				Description: "Path to a keytab file to be used instead of a password",
			},
			"winrm_pass_credentials": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_WINRM_PASS_CREDENTIALS", false),
				Description: "Pass credentials in WinRM session to create a System.Management.Automation.PSCredential. (default: false, environment variable: WINDOWSAD_WINRM_PASS_CREDENTIALS)",
			},
			"domain_controller": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINDOWSAD_DC", ""),
				Description: "Use a specific domain controller. (default: none, environment variable: WINDOWSAD_DC)",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			// Primary names (recommended)
			"windowsad_user":     dataSourceADUser(),
			"windowsad_group":    dataSourceADGroup(),
			"windowsad_gpo":      dataSourceADGPO(),
			"windowsad_computer": dataSourceADComputer(),
			"windowsad_ou":       dataSourceADOU(),
			// Legacy aliases for migration from hashicorp/ad provider
			// Deprecated: Use windowsad_* names for new configurations
			"ad_user":     dataSourceADUser(),
			"ad_group":    dataSourceADGroup(),
			"ad_gpo":      dataSourceADGPO(),
			"ad_computer": dataSourceADComputer(),
			"ad_ou":       dataSourceADOU(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// Primary names (recommended)
			"windowsad_user":             resourceADUser(),
			"windowsad_group":            resourceADGroup(),
			"windowsad_group_membership": resourceADGroupMembership(),
			"windowsad_gpo":              resourceADGPO(),
			"windowsad_gpo_security":     resourceADGPOSecurity(),
			"windowsad_computer":         resourceADComputer(),
			"windowsad_ou":               resourceADOU(),
			"windowsad_gplink":           resourceADGPLink(),
			// Legacy aliases for migration from hashicorp/ad provider
			// Deprecated: Use windowsad_* names for new configurations
			"ad_user":             resourceADUser(),
			"ad_group":            resourceADGroup(),
			"ad_group_membership": resourceADGroupMembership(),
			"ad_gpo":              resourceADGPO(),
			"ad_gpo_security":     resourceADGPOSecurity(),
			"ad_computer":         resourceADComputer(),
			"ad_ou":               resourceADOU(),
			"ad_gplink":           resourceADGPLink(),
		},
		ConfigureFunc: initProviderConfig,
	}
}

func initProviderConfig(d *schema.ResourceData) (interface{}, error) {
	proto := d.Get("winrm_proto").(string)
	hostname := d.Get("winrm_hostname").(string)
	if strings.ToLower(proto) == "http" && hostname != "localhost" && hostname != "127.0.0.1" {
		log.Println("[WARN] Using HTTP protocol for WinRM is insecure and deprecated. " +
			"Credentials are transmitted in cleartext (base64 encoded). Please use HTTPS (winrm_proto = \"https\").")
	}

	// Warning for non-Windows platforms without Kerberos
	krbRealm := d.Get("krb_realm").(string)
	if runtime.GOOS != "windows" && krbRealm == "" {
		log.Println("[WARN] Running on non-Windows without krb_realm configured. " +
			"Kerberos authentication is required for non-Windows clients.")
	}

	cfg, err := config.NewConfig(d)
	if err != nil {
		return nil, err
	}
	pcfg := config.NewProviderConf(cfg)
	return pcfg, nil
}

func suppressCaseDiff(k, old, new string, d *schema.ResourceData) bool {
	// k is ignored here, but wee need to include it in the function's
	// signature in order to match the one defined for DiffSuppressFunc
	return strings.EqualFold(old, new)
}
