package tools

//Response representation response only with message
type Response struct {
	Message string `json:"message"`
}

//ResponseValidation representation Validation response
type ResponseValidation struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

//ResponseData representation response message with data
type ResponseData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//ResponseDataTotal representation response message with total data and data it self
type ResponseDataTotal struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
}

//ResponseLogin representation response message with data and token
type ResponseLogin struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   JWTResponse `json:"token"`
}
