resource "stacklet_configuration_profile_msteams" "example" {
  customer_config_input = {
    prefix = "myprefix"
    tags = {
      "app" = "stacklet-msteams-bot"
    }
  }
  access_config_input = {
    client_id        = "00000000-1111-2222-3333-444444444444"
    roundtrip_digest = "724ba7cc82663bc247b5a100b3ca2ece"
    tenant_id        = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  }

  channel_mapping {
    name       = "ch1"
    channel_id = "19:hZZSubNbJL7A5cYMGLnK_AiL3ytC2gNl6yF08_LVzbM1@thread.tacv2"
    team_id    = "ffffffff-eeee-dddd-cccc-bbbbbbbbbbbb"
  }

  channel_mapping {
    name       = "ch2"
    channel_id = "19:dFb954569d964adf94ag0481755ce4b9@thread.tacv2"
    team_id    = "99999999-9999-7777-6666-555555555555"
  }
}
