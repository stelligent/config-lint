module "test_module" {
   source = "./modules/foo"
   name = "counter_1"
   count = 4
}

module "test_module_with_description" {
   source = "./modules/foo"
   name = "counter_2"
   description = "here is a description"
}
