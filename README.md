# Terraform Environment Binding Plugin 

- Website: https://www.terraform.io
- [![CircleCI](https://circleci.com/gh/dihedron/terraform-provider-environment.svg?style=svg)](https://circleci.com/gh/dihedron/terraform-provider-environment)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px">

## Installation

You can easily install the latest version with the following commands:

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
It might be  necessary to ensure that the `$GOPATH` variable is properly configured and `$GOPATH/bin` is accessible under the `$PATH`.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

## Provider declaration example

```
provider "environment" {
    # restrict provider version in 0.1.x
    version = "~> 0.1"
    # specify potential bindings, each with its own label and URL
    bindings = [{ 
             name    = "production"      
             url     = "http://www.example.com/environments/production?format=json"
        }, {
             name    = "certification"      
             url     = "http://www.example.com/environments/certification?format=toml"
        }, {
             name    = "quality"      
             url     = "http://www.example.com/environments/quality?format=text"
        }, {
             name    = "integration"      
             url     = "http://www.example.com/environments/integration?format=toml"
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

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/dihedron/terraform-provider-environment`

```sh
$ mkdir -p $GOPATH/src/github.com/dihedron; cd $GOPATH/src/github.com/dihedron
$ git clone git@github.com:dihedron/terraform-provider-environment
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/dihedron/terraform-provider-environment
$ make build
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory; make sure it is under your `$PATH` if you want the provider to be accessible from Terraform.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-environment
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests do not create real resources, so it is safe to run them at will.

```sh
$ make testacc
```
