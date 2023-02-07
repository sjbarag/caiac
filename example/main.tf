terraform {
  required_providers {
    caiac = {
      source = "registry.terraform.io/sjbarag/caiac"
    }
  }
}

provider "caiac" {}

resource "caiac_go_source" "main" {
  filename = "./main.go"
  package_name = "main"
  import {
    path = "fmt"
  }

  func {
    name = "main"

    signature{}

    body{
      statement {
        kind = "expression"
        expression {
          kind = "call"
          call {
            func {
              from = "fmt"
              prop = "Println"
            }
            arg {
              kind = "string"
              value = "Hello, world"
            }
          }
        }
      }
    }
  }
}
