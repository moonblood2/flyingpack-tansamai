package ordproduct

import (
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestPostgres_Insert(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	count := 3
	modp := mock.OrderProduct(count)
	odp, err := repo.Insert(modp)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		if odp[i].Id == 0 {
			t.Errorf("odp[%v].Id=%v, want: other", i, odp[i].Id)
		}
		if odp[i].UserId != modp[i].UserId {
			t.Errorf("odp[%v].ContactId=%v, want: %v", i, odp[i].UserId, modp[i].UserId)
		}
		if odp[i].SenderId != modp[i].SenderId {
			t.Errorf("odp[%v].SenderId=%v, want: %v", i, odp[i].SenderId, modp[i].SenderId)
		}
		if odp[i].ProductId != modp[i].ProductId {
			t.Errorf("odp[%v].ProductId=%v, want: %v", i, odp[i].ProductId, modp[i].ProductId)
		}
		if odp[i].Quantity != modp[i].Quantity {
			t.Errorf("odp[%v].Quantity=%v, want: %v", i, odp[i].Quantity, modp[i].Quantity)
		}
		if odp[i].PaymentMethod != modp[i].PaymentMethod {
			t.Errorf("odp[%v].PaymentMethod=%v, want: %v", i, odp[i].PaymentMethod, modp[i].PaymentMethod)
		}
		t.Logf("odp[%v] = %v\n", i, odp[i])
		t.Logf("modp[%v] = %v\n", i, modp[i])
	}
}

func TestRepository_Find(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)
	r, err := repo.Find("", "2020-11-01", "2021-01-24")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", r)
}
