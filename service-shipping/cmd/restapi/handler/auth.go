package handler

import (
	"encoding/json"
	"github.com/jna-distribution/service-shipping/cmd/restapi/authentication"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/user"
	"net/http"
)

// BuildAuthHandler build handler for /auth/ path.
func BuildAuthHandler(rootMux *http.ServeMux, svc user.Service) {
	rootMux.Handle("/auth/login/", login(svc))
}

// loginResponse response for login().
type loginResponse struct {
	Token      string          `json:"token"`
	Id         string          `json:"id"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	Role       entity.Role     `json:"role"`
	RoleString string          `json:"role_string"`
	Contact    *entity.Contact `json:"contact,omitempty"`
}

// login handler for '/auth/login/'.
func login(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Create new JSON Response Writer.
		rw := responses.NewJSONResponseWriter(w)
		//Check method.
		if r.Method == "POST" {
			//Retrieve data from JSON body.
			input := user.LoginInput{}
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			//Call service for login.
			o, err := svc.Login(input)
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			//Password matched, create new token for user.
			ss, err := authentication.CreateToken(o.User.Id, o.User.Role, o.Contact.Id)
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			//Write loginResponse.
			res := loginResponse{}
			res.Token = ss
			res.Id = o.User.Id
			res.Email = o.User.Email
			res.Name = o.User.Name
			res.Role = o.User.Role
			res.RoleString = o.User.GetRoleString()
			if o.User.IsShop() {
				res.Contact = &o.Contact
			} else {
				res.Contact = nil
			}
			rw.WriteWithStatus(http.StatusOK, res)
			return
		}

		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
	})
}
