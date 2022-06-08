package contact

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"testing"
)

func TestPostgres_Create(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Error("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)

	//Case 1: user id = 3, should return err, about unique key.
	s, err := repo.Create(entity.Contact{
		UserId: "3",
	})
	if err == nil {
		t.Errorf("NO: %v, %v", s, err)
		return
	}
	fmt.Println(s, err)
	//Case 2: user id = 30000, should return err, about foreign key.
	s, err = repo.Create(entity.Contact{
		UserId: "30000",
	})
	if err == nil {
		t.Errorf("NO: %v %v", s, err)
		return
	}
	fmt.Println(s, err)
	//Case 3: user id = 4, should success.
	s, err = repo.Create(entity.Contact{
		UserId: "7",
	})
	if err != nil {
		t.Errorf("err: %v", err)
		return
	}
	fmt.Println(s, err)
}

func TestPostgres_FindByUserId(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Error("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: userId 2, should found.
	s, err := repo.FindByUserId("2")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if s.UserId != "" {
		t.Errorf("user id != 2\n")
	}
	fmt.Printf("%#v\n", s)

	//Case 2: userId 1, should not found.
	s, err = repo.FindByUserId("")
	if err != nil {
		switch err.(type) {
		case internal.ErrNotFound:
			fmt.Printf("%#v\n", s)
		default:
			t.Errorf("err: %v", err)
		}
	}
}
