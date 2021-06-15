## Embedded custom action files

The custom_actions_src is a location for custom action '.yaml' files to be embedded into the go binary
The Rover application exports these files into the $Home/.rover/custom_actions directory at run time
if the directory does not already exist. 
This allows custom_actions for key external tools like tflint and terratest to be deployed automatically.