package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/i4ba1/CustomerOrderAPI/helper"
	"github.com/i4ba1/CustomerOrderAPI/user"
	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
	"github.com/twinj/uuid"
	"log"
	"math/rand"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"
	_ "strconv" // package used to covert string into int type
	"strings"
	"time"
)

var client *redis.Client
func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

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
	response["message"] = "Email is " + input.Email + " with customer id " + insertID
	SuccessRespond(response, w)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&mapToken); err != nil {
		ErrorResponse(http.StatusUnprocessableEntity, err.Error(), w)
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		fmt.Println("the error: ", err)
		ErrorResponse(http.StatusUnauthorized, "Refresh token expired", w)
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		ErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			ErrorResponse(http.StatusUnprocessableEntity, err.Error(), w)
			return
		}

		//userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		var userId string
		userId = claims["user_id"].(string)
		if len(userId) <= 0 {
			ErrorResponse(http.StatusUnprocessableEntity, "Error occurred", w)
			return
		}

		//Delete the previous Refresh Token
		deleted, delErr := deleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			ErrorResponse(http.StatusUnauthorized, "unauthorized", w)
			return
		}

		//Create new pairs of refresh and access tokens
		ts, createErr := createAccessAndRefreshToken(userId)
		if  createErr != nil {
			ErrorResponse(http.StatusForbidden, createErr.Error(), w)
			return
		}
		//save the tokens metadata to redis
		saveErr := createAuth(userId, ts)
		if saveErr != nil {
			ErrorResponse(http.StatusForbidden, saveErr.Error(), w)
			return
		}

		response := make(map[string]interface{})
		/*tokens := map[string]interface{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}*/
		response["access_token"] = ts.AccessToken
		response["refresh_token"] = ts.RefreshToken
		SuccessRespond(response, w)
	} else {
		ErrorResponse(http.StatusUnauthorized, "refresh expired", w)
	}
}

func createAuth(userid string, td *user.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty input of type models.User
	//var input user.User
	input := &user.LoginDto{}
	// decode the json request to input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		ErrorResponse(http.StatusUnprocessableEntity, "Invalid JSON", w)
		return
	}

	if ok, errors := helper.ValidateInputs(*input); !ok {
		validationResponse(errors, w)
	}

	data, err := loginUser(input)
	var ts *user.TokenDetails
	if helper.DoPasswordsMatch(data[1].PasswordHash, input.Password, []byte(data[1].SaltSize)) {
		ts, err = createAccessAndRefreshToken(data[1].CustomerId)

		if err != nil {
			ErrorResponse(http.StatusUnprocessableEntity, err.Error(), w)
			return
		}
	}

	response := make(map[string]interface{})
	response["accessToken"] = ts.AccessToken
	response["refreshToken"] = ts.RefreshToken
	SuccessRespond(response, w)
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

func generateCustomerId() string {
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

// Define the size of the salt
const saltSize = 16

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

	salt := helper.GenerateRandomSalt(saltSize)
	password := helper.HashPassword(generateCustomerId(), salt)

	err := db.QueryRow(sqlStatement, customerId, user.CustomerName, user.PhoneNumber, user.Email,
		user.DateOfBird, user.Sex, salt, password, user.CreatedAt).Scan(&customerId)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", customerId)

	// return the inserted id
	return customerId
}

type ResponseLoginDto struct {
	PasswordHash string
	SaltSize     string
	CustomerId   string
}

func loginUser(login *user.LoginDto) (map[int]ResponseLoginDto, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	result1 := strings.Index(login.Username, "@")

	//If the result below 0, then it mean the username is phone number
	sqlStatement := ""
	if result1 < 0 {
		// create the select sql query
		sqlStatement = `SELECT customer_id,password,salt FROM tbl_customer WHERE email=? and password=?`
	} else {
		sqlStatement = `SELECT customer_id,password,salt FROM tbl_customer WHERE phone_number=? and password=?`
	}

	// execute the sql statement
	stmt, err := db.Prepare(sqlStatement)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var password string
	var salt string
	var customerId string
	err = stmt.QueryRow(1, 2).Scan(&password, &salt, &customerId)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Password: " + password + " SaltSize: " + salt)

	a1 := ResponseLoginDto{PasswordHash: password, SaltSize: salt, CustomerId: customerId}
	data := map[int]ResponseLoginDto{
		1: a1,
	}

	// return empty u on error
	return data, err
}

func SuccessRespond(fields map[string]interface{}, writer http.ResponseWriter) {
	fields["status"] = "success"
	message, err := json.Marshal(fields)
	if err != nil {
		//An error occurred processing the json
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("An error occurred internally"))
	}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(message)
}

func ErrorResponse(statusCode int, error string, writer http.ResponseWriter) {
	//Create a new map and fill it
	fields := make(map[string]interface{})
	fields["status"] = "error"
	fields["message"] = error
	message, err := json.Marshal(fields)

	if err != nil {
		//An error occurred processing the json
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("An error occured internally"))
	}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
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

func createAccessAndRefreshToken(userId string) (*user.TokenDetails, error) {
	td := &user.TokenDetails{

	}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userId

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	atClaims["access_uuid"] = td.AccessUuid
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func extractTokenMetadata(r *http.Request) (*user.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &user.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func tokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r)

		if err != nil {
			ErrorResponse(http.StatusUnauthorized, err.Error(), w)
			next.ServeHTTP(w, r)
			return
		}
	})
}

func deleteAuth(givenUuid string) (int64,error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}