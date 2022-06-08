package requests

import "testing"

func TestNewPathParam(t *testing.T) {
	params, err := NewPathParam("/users/:id/", "/users/1/")
	if err != nil {
		t.Errorf("errs: %v\n", err)
	}

	id, ok := params["id"]
	if !ok {
		t.Errorf("can't retrieve id: %v\t, %v", ok, params)
	}
	if id != "1" {
		t.Errorf("id not equal 1\n")
	}
}
