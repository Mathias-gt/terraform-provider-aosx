# Terraform RESTCONF Provider

The Terraform AOS-X RESTCONF provider allows you to manage configuration of Alcatel-Lucent Enterprise OmniSwitch using the RESTCONF API.

## Building the Provider

1. Clone the provider repository:

```bash
git clone https://github.com/Mathias-gt/terraform-provider-aosx
```

2. Move into the repository directory:

```bash
cd terraform-provider-aosx
```

3. Build the provider binary:

```bash
go build -o terraform-provider-aosx
```

4. Move the provider binary to the Terraform plugins directory:

For Terraform 0.12.x and earlier, you can put the provider binary in the same directory as your Terraform configuration files or in the user plugins directory, which is usually `~/.terraform.d/plugins` on UNIX-like systems and `%APPDATA%\terraform.d\plugins` on Windows.

Example for Linux:
```bash
mkdir -p ~/.terraform.d/plugins/spacewalkers.com/ale/aosx/1.0.0/linux_amd64
mv -f terraform-provider-aosx ~/.terraform.d/plugins/spacewalkers.com/ale/aosx/1.0.0/linux_amd64
```

## Example of main.tf file

```hcl
terraform {
  required_providers {
    aosx = {
      source  = "spacewalkers.com/ale/aosx"
    }
  }
}

provider "aosx" {
  username = "admin"
  password = "switch"
}

  resource "aosx_restconf" "example" {
    path    = "https://172.26.9.40/restconf/data/ale-vlan:ale-vlan"
    delete_path    = "https://172.26.9.40/restconf/data/ale-vlan:ale-vlan/VLAN/VLAN_LIST=100"
    content = jsonencode({
            "ale-vlan:ale-vlan": {
                "VLAN": {
                    "VLAN_LIST": [
                        {
                            "vlan_id": 100,
                            "description": "VLAN_100",
                            "admin_status": "enable"
                        }
                    ]
                }
            }
        })
    }
```