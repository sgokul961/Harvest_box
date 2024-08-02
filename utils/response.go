package utils

import (
	"HarvestBox/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HandleError logs the error and sends an error response to the client.
func HandleError(w http.ResponseWriter, err error, description string, statusCode int) {
	log.Println(description+":", err)
	jsonResponse(w, statusCode, "Failed", description, err.Error())
}

// SuccessResponse sends a success response to the client.
func SuccessResponse(w http.ResponseWriter, description string, responseData interface{}, statusCode int) {
	jsonResponse(w, statusCode, "Success", description, responseData)
}

// jsonResponse sends a JSON response to the client.
func jsonResponse(w http.ResponseWriter, status int, responseStatus, responseDescription string, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := models.GenaralResponse{
		ResponseStatus:      responseStatus,
		ResponseDescription: responseDescription,
		ResponseData:        responseData,
		ResponseCode:        status,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("error encoding JSON response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func FormatDateFromDate(inputDate time.Time) string {
	// Format the input time as per the desired date format
	formattedDate := inputDate.Format("2006-01-02")
	return formattedDate
}
