package product

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"testing"
)

func TestPostgres_Insert(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)

	//Case 1
	newProduct := entity.Product{
		ContactId: 3,
		Name:      "Package 1",
		Price:     99,
	}
	p, err := repo.Insert(newProduct)
	if err != nil {
		t.Errorf("NO: %v, %v", p, err)
		return
	}
	//Check Id
	if p.Id == 0 {
		t.Errorf("NO: %v, %v", p, err)
		return
	}
	//Check ContactId
	if p.ContactId != newProduct.ContactId {
		t.Errorf("NO: %v, %v", p, err)
		return
	}
	//Check Name
	if p.Name != newProduct.Name {
		t.Errorf("NO: %v, %v", p, err)
		return
	}
	//Check Price
	if p.Price != newProduct.Price {
		t.Errorf("NO: %v, %v", p, err)
		return
	}
	fmt.Printf("OK: %v, %v\n", p, err)
}

func TestPostgres_DoesIdsExistByContactId(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)

	//Case 1
	ids := []int32{2, 1, 2}
	existList, err := repo.DoesIdsExistByContactId(ids, 3)

	if err != nil {
		t.Errorf("NO: %v", err)
		return
	}
	if existList[0] != true {
		t.Errorf("NO: %v", existList)
		return
	}
	if existList[1] == true {
		t.Errorf("NO: %v", existList)
		return
	}
	if existList[2] != true {
		t.Errorf("NO: %v", existList)
		return
	}
	fmt.Printf("OK: %v\n", existList)
}
