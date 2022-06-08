package destination

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

func TestPostgres_UpsertByPhoneNumber(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)

	nd := mock.Destination(10)
	for i, v := range nd {
		t.Logf("newDestinations[%v]: %v\n", i, v)
	}
	d, err := repo.UpsertByPhoneNumber(nd)
	if err != nil {
		t.Fatal(err)
	}
	for i := range d {
		t.Logf("d[%v]: %+v\n", i, d[i])
	}
}
