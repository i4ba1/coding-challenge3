# How to
 - run docker-compose up --build
 - open up http://localhost:15432, to open pgadmin, login with password: "postgresql"
 - create new role in pgadmin, username: mni & password: mni123!# and set privileges to superuser
 - create new database with name "customer_db", after that right click -> Restore... and select the file customer.sql
 - open up terminal and run "go run main.go" or for unit test "go run main_test.go"
 
### I put the assignment logic in GolangLogicChallenge
### For query in AssignmentSQL