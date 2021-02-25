# immudb-homelab

This is a Go application that makes use of the immudb's SDK to 
insert records into the database.

You can specify the number of records you'd like to add by passing
a valid integer:

$ ./immudb-homelab 200

The application will generate a list of random key-value pairs and 
pass it to the immudb client for insertion. It uses the VerifiedSet
methods and prints the measured time it takes to insert the entries.

This utility is good for testing immudb with small samples. However, 
if you'd like to perform operations with large (10k+) entrires, this
utility might take a little longer, as it is not taking advantage of
immudb's SetAll batch operations. This might be added in the future. 

