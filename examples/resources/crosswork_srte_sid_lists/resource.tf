resource "crosswork_srte_sid_lists" "example" {
  name = "cwtest-sid-list"
  sids = [
    {
      index = 1
      label = 16001
    }
  ]
}
