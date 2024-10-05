package response

type CustomResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type IDResponse struct {
	ID uint `json:"id"`
}
