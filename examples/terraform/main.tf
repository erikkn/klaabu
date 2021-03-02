module "klaabu" {
  source = "./klaabu"
}

output "aws" {
  value = module.klaabu.aliases["aws"]
}
