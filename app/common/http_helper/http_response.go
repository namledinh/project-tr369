package httphelper

import "encoding/json"

type HttpResponse struct {
	Message   string      `json:"message" `         // message of the response
	StatusKey string      `json:"status_key"`       // status key to identify the error type or success
	Data      interface{} `json:"data,omitempty"`   // data of the response
	Paging    interface{} `json:"paging,omitempty"` // paging of the response
	Filter    interface{} `json:"filter,omitempty"` // filter of the response
}

func (r *HttpResponse) ToJsonBytes() []byte {
	jsonBytes, _ := json.Marshal(r)
	return jsonBytes
}

type emptyStruct struct{}

func NewSuccessResponse(data, paging, filter interface{}) *HttpResponse {
	return &HttpResponse{
		Message:   "request successfully",
		StatusKey: "Success",
		Data:      data,
		Paging:    paging,
		Filter:    filter,
	}
}

func NewErrorHTTPResponse(data interface{}, message, statusKey string) *HttpResponse {
	if data == nil {
		data = emptyStruct{}
	}

	return &HttpResponse{
		Message:   message,
		StatusKey: statusKey,
		Data:      data,
		Paging:    nil,
		Filter:    nil,
	}
}
