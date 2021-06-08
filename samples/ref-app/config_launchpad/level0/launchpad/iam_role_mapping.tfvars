#
# Services supported: subscriptions, storage accounts and resource groups
# Can assign roles to: AD groups, AD object ID, AD applications, Managed identities
#
role_mapping = {
  built_in_role_mapping = {
    subscriptions = {
      # Required both to create Azure resources on their respective levels.
      # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/managed_service_identity
      logged_in_subscription = {
        "Contributor" = {
          managed_identities = {
            keys = ["level0", "level1", "level2", "level3", "level4"]
          }
        }
        "User Access Administrator" = {
          managed_identities = {
            keys = ["level0", "level1", "level2", "level3", "level4"]
          }
        }
      }
    }
    # Access storage accounts without an access key and read to next level.
    # Separation of responsibility from subscription roles for best practices.
    storage_accounts = {
      level0 = {
        "Storage Blob Data Contributor" = {
          managed_identities = {
            keys = ["level0"]
          }
        }
        "Storage Blob Data Reader" = {
          managed_identities = {
            keys = ["level1"]
          }
        }
      }

      level1 = {
        "Storage Blob Data Contributor" = {
          managed_identities = {
            keys = ["level1"]
          }
        }
        "Storage Blob Data Reader" = {
          managed_identities = {
            keys = ["level2"]
          }
        }
      }

      level2 = {
        "Storage Blob Data Contributor" = {
          managed_identities = {
            keys = ["level2"]
          }
        }
        "Storage Blob Data Reader" = {
          managed_identities = {
            keys = ["level3"]
          }
        }
      }

      level3 = {
        "Storage Blob Data Contributor" = {
          managed_identities = {
            keys = ["level3"]
          }
        }
        "Storage Blob Data Reader" = {
          managed_identities = {
            keys = ["level4"]
          }
        }
      }

      level4 = {
        "Storage Blob Data Contributor" = {
          managed_identities = {
            keys = ["level4"]
          }
        }
      }
    }
  }
}
