package user

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/contact"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"testing"
)

func TestService_RegisterInput_Validation(t *testing.T) {
	i := RegisterInput{
		Name:            "Kozuki Oden",
		Email:           "oden@gmail.com",
		Role:            entity.RoleAdmin,
		Password:        "1234567891",
		ConfirmPassword: "1234567891",
	}
	if i.Validate() != nil {
		t.Errorf("i.Validate()=%v, nil\n", i.Validate())
	}
	i.Role = entity.RoleShop
	if i.Validate() == nil {
		t.Errorf("i.Validate()=%v, not nil\n", i.Validate())
	}
}

func TestService_Register(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	userRepo := NewPostgresRepository(db)
	contactRepo := contact.NewPostgresRepository(db)
	userSvc := NewService(userRepo, contactRepo)

	//Case 1: Create ADMIN, should fail, duplicate email.
	o, err := userSvc.Register(RegisterInput{
		Name:            "123 Contact",
		Email:           "admin01@gmail.com",
		Role:            1,
		Password:        "111",
		ConfirmPassword: "111",
	})
	if err == nil {
		t.Errorf("NO: %v, %v", o, err)
		return
	}
	fmt.Printf("OK: %v, %v\n", o, err)
}
