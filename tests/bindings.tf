
provider "environment" {
    # restrict provider version in 0.1.x
#    version = "~> 1.0.0"
#    version = "1.0.0"
    # specify potential bindings, each with its own label and URL
    environments = [{ 
             name    = "production"      
#             url     = "http://www.example.com/environments?binding=production"
            url = "http://localhost:8000/variables.txt"
        }, {
             name    = "certification"      
             url     = "http://www.example.com/environments?binding=certification"
        }, {
             name    = "quality"      
             url     = "http://www.example.com/environments?binding=quality"
        }, {
             name    = "integration"      
             url     = "http://www.example.com/environments?binding=integration"
        },
    ]
}

data "environment_bindings" "my" {
	name = "production",
    filters = [{
            name = "PATH"
            override = true
            default = "/home/andrea/bin"
        },
    ]
}


output "body" {
  value = "${data.environment_bindings.my.name}"
}


