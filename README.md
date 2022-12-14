# Klaabu

Klaabu is an IP network documenting tool, usually referred to as an IPAM (IP Address Management) tool, with a very strong focus on simplicity, convenience, and operating in Cloud-Native environments.

## Why invent something new?

There are a ton of different IPAM tools available out there, but with Klaabu you don’t have to <i>host</i> and <i>maintain</i> yet another tool. NetBox, for instance, is a great tool, but in case NetBox is your source-of-truth, you will need to have some sort of an SLA, strong observability, DR plan, and more.
Klaabu tries to solve this unnecessary TOIL by using a file-based schema that you can store in Git and a powerful `validate` function to make sure your schema is in the right state. The validate function can be used in your CI pipeline, or standalone, to make sure you don’t have duplicated, overlapping CIDRs, or some other schema mistakes. The combination of storing your schema in Git and the validate function offers exactly the kind of workflow you are already used to, with code reviews, audit trails, and more.

### Why another markup language?

The first couple of iterations of the Klaabu program were using a YAML-based schema. However, after using Klaabu in a live environment (with 100 VPCs, ~500 subnets, BGP ASNs, and more) we decided that this is not the right format. Using a YAML-based schema in such an environment will result in hundreds of lines, which becomes very hard to read. We believe that the ‘Klaabu Markup Language’ (`kml`) is much easier and more convenient to the user.

### Terraform

Apart from some other powerful (CLI) functions, Klaabu offers the `export-terraform` function which exports your schema to a Terraform supported module. In turn, you can import this module from any other Terraform module and lookup the CIDRs, attributes, or labels. Having a single source-of-truth that also works natively with Terraform is very convenient and will simplify the usage of your VPC, Subnet, SG, and other Terraform modules a lot.

## Installation

With GO installed:

```bash
go get github.com/transferwise/klaabu
```

Alternatively, you can also import the package directly:

```bash
import “github.com/transferwise/klaabu”
```

In case you want to build the package yourself

```bash
make build
```

Last but not least, you can also just download the binary directly from the release page.

## CLI usage

```
Usage: klaabu <command> [args]
```

At the moment the Klaabu CLI supports the following commands:
* `find`: recursively search the schema for any object that matches your search pattern.
* `get`: in contrast to the `find` command, `get` retrieves a single object in the schema.
* `list`: the `list` command shows all the child objects of a certain instance.
* `space`: use this command to see the available IP space within a certain prefix/object.
* `init`: initializes a new schema.
* `validate`: validates your schema, including the actual content of the objects (e.g. valid CIDR, no overlapping CIDRs, and more).
* `fmt`: used to rewrite Klaabu configuration files to a canonical format and style.
* `export-terraform`: exports your schema to a valid Terraform module; your schema is stored in a file with the `tf.json` notation, which is a valid input module for Terraform.

### Examples

The examples assume your schema lives in the current working directory / you have an environment variable set pointing to the location of your schema.

```bash
klaabu space 192.168.0.0/20
```

```bash
klaabu find -label az=euc1-az1

klaabu find -label vpc=foobar,env=production
```

## Workflow

- Create a new private repository and store your `schema.kml` in there
- Use your IDE to add a new CIDR to your schema
- Run `klaabu fmt` to produce configuration files that conform to the imposed style
- Run `klaabu validate` to make sure your schema is in a valid state and that you don’t have overlapping CIDRs for instance
- Follow your personal/company’s process for committing and merging your changes. You probably want to follow the traditional code review process with a CI pipeline that also uses the validate function.

## Contributing

Check out the [CONTRIBUTING](./CONTRIBUTING.md/) guide if you want to contribute.


## Acknowledgements

A big THANKS to **Taras Burko**, for mentoring me, and for all your time & effort in writing this program with me.

A big THANKS to **Taavi Tuisk** for mentoring & helping me these past years.
