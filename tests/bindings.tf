
provider "environment" {
    # restrict provider version in 0.1.x
#    version = "~> 1.0.0"
#    version = "1.0.0"
    # specify potential bindings, each with its own label and URL
    bindings = [{ 
             name    = "production"      
             url     = "http://www.example.com/environments?binding=production"
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

data "environment" "http_test" {
	name = "production"
  #url = "http://www.google.com"
}

/*
output "body" {
  value = "${data.http.http_test.body}"
}

data "environment" 
*/
