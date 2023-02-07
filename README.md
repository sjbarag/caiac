> Infrastructure as code (IaC) tools allow you to manage infrastructure with
> configuration files rather than through a graphical user interface. IaC
> allows you to build, change, and manage your infrastructure in a safe,
> consistent, and repeatable way by defining resource configurations that you
> can version, reuse, and share.

— https://developer.hashicorp.com/terraform/tutorials/aws-get-started/infrastructure-as-code

Why should infrastructure get to have all the fun? We should also be able to
source code that we can version, reuse, and share in a safe, consistent, and
repeatable way. Introducing…

# Terraform Provider CaIaC
A Terraform provider enabling Code as Infrastructure as Code!

## How it Works
Instead of manually editing code through a graphical user interface, the CaIaC
Terraform Provider lets you manage your source code through configuration
files. No more arguing about code styling or installing Go tools for your
IDE: simply describe your code's structure and configuration in HCL and let
Terraform take care of the rest.

### Without CaIaC
Managing go files requires manually writing, maintaining, and understanding the
nuances of Go syntax. This gets increasingly complex as more languages are
introduced to a project. Consider the standard Go "hello world":

```go
package main

import "fmt"

func main() {
  fmt.Println("Hello, world")
}
```

### With CaIaC
With CaIaC, only HCL is required! Any language — not just Go — can be managed
by Terraform, freeing developers to think in the consistent, reusable interface
of HCL. Changes can be previewed before they're applied to your code,
preventing accidental modification or deletion of files.

With CaIaC, that same "hello world" doesn't require writing Go at all:

```hcl
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
```

### No Code
CaIaC can even power no-code workflows! Since Terraform can consistently
recreate files, there's no need to even commit source code any more. Simply
commit an `.hcl` file, add `*.go` to `.gitignore`, and leverage Terraform to
ensure everyone on your team can reliably create the same code.

## Common Questions
### Does this actually work?
Goodness, I wish it didn't.

### How complicated can my file get?
Not very. Most of the Go AST isn't modeled in the `caciac_go_source` Terraform
resource because this is a terrible idea. Also, Terraform doesn't allow a
resource's schema to contain a cycle, so a function call can't currently
contain an arbitrary expression.

### Why isn't this published to the Terraform Registry?
Because it's horrible.

## Development
### Build provider
Please don't. But if you need to:

```sh
cat <<EOF
provider_installation {
  dev_overrides {
    "registry.terraform.io/sjbarag/caiac" = "/some/directory/on/your/path"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF > ~/.terraformrc
```

Then run `go install .` from the root of this directory.

### Test sample configuration
I wouldn't recommend it, but okay:

```sh
# 1. Install the provider (see above)

# 2. Initialize a Terraform workspace
terraform -chdir=example init

# 3. View a plan
terraform -chdir=example plan

# 4. Apply the plan
terraform -chdir=example apply

# 5. View the produced file
cat ./example/main.go
```
