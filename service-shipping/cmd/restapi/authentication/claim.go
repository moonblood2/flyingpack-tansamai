package authentication

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"net/http"
	"strconv"
)

const (
	headerUserId    = "user-id"
	headerUserRole  = "user-role"
	headerContactId = "contact-id"
)

type Claims struct {
	UserId    string      `json:"uid"`
	UserRole  entity.Role `json:"rol"`
	ContactId int16       `json:"contact_id,omitempty"`
	jwt.StandardClaims
}

func AttachClaimsWithHeader(claims *Claims, r *http.Request) {
	r.Header.Set(headerUserId, fmt.Sprintf("%v", claims.UserId))
	r.Header.Set(headerUserRole, fmt.Sprintf("%v", claims.UserRole))
	if claims.ContactId > 0 {
		r.Header.Set(headerContactId, fmt.Sprintf("%v", claims.ContactId))
	}
}

func DetachClaimsFromHeader(r *http.Request) (Claims, error) {
	claims := Claims{}
	if r.Header.Get(headerUserId) != "" {
		claims.UserId = r.Header.Get(headerUserId)
	} else {
		return Claims{}, fmt.Errorf("Not have %v in Header. ", headerUserRole)
	}
	if r.Header.Get(headerUserRole) != "" {
		role, err := strconv.ParseInt(r.Header.Get(headerUserRole), 10, 16)
		if err != nil {
			return Claims{}, err
		}
		claims.UserRole = entity.Role(role)
	} else {
		return Claims{}, fmt.Errorf("Not have %v in Header. ", headerUserRole)
	}
	if r.Header.Get(headerContactId) != "" {
		id, err := strconv.ParseInt(r.Header.Get(headerContactId), 10, 16)
		if err != nil {
			return Claims{}, err
		}
		claims.ContactId = int16(id)
	}
	return claims, nil
}
