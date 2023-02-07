terraform {
  required_providers {
    caiac = {
      source = "registry.terraform.io/sjbarag/caiac"
    }
  }
}

provider "caiac" {}

resource "caiac_go_source" "some_new_file" {
  filename = "./some_new_file.go"
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
