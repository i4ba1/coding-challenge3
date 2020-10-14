package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"github.com/i4ba1/CustomerOrderAPI/helper"
	"github.com/i4ba1/CustomerOrderAPI/user"
	"log"
	"math/rand"
	"net/http"  // used to access the request and response object of the api
	"os"        // used to read the environment variable
	_ "strconv" // package used to covert string into int type
	"strings"
	"time"

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateUser create a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty input of type models.User
	//var input user.User
	input := &user.UserDto{}

	// decode the json request to input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		ErrorResponse(http.StatusUnprocessableEntity, "Invalid JSON", w)
		return
	}

	if ok, errors := helper.ValidateInputs(*input); !ok {
		validationResponse(errors, w)
	}

	_, err = getUser(input.Email, input.PhoneNumber)
	if err != nil {
		ErrorResponse(http.StatusConflict, "Email or Phone Number already used", w)
	}

	// call insert input function and pass the input
	insertID := insertUser(input)

	response := make(map[string]interface{})
	response["message"] = "Email is " + input.Email + " with customer id "+insertID
	SuccessRespond(response,w)
}

// get one user from the DB by its userid
func getUser(email string, phoneNumber string) (user.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a u of models.User type
	var u user.User

	// create the select sql query
	sqlStatement := `SELECT * FROM tbl_customer WHERE email=$1 or phone_number=$2`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, email, phoneNumber)

	// unmarshal the row object to u
	err := row.Scan(&u.CustomerId, &u.Email, &u.PhoneNumber)

	switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return u, nil
		case nil:
			return u, nil
		default:
			log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty u on error
	return u, err
}

func generateCustomerId() string{
	rand.Seed(time.Now().Unix())
	//Only lowercase
	charSet := "abcdedfghijklmnopqrstvwxyz"
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	fmt.Println(output.String())
	return output.String()
}


// insert one user in the DB
func insertUser(user *user.UserDto) string {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	customerId := generateCustomerId()
	sqlStatement := "INSERT INTO tbl_customer (customer_id, customer_name, phone_number, email, dob, sex, salt, password, created_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING customer_id"

	// the inserted id will store in this id
	//var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	config := &helper.PasswordConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
	password, _ := helper.GeneratePassword(config, generateCustomerId())

	err := db.QueryRow(sqlStatement, customerId, user.CustomerName, user.PhoneNumber, user.Email,
		user.DateOfBird, user.Sex, password["salt"], password["password"], user.CreatedAt).Scan(&customerId)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", customerId)

	// return the inserted id
	return customerId
}

func SuccessRespond(fields map[string]interface{}, writer http.ResponseWriter){
	fields["status"] = "success"
	message,err := json.Marshal(fields)
	if err != nil{
		//An error occurred processing the json
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("An error occurred internally"))
	}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type","application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(message)
}

func ErrorResponse(statusCode int, error string, writer http.ResponseWriter){
	//Create a new map and fill it
	fields := make(map[string]interface{})
	fields["status"] = "error"
	fields["message"] = error
	message,err := json.Marshal(fields)

	if err != nil{
		//An error occurred processing the json
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("An error occured internally"))
	}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type","application/json")
	writer.WriteHeader(statusCode)
	writer.Write(message)
}

func validationResponse(fields map[string][]string, writer http.ResponseWriter) {
	//Create a new map and fill it
	response := make(map[string]interface{})
	response["status"] = "error"
	response["message"] = "validation error"
	response["errors"] = fields
	message, err := json.Marshal(response)

	if err != nil {
		//An error occurred processing the json
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("An error occured internally"))
	}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnprocessableEntity)
	writer.Write(message)
}