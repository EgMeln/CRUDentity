package handlers

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestParkingLotHandler_Add(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"num":        1111,
		"in_parking": false,
		"remark":     "vvvv",
	})
	require.NoError(t, err)
	request, err := http.NewRequest("POST", "http://localhost:8081/admin/park", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"num":1111,"in_parking":false,"remark":"vvvv"}`, strings.Trim(string(byteBody), "\n"))
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close parking lot create %v", err)
	}
}
func TestParkingLotHandler_GetAll(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8081/user/park", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `[{"num":1111,"in_parking":false,"remark":"vvvv"}]`, strings.Trim(string(byteBody), "\n"))
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close parking lot get all %v", err)
	}
}
func TestParkingLotHandler_GetByNum(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8081/user/park/1111", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"num":1111,"in_parking":false,"remark":"vvvv"}`, strings.Trim(string(byteBody), "\n"))
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close parking lot get all %v", err)
	}
}
func TestParkingLotHandler_Update(t *testing.T) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"num":        1111,
		"in_parking": true,
		"remark":     "wwww",
	})
	require.NoError(t, err)
	request, err := http.NewRequest("PUT", "http://localhost:8081/admin/park", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"num":1111,"in_parking":true,"remark":"wwww"}`, strings.Trim(string(byteBody), "\n"))
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close parking lot update %v", err)
	}
}
func TestParkingLotHandler_Delete(t *testing.T) {
	request, err := http.NewRequest("DELETE", "http://localhost:8081/admin/park/1111", nil)
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
		log.Warnf("can't body close parking lot delete %v", err)
	}
}
