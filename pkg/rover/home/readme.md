# Rover Home Defaults

The files in this directory are used to create the default rover home directory which is in `$HOME/.rover`

At runtime if the rover home directory doesn't exist it will be created when rover starts, and the contents of this directory used to populate it

Note 1. The contents of this directory are embedded into the rover binary at build time, it has purpose at runtime.

Note 2. This readme is not copied to the rover home directory 