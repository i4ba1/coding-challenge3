package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/i4ba1/CustomerOrderAPI/router"
	"github.com/i4ba1/CustomerOrderAPI/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSuccessRegister(t *testing.T){
	newRouter := router.SetupRouter()
	t.Run("Register new Customer", func(t *testing.T) {
		customer, _ := json.Marshal(user.UserDto{
			CustomerName: "Muhamad Nizar Iqbal",
			PhoneNumber: "087882458839",
			Email: "muhamad.iqbal1981@gmail.com",
			DateOfBird: time.Now(),
			CreatedAt: time.Now(),
			Sex: true,
		})

		req, err := http.NewRequest("POST", "/api/register", bytes.NewReader(customer))
		if err != nil {
			fmt.Println(err.Error())
		}

		assert.Equal(t, nil, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		newRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRegisterConflicted(t *testing.T){
	newRouter := router.SetupRouter()
	t.Run("Register new Customer but conflicted", func(t *testing.T) {
		customer, _ := json.Marshal(user.UserDto{
			CustomerName: "Muhamad Nizar Iqbal",
			PhoneNumber: "087882458829",
			Email: "muhamad.iqbal1983@gmail.com",
			DateOfBird: time.Now(),
			CreatedAt: time.Now(),
			Sex: true,
		})

		req, err := http.NewRequest("POST", "/api/register", bytes.NewReader(customer))
		if err != nil {
			fmt.Println(err.Error())
		}

		assert.Equal(t, nil, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		newRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestSuccessLogin(t *testing.T){
	newRouter := router.SetupRouter()
	t.Run("Successfully Login with Email or Phone Number", func(t *testing.T) {
		customer, _ := json.Marshal(user.LoginDto{
			Username: "muhamad.iqbal1981@gmail.com",
			Password: "hkizilwxkx",
		})

		req, err := http.NewRequest("POST", "/api/login", bytes.NewReader(customer))
		if err != nil {
			fmt.Println(err.Error())
		}

		assert.Equal(t, nil, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		newRouter.ServeHTTP(w, req)
		fmt.Println(
			"The result on the Body => ",
			w.Body,
		)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}


