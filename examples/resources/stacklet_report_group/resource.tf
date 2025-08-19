data "stacklet_binding" "b1" {
  name = "binding-1"
}

data "stacklet_binding" "b2" {
  name = "binding-2"
}

resource "stacklet_report_group" "example" {
  name     = "example"
  schedule = "0 12 * * *"
  bindings = [
    data.stacklet_binding.b1.uuid,
    data.stacklet_binding.b2.uuid,
  ]
  group_by = ["account", "region"]
}
