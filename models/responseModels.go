package models

type GenaralResponse struct {
	ResponseStatus      string      `json:"responseStatus"`
	ResponseData        interface{} `json:"responseData"`
	ResponseCode        int         `json:"responseCode"`
	ResponseDescription string      `json:"responseDescription"`
}

// type MasterTableResponse struct {
// 	Id   string `json:"id"`
// 	Name string `json:"name"`
// 	Date string `json:"date"`
// 	Time string `json:"time"`
// }

// type TaxRate struct {
// 	Id         string `json:"id"`
// 	Name       string `json:"name"`
// 	Percentage string `json:"percentage"`
// }
