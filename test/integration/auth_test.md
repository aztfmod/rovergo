Enumeration of test cases for authentication.  

| func | Identities | Perms | Username | Priority | Status |
| ---- | ---------- | ----- | -------- | -------- | ------ |
| Test_VM_SystemAssigned_No_Role | system-assigned MI | none | N/A | low | done |
| Test_VM_SystemAssigned_SubOwner_Role | system-assigned MI | subscription owner | N/A | high | done |
| *** | user principal | subscription owner | userid/pw | high | - |
| *** | service principal | subscription owner | client-id | high | - |
| *** | 1 user-assigned MI | none | client-id | high | - |
| *** | 1 user-assigned MI | none | object-id | low | - |
| *** | 1 user-assigned MI | none | resource-id | low | - |
| *** | 1 user-assigned MI | subscription owner | client-id | low | - |
| *** | 1 user-assigned MI | subscription owner | object-id | low | - |
| *** | 1 user-assigned MI | subscription owner | resource-id | low | - |
| *** | 1 user-assigned MI | rg owner | client-id | low | - |
| *** | 1 user-assigned MI | rg owner | object-id | low | - |
| *** | 1 user-assigned MI | rg owner | resource-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | none | client-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | none | object-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | none | resource-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | subscription owner | client-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | subscription owner | object-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | subscription owner | resource-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | rg owner | client-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | rg owner | object-id | low | - |
| *** | system-assigned MI + 1 user-assigned MI | rg owner | resource-id | low | - |
| *** | >1 user-assigned MI | none | client-id | low | - |
| *** | >1 user-assigned MI | none | object-id | low | - |
| *** | >1 user-assigned MI | none | resource-id | low | - |
| *** | >1 user-assigned MI | subscription owner | client-id | low | - |
| *** | >1 user-assigned MI | subscription owner | object-id | low | - |
| *** | >1 user-assigned MI | subscription owner | resource-id | low | - |
| *** | >1 user-assigned MI | rg owner | client-id | low | - |
| *** | >1 user-assigned MI | rg owner | object-id | low | - |
| *** | >1 user-assigned MI | rg owner | resource-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | none | client-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | none | object-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | none | resource-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | subscription owner | client-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | subscription owner | object-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | subscription owner | resource-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | rg owner | client-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | rg owner | object-id | low | - |
| *** | system-assigned MI + >1 user-assigned MI | rg owner | resource-id | low | - |
