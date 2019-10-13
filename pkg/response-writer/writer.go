package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SuccessWriter(w http.ResponseWriter, data interface{}) error {

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataJSON)

	return nil
}

func ErrorWriter(w http.ResponseWriter, status int, errors interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status-Code", fmt.Sprintf("%d", status))
	w.WriteHeader(status)

	errorsResponse := map[string]interface{}{
		"status": http.StatusText(status),
		"error":  errors,
	}
	errorsBytes, err := json.Marshal(errorsResponse)
	if err != nil {
		log.Println(err)
		return err
	}

	w.Write(errorsBytes)

	return nil
}
