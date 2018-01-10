# update.go
This file contains methods that modify deployment attributes through the kubernetes client-go. This is currently under development, and all instructions for completed methods will be found here.

### replicaUpdate
This method updates the replica count of a specific deployment. The string argument `metaname` refers to the name of the deployment specified in the metadata. The string argument `magnitude` specifies how much the replica count will increase by. If a negative number is given, then the replica count will be decreased by `magnitude`. For an example, refer to the main function.
Note that there is currently 2 versions of the method. The first is used out-of-cluster while the second with `v1` is used in-cluster. Only the out-of-cluster method has been verified.

### main
This method will conain unit tests for each of the methods described above. Feel free to comment out any lines.
