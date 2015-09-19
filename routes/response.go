package routes

type BasicResponse struct {
	Success bool   `json:"requestSuccess"`
	Message string `json:"message"`
}
