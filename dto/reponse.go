package dto

type Response struct {
	ErrorCode int    `json:"errorCode"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

type ResponseData struct {
	Response
	Data interface{} `json:"data"`
}
