package origin

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"testing"
)

func TestPostgres_DoesPhoneNumberExist(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	repo := NewPostgresRepository(db)
	//Case 1: Find 0902504995, should exists.
	exist, err := repo.DoesPhoneNumberExist("0902504995")
	if err != nil {
		t.Errorf("NO: %v, %v", exist, err)
		return
	}
	if !exist {
		t.Errorf("NO: %v, %v", exist, err)
		return
	}
	fmt.Printf("OK: %v, %v\n", exist, err)
	//Case 2: Find 191, should not exists.
	exist, err = repo.DoesPhoneNumberExist("191")
	if err != nil {
		t.Errorf("NO: %v, %v", exist, err)
		return
	}
	if exist {
		t.Errorf("NO: %v, %v", exist, err)
		return
	}
	fmt.Printf("OK: %v, %v\n", exist, err)
}

func TestPostgres_DoesPhoneNumbersExist(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	repo := NewPostgresRepository(db)

	//Case 1
	existList, err := repo.DoesPhoneNumbersExist([]string{"0902504995", "022700499", "022700499", "0862828903", "0902504995", "0976493333"})
	correctAns := []bool{true, true, true, false, true, false}
	if err != nil {
		t.Errorf("NO: %v", err)
		return
	}
	for i, v := range correctAns {
		if existList[i] != v {
			t.Errorf("NO: %v %v", i, existList)
			return
		}
	}
	fmt.Printf("OK: %v\n", existList)
}

func TestPostgres_UpsertByPhoneNumber(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)

	no := mock.Origin(10)
	for i, v := range no {
		t.Logf("newOrigins[%v]: %v\n", i, v)
	}
	o, err := repo.UpsertByPhoneNumber(no)
	if err != nil {
		t.Fatal(err)
	}
	for i := range o {
		t.Logf("o[%v]: %+v\n", i, o[i])
	}
}
