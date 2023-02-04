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
  contents = <<-EOT
    package foo
    import "fmt"
    func main() {
      fmt.Println("hello")
    }
  EOT
}

output "main-file" {
  value = data.caiac_go_source.main_go
}

output "foo-file" {
  value = caiac_go_source.foo_go
}
