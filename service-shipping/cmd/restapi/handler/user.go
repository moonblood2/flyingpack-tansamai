package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jna-distribution/service-shipping/cmd/restapi/middleware"
	"github.com/jna-distribution/service-shipping/cmd/restapi/requests"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/user"
	"net/http"
)

// userResponse struct for responses in /users/ endpoint.
type userResponse struct {
	Id         string          `json:"id"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	Role       entity.Role     `json:"role"`
	RoleString string          `json:"role_string"`
	Contact    *entity.Contact `json:"contact,omitempty"`
}

// BuildUserHandler receive RootMux and register path with handler.
func BuildUserHandler(rootMux *http.ServeMux, svc user.Service) {
	rootMux.Handle("/users/", usersHandler(
		middleware.WithPermission(getUser(svc), entity.RoleAdmin, entity.RoleAccounting),
		middleware.WithPermission(register(svc), entity.RoleAdmin),
		middleware.WithPermission(editUser(svc), entity.RoleAdmin),
		middleware.WithPermission(removeUser(svc), entity.RoleAdmin),
	))
}

// usersHandler map each handler of /users/ endpoint for each method.
func usersHandler(get, post, patch, delete http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get.ServeHTTP(w, r)
		case http.MethodPost:
			post.ServeHTTP(w, r)
		case http.MethodPatch:
			patch.ServeHTTP(w, r)
		case http.MethodDelete:
			delete.ServeHTTP(w, r)
		default:
			rw := responses.NewJSONResponseWriter(w)
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
		}
	})
}

// getUser handle request that match for "/users/" path and GET method, query all users or specific user id.
func getUser(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO 1, Pagination for GetAll().
		//Not found id, find all users.
		rw := responses.NewJSONResponseWriter(w)
		if r.RequestURI == "/users/" {
			users, err := svc.GetAll()
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			//Make new user responses, omit password,
			//created_at and deleted_at from UserEntity and translate role number to string.
			var ur []userResponse
			for _, u := range users {
				ur = append(ur, userResponse{
					Id:         u.Id,
					Name:       u.Name,
					Email:      u.Email,
					Role:       u.Role,
					RoleString: u.GetRoleString(),
				})
			}
			rw.WriteWithStatus(http.StatusOK, ur)
		} else {
			//Find by specific user id.
			params, err := requests.NewPathParam("/users/:id/", r.RequestURI)
			if _, ok := params["id"]; !ok && err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			id := params["id"]
			//Call user service for get user by id.
			u, err := svc.GetById(id)
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			//Create new user responses.
			ur := userResponse{
				Id:         u.Id,
				Email:      u.Email,
				Name:       u.Name,
				Role:       u.Role,
				RoleString: u.GetRoleString(),
			}
			rw.WriteWithStatus(http.StatusOK, ur)
		}
	})
}

// register create new user if e-mail not used.
func register(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		//Retrieve input
		input := user.RegisterInput{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}

		//Register user.
		output, err := svc.Register(input)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		//Create new user responses.
		res := userResponse{
			Id:         output.User.Id,
			Name:       output.User.Name,
			Email:      output.User.Email,
			Role:       output.User.Role,
			RoleString: output.User.GetRoleString(),
		}
		if output.User.Role == entity.RoleShop {
			res.Contact = &output.Contact
		}
		rw.WriteWithStatus(http.StatusCreated, res)
	})
}

// editUser handle request that match for "/users/" path and PATCH method.
func editUser(s user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		params, err := requests.NewPathParam("/users/:id/", r.RequestURI)
		if _, ok := params["id"]; err != nil && !ok {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		//Receive input from JSON.
		input := user.EditUserInput{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		id := params["id"]
		//Add id uint16 to input.
		input.Id = id
		//Call edit service.
		u, err := s.Edit(input)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		//Create new user response.
		ur := userResponse{
			Id:         u.Id,
			Email:      u.Email,
			Name:       u.Name,
			Role:       u.Role,
			RoleString: u.GetRoleString(),
		}
		rw.WriteWithStatus(http.StatusOK, ur)
	})
}

// removeUser handle request that match for "/users/" path and DELETE method.
func removeUser(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		params, err := requests.NewPathParam("/users/:id/", r.RequestURI)
		if _, ok := params["id"]; err != nil && !ok {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		id := params["id"]
		//Call service to remove user.
		if err := svc.Remove(id); err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, responses.Message{
			Message: fmt.Sprintf("User %v removed.", id),
		})
	})
}
