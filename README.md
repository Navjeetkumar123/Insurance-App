# Insurance-App
Blockchain aaplication
go to Insurance directory.
Run script ./runnApp.sh
Open Duplicate command line and run ./testApi.sh

Above two commands will do the network setup,chaincode installation and init.
To test the application please  import Insurance .postman_collection2 file in postman.
Before testing in postman do generate JWT token of both organisation and token generation collection is already defined in this file.
Use this token in Authorisation header for rest of the APIs with respect to organisation'
In this Network Org1 is Insurer and Org2 is Customer.
