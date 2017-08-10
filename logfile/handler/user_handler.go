package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"fmt"
	"go.web.red/module"
)

var User_default *User_entity

func init() {
	User_default = &User_entity{
		module: new(module.User),
		 Red_com_handle: Red_com_handle{
			Path: "user",
			Handler_works: make(map[string]func(echo.Context)(error)),
		},
	}

	User_default.Handler_works[echo.GET] = User_default.get
	User_default.Handler_works[echo.POST] = User_default.post
	User_default.Handler_works[echo.DELETE] = User_default.delete
	User_default.Handler_works[echo.PUT] = User_default.put
}

type User_entity struct {
	Red_com_handle
	module *module.User
}

type (
	User_param struct {
		Name     string			`json:"name" validate:"required,lte_len=50"`
		Age      int                    `json:"age"`
		Address  string                 `json:"address"`
		Status   int			`json:"status"`
		Created  string			`json:"create_time"`
		Modified string			`json:"update_time"`
	}

	User_param_all struct {
		User_param
		Id		int		`json:"id" validate:"required,lte_len=10"`
	}

	User_get_req struct {
		Base_req_entity
		Id		int
		Status		int
	}

	User_get_res struct {
		Base_res_entity
		Users []User_param_all `json:"users"`
	}

	User_post_req struct {
		Base_req_entity
		User_param
	}

	User_put_req struct {
		Base_req_entity
		User_param_all
	}

	User_delete_req struct {
		Base_req_entity
		Id	int	`validate:"required"`
	}
)

func (entity *User_entity) get(cont echo.Context) error {
	req := User_get_req{
		Base_req_entity: Base_req_entity{
			TrackId: cont.QueryParam("trackId"),
		},
		Id: cont.QueryParam("id"),
		Status: cont.QueryParam("status"),
	}

	res := User_get_res{}

	if err := Validate.Struct(&req); err != nil {
		res.Status_code = REQUESTERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	queryParams := map[string]interface{}{}

	queryParams["status"] = req.Status

	users, err := entity.module.Query(req.TrackId, queryParams)

	if err != nil {
		res.Status_code = DBERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	if len(users) == 0 {
		res.Users = []User_param_all{}
	} else {
		for _, user := range users {
			res.Users = append(res.Users, User_param_all{
				User_param: User_param{
					Name: user.Name,
					Age: user.Age,
					Address: user.Address,
					Status: user.Status,
					Modified: user.Modified.Format(DATAFORMAT),
					Created: user.Created.Format(DATAFORMAT),
				},
				Id: user.Id,
			})
		}
	}

	res.Status_code = SUCCESS

	return cont.JSON(http.StatusOK, res)
}

func (entity *User_entity) post(cont echo.Context) error {
	req := new(User_post_req)
	res := new(Base_res_entity)

	if err := cont.Bind(req); err != nil {
		fmt.Println(err.Error())
		res.Status_code = REQUESTERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	req.TrackId = cont.QueryParam("trackId")

	if err := Validate.Struct(req); err != nil {
		fmt.Println(err.Error())
		res.Status_code = REQUESTERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	user := module.User{
		Name: req.Name,
		Age: req.Age,
		Address: req.Address,
		Status: req.Status,
	}

	err := entity.module.Save(req.TrackId, user)

	if err != nil {
		fmt.Println(err.Error())
		res.Status_code = DBERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	res.Status_code = SUCCESS
	return cont.JSON(http.StatusOK, res)
}

func (entity *User_entity) delete(cont echo.Context) error {
	req := User_delete_req{
		Base_req_entity: Base_req_entity{
			TrackId: cont.QueryParam("trackId"),
		},
		Id: cont.QueryParam("id"),
	}

	res := Base_res_entity{}


	if err := Validate.Struct(&req); err != nil {
		res.Status_code = REQUESTERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	err := entity.module.Delete(req.TrackId, req.Id)

	if err != nil {
		res.Status_code = DBERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	res.Status_code = SUCCESS
	return cont.JSON(http.StatusOK, res)
}

func (entity *User_entity) put(cont echo.Context) error {
	req := new(User_put_req)
	res := new(Base_res_entity)

	if err := cont.Bind(req); err != nil {
		fmt.Println(err.Error())
		res.Status_code = REQUESTERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	req.TrackId = cont.QueryParam("trackId")

	if err := Validate.Struct(req); err != nil {
		fmt.Println(err.Error())
		res.Status_code = REQUESTERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	user := module.User{
		Id: req.Id,
		Name: req.Name,
		Age: req.Age,
		Address: req.Address,
		Status: req.Status,
	}

	err := entity.module.Update(req.TrackId, user)

	if err != nil {
		fmt.Println(err.Error())
		res.Status_code = DBERROR
		return cont.JSON(http.StatusBadRequest, res)
	}

	res.Status_code = SUCCESS
	return cont.JSON(http.StatusOK, res)
}
