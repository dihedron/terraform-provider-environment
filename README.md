# Terraform Environment Binding Plugin 

[![CircleCI](https://circleci.com/gh/dihedron/terraform-provider-environment.svg?style=svg)](https://circleci.com/gh/dihedron/terraform-provider-environment)


## Installation

You can easily install the latest version with the following :

```
go get -u github.com/dihedron/terraform-provider-environment
```

Then add the plugin to your local `.terraformrc` :

```
cat >> ~/.terraformrc <<EOF
providers {
    environment = "${GOPATH}/bin/terraform-provider-environment"
}
EOF
```

## Provider example

```
provider "environment" {
    bindings = [{ 
             name    = "production"      
             url     = "http://www.example.com/environments/production?format=json",
             format  = "json"
        }, {
             name    = "certification"      
             url     = "http://www.example.com/environments/certification?format=toml",
             format  = "toml"
        }, {
             name    = "quality"      
             url     = "http://www.example.com/environments/quality?format=text",
             format  = "text"
        }, {
             name    = "integration"      
             url     = "http://www.example.com/environments/integration?format=toml",
             format  = "text"
        },
    ]
}
```

## Resource "Environment Binding" example

```
resource "environment_binding" "my" {
    name        = "certification",
    variables   = [{
	    name    = "PATH",
	    force   = false,
	    default = "/bin:/usr/bin"
	}, {
	    name    = "CLASSPATH",
	    force   = true,
	    default = ".:/opt/java/classes"
	}
    ]
}	 
```

The name of the binding must correspond to one of the bindings declared in the "environment" provider; 
the resource available at the binding URL must be a plain text file, with the same format you would use in
shell scripts, e.g.
```
# this is a comment
VAR1=value1
var2=value2
Var3=value3
# and so on
```
Future versions may support additional formats such as TOML, YAML and JSON.