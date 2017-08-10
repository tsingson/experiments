package handler

import "github.com/labstack/echo"

// http response status_code
const(
	SUCCESS 		= "0"
	DBERROR 		= "1"
	SERVERERROR 	= "2"
	REQUESTERROR 	= "3"
)

const DATAFORMAT = "2006-01-02 15:04:05"

type (
	Base_req_entity struct {
		TrackId		string	`validate:"lte_len=20"`
	}

	Base_res_entity struct {
		Status_code		string `json:"status_code"`
	}

)
type Red_handle interface {
	Access_path() string
	Support_method() []string
	Handler(cont echo.Context) error
}

type Red_com_handle struct {
	Path          string
	Handler_works map[string]func(echo.Context)(error)
}

func (this *Red_com_handle)Access_path() string {
	return this.Path
}

func (this *Red_com_handle)Support_method() []string {
	sup_methods := []string{}
	for method, _ := range this.Handler_works {
		sup_methods = append(sup_methods, method)
	}
	return sup_methods
}

func (this *Red_com_handle)Handler(cont echo.Context) (ret error) {
	return 	this.Handler_works[cont.Request().Method](cont)
}
