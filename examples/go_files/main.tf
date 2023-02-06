terraform {
  required_providers {
    caiac = {
      source = "registry.terraform.io/sjbarag/caiac"
    }
  }
}

provider "caiac" {}

data "caiac_go_source" "main_go" {
  filename = "main.go"
}

resource "caiac_go_source" "foo_go" {
  filename = "./foo/foo.go"
  package_name = "main"
  import {
    path = "fmt"
  }
}

resource "caiac_go_source" "some_new_file" {
  filename = "./some_new_file.go"
  package_name = "main"
  import {
    path = "fmt"
    name = "format"
  }

  func {
    name = "foo"
    signature {
      param {
        name = "s"
        type = "string"
      }
      result {
        type = "string"
      }
    }
  }

  func {
    name = "main"
    signature{}
  }
}

output "result" {
  value = caiac_go_source.some_new_file.contents
}
