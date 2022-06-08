package sender

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"testing"
)

func TestPostgres_DoesPhoneNumberExist(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Errorf("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: phone_number 0817289398 should not exist.
	exist, err := repo.DoesPhoneNumberExist("0817289398")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if exist {
		t.Errorf("Sender phone_number %v should not exist. exist: %v", "0817289398", exist)
	}
	//Case 2: phone_number 0924925265 should exist.
	exist, err = repo.DoesPhoneNumberExist("0924925265")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if !exist {
		t.Errorf("Sender phone_number %v should exist. exist: %v", "0924925265", exist)
	}
}

func TestPostgres_DoesIdExist(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Errorf("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: id 99 should not exist.
	exist, err := repo.DoesIdExist(99)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if exist {
		t.Errorf("Sender id %v should not exist. exist: %v", 99, exist)
	}
	//Case 2: id 1 should exist.
	exist, err = repo.DoesIdExist(1)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if !exist {
		t.Errorf("Sender id %v should exist. exist: %v", 1, exist)
	}
}

func TestPostgres_FindById(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Errorf("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: id 99 should not found.
	s, err := repo.FindById(99)
	if err != nil {
		switch err.(type) {
		case internal.ErrNotFound:
			fmt.Printf("OK: %v\n", s)
		case internal.ErrDatabase:
			t.Errorf("err: %v", err)
		default:
			t.Errorf("Sender id %v should not found", 99)
		}
	}
	//Case 2: Sender id 1 should found.
	s, err = repo.FindById(1)
	if err != nil {
		switch err.(type) {
		case internal.ErrNotFound:
			t.Errorf("Sender id %v should found", 1)
		default:
			t.Errorf("err: %v", err)
		}
	}
	fmt.Printf("OK: %v\n", s)

}

func TestPostgres_Create(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Errorf("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: should success.
	s, err := repo.Create(entity.Sender{
		SenderType:       2,
		NationalIdNumber: "6264728304725",
		BirthDate:        "1970-12-01",
		ContactInfo: entity.ContactInfo{
			Name:        "JNA Distribution 3",
			PhoneNumber: "0924925270",
			Address:     "19 ถนน Rama II Soi 30",
			District:    "Bang Mot",
			State:       "Chom Thong",
			Province:    "Bangkok",
			Postcode:    "10150",
		},
	})
	if err != nil {
		t.Errorf("NO: %+v", err)
		return
	}
	if s.Id == 0 {
		t.Errorf("NO: %+v", err)
		return
	}
	fmt.Printf("OK: %+v\n", s)
	//Case 2, duplicate phone_number should fail.
	s, err = repo.Create(entity.Sender{
		SenderType:       2,
		NationalIdNumber: "6264728304725",
		BirthDate:        "1970-12-01",
		ContactInfo: entity.ContactInfo{
			Name:        "JNA Distribution 3",
			PhoneNumber: "0924925270",
			Address:     "19 ถนน Rama II Soi 30",
			District:    "Bang Mot",
			State:       "Chom Thong",
			Province:    "Bangkok",
			Postcode:    "10150",
		},
	})
	if err == nil {
		t.Errorf("NO: %+v", err)
		return
	}
	fmt.Printf("OK: %+v, %+v\n", s, err)
}

func TestPostgres_UpdateById(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Errorf("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: Update Name to "JNA Distribution XX"
	newName := "JNA Distribution 111xxx"
	s, err := repo.UpdateById(entity.Sender{
		ContactInfo: entity.ContactInfo{
			Name: newName,
		},
	}, 1)
	if err != nil {
		t.Errorf("NO: %v, %v", s, err)
		return
	}
	if s.Id != 1 {
		t.Errorf("NO: %v, %v", s, err)
	}
	if s.Name != newName {
		t.Errorf("NO: Name not equal: %v, name: %v\n", "JNA Distribution 111xx", s.Name)
		return
	}
	fmt.Printf("OK: %v, %v\n", s, err)
	//Case 2: Update Postcode to "10151"
	newPostCode := "10151"
	s, err = repo.UpdateById(entity.Sender{
		Id: 1,
		ContactInfo: entity.ContactInfo{
			Postcode: newPostCode,
		},
	}, 1)
	if err != nil {
		t.Errorf("NO: %v, %v", s, err)
		return
	}
	if s.Id != 1 {
		t.Errorf("NO: %v, %v", s, err)
	}
	if s.Postcode != newPostCode {
		t.Errorf("NO: Postcode not equal: %v, postcode: %v\n", newPostCode, s.Postcode)
		return
	}
	fmt.Printf("OK: %v, %v\n", s, err)
}

func TestPostgres_UpdateByPhoneNumber(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Errorf("init db failed.")
		log.Fatal(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1
	newName := "JNA Distribution XYZ"
	s, err := repo.UpdateByPhoneNumber(entity.Sender{
		SenderType:       2,
		NationalIdNumber: "6264728304725",
		BirthDate:        "1970-12-01",
		ContactInfo: entity.ContactInfo{
			Name:        newName,
			PhoneNumber: "0924925268",
			Address:     "19 ถนน Rama II Soi 30",
			District:    "Bang Mot",
			State:       "Chom Thong",
			Province:    "Bangkok",
			Postcode:    "10150",
		},
	}, "0924925268")
	if err != nil {
		t.Errorf("NO: %+v", err)
		return
	}
	if s.Id == 0 {
		t.Errorf("NO: %+v", s)
		return
	}
	if s.Name != newName {
		t.Errorf("NO: %+v", s)
		return
	}
	fmt.Printf("OK: %+v\n", s)
	//Case 2
	newPostCode := "10150"
	s, err = repo.UpdateByPhoneNumber(entity.Sender{
		SenderType:       2,
		NationalIdNumber: "6264728304725",
		BirthDate:        "1970-12-01",
		ContactInfo: entity.ContactInfo{
			Name:        newName,
			PhoneNumber: "0924925268",
			Address:     "19 ถนน Rama II Soi 30",
			District:    "Bang Mot",
			State:       "Chom Thong",
			Province:    "Bangkok",
			Postcode:    newPostCode,
		},
	}, "0924925268")
	if err != nil {
		t.Errorf("NO: %+v", err)
		return
	}
	if s.Id == 0 {
		t.Errorf("NO: %+v", s)
		return
	}
	if s.Postcode != newPostCode {
		t.Errorf("NO: %+v", s)
		return
	}
	fmt.Printf("OK: %+v\n", s)
}
