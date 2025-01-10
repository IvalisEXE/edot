package response

// Response is the response struct
type Response struct {
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}
