terraform {
  required_providers {
    caiac = {
      source = "registry.terraform.io/sjbarag/caiac"
    }
  }
}

provider "caiac" {}

data "caiac_source_go" "main_go" {
  filename = "main.go"
}

output "main-file" {
  value = data.caiac_source_go.main_go
}
