landingzone = {
  backend_type        = "azurerm"
  global_settings_key = "launchpad"
  level               = "level2"
  key                 = "foundations"
  tfstates = {
    caf_foundations = {
      level   = "lower"
      tfstate = "foundations.tfstate"
    }
    launchpad = {
      level   = "lower"
      tfstate = "launchpad.tfstate"
    }
  }
}
