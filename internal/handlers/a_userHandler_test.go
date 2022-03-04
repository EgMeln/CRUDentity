package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestUserHandler_SignIn(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "handler",
		"password": "1234",
	})
	require.NoError(t, err)
	request, err := http.NewRequest("POST", "http://localhost:8081/auth/sign-in", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	log.Infof(string(byteBody))
	require.NoError(t, err)
	_, err = fmt.Sscanf(string(byteBody), "access token: %s , refresh token: %s", &accessToken, &refreshToken)
	err = request.Body.Close()
	if err != nil {
		log.Warnf("can't user sign in %v", err)
	}
}
func TestUserHandler_Add(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "test2",
		"password": "1234",
	})
	require.NoError(t, err)
	request, err := http.NewRequest("POST", "http://localhost:8081/auth/sign-up", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"username":"test2","password":"1234"}`, strings.Trim(string(byteBody), "\n"))
	err = request.Body.Close()
	if err != nil {
		log.Warnf("can't user sign up %v", err)
	}
}
func TestUserHandler_Get(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8081/admin/users/test2", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(byteBody), "test2")

	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't user get %v", err)
	}
}
func TestUserHandler_Update(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "test2",
		"password": "123",
	})
	require.NoError(t, err)
	request, err := http.NewRequest("PUT", "http://localhost:8081/admin/users", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(byteBody), "123")
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close user update %v", err)
	}
}
func TestUserHandler_GetAll(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8081/admin/users", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(byteBody), "test")
	require.Contains(t, string(byteBody), "test2")
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't user get all %v", err)
	}
}

func TestUserHandler_Refresh(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"username": "test",
	})
	request, err := http.NewRequest("POST", "http://localhost:8081/refresh", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	var access, refresh string
	_, err = fmt.Sscanf(string(byteBody), "new access token: %s , new refresh token: %s", &access, &refresh)
	if access == "" && refresh == "" {
		log.Warnf("can't refresh tokens %v", err)
	}
	err = request.Body.Close()
	if err != nil {
		log.Warnf("can't refresh tokens %v", err)
	}
}
func TestUserHandler_Delete(t *testing.T) {
	request, err := http.NewRequest("DELETE", "http://localhost:8081/admin/users/test2", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "{}\n", string(byteBody))
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close user delete %v", err)
	}
}
