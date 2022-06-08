package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jna-distribution/service-shipping/cmd/restapi/authentication"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg"
	"net/http"
	"strings"
)

// WithPermission middleware for check each handlers.
func WithPermission(next http.Handler, permittedRoles ...entity.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Create new JSON response writer.
		rw := responses.NewJSONResponseWriter(w)
		env, err := pkg.LoadEnv()
		if err != nil {
			e := responses.BuildErrorResponse(responses.ErrInternalServer{InternalError: err})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		signingKey := []byte(env["JWT_SIGNING_KEY"])
		//Trim out "Bearer " from HEADER["Authorization"].
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		//Parse token to verify.
		token, err := jwt.ParseWithClaims(tokenString, &authentication.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		//Check for error if err != nil, in line 84 of "jwt-go@v3.2.0+incompatible\parser.go",
		//tells if ParseWithClaims() return nil that means token.Valid is true (token.Valid = true),
		//thus it doesn't matter check token.Valid.
		if err != nil {
			e := responses.BuildErrorResponse(responses.ErrUnauthorized{InternalError: err, Details: err.Error()})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		//Retrieve claims, convert token.Claims to concrete UserClaims struct.
		claims, ok := token.Claims.(*authentication.Claims)
		//Check results of retrieve claims.
		if !ok {
			e := responses.BuildErrorResponse(responses.ErrUnauthorized{InternalError: err})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		//Check if not permit.
		//eg. permittedRole = 0011
		//    UserRole = 1000
		//	  permittedRole & userRole = 0
		if len(permittedRoles) > 0 {
			permitted := permittedRoles[0]
			for _, v := range permittedRoles {
				permitted = permitted | v
			}
			if permitted&claims.UserRole == 0 {
				e := responses.BuildErrorResponse(responses.ErrForbidden{})
				rw.WriteWithStatus(e.Status, e)
				return
			}
		}
		//Pass all conditions, let's go to the next handler.
		//Save userId and userRole in header, reduce time in decryption token.
		r.Header.Set("user-id", fmt.Sprintf("%v", claims.UserId))
		r.Header.Set("user-role", fmt.Sprintf("%v", claims.UserRole))
		r.Header.Set("contact-id", fmt.Sprintf("%v", claims.ContactId))
		next.ServeHTTP(w, r)
	})
}
