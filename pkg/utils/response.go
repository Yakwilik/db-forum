package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type InterfaceMap map[string]interface{}

func JSONResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	logrus.Println(message)
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(message)
	logrus.Printf("error: %+v", err)
}
