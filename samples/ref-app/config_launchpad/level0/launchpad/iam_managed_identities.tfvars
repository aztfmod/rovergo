managed_identities = {
  level0 = {
    # Used by the release agent to deploy resources to level 0
    name               = "runner-level-0"
    resource_group_key = "security"
    tags = {
      level = "level0"
    }
  }
  level1 = {
    # Used by the release agent to deploy resources to level 1
    name               = "runner-level-1"
    resource_group_key = "security"
    tags = {
      level = "level1"
    }
  }
  level2 = {
    # Used by the release agent to deploy resources to level 2
    name               = "runner-level-2"
    resource_group_key = "security"
    tags = {
      level = "level2"
    }
  }
  level3 = {
    # Used by the release agent to deploy resources to level 3
    name               = "runner-level-3"
    resource_group_key = "security"
    tags = {
      level = "level3"
    }
  }
  level4 = {
    # Used by the release agent to deploy resources to level 4
    name               = "runner-level-4"
    resource_group_key = "security"
    tags = {
      level = "level4"
    }
  }
}
