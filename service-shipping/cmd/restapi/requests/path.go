package requests

import (
	"github.com/jna-distribution/service-shipping/internal"
	"strings"
)

// PathParam map key is parameter.
// eg. users/:id/
// PathParam = {"id": 1}
type PathParam map[string]string

// NewPathParam make new map[string]string for each placeholder ":" parameter.
func NewPathParam(pattern, requestURI string) (PathParam, error) {
	uriArr := strings.Split(requestURI, "/")
	patArr := strings.Split(pattern, "/")
	uriArrLen := len(uriArr)
	patArrLen := len(patArr)

	if uriArrLen != patArrLen {
		return nil, internal.ErrNotFound{}
	}

	params := PathParam{}
	for i, v := range patArr {
		if strings.HasPrefix(v, ":") {
			params[strings.TrimPrefix(v, ":")] = uriArr[i]
		} else if strings.Compare(uriArr[i], v) != 0 {
			return nil, internal.ErrNotFound{}
		}
	}

	return params, nil
}

func MatchPattern(pattern, requestURI string) bool {
	patArr := strings.Split(pattern, "/")
	uriArr := strings.Split(requestURI, "/")
	uriArrLen := len(uriArr)
	patArrLen := len(patArr)

	if uriArrLen != patArrLen {
		return false
	}
	for i := 0; i < uriArrLen; i++ {
		//Skip placeholder
		if strings.HasPrefix(patArr[i], ":") {
			continue
		}
		if patArr[i] != uriArr[i] {
			return false
		}
	}
	return true
}
