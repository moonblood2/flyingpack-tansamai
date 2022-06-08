package ordparcel

import (
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/pkg"
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
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)
	count := 1
	modp := mock.OrderParcel(count, "")
	odp, err := repo.Insert(modp)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		if odp[i].OrderParcelId == "" {
			t.Errorf("odp[%v].OrderParcelId=%v, want: other", i, odp[i].OrderParcelId)
		}

		if odp[i].ProviderCode != modp[i].ProviderCode {
			t.Errorf("odp[%v].ProviderCode=%v, want: %v", i, odp[i].ProviderCode, modp[i].ProviderCode)
		}
		if odp[i].Price != modp[i].Price {
			t.Errorf("odp[%v].Price=%v, want: %v", i, odp[i].Price, modp[i].Price)
		}
		if odp[i].PaymentMethod != modp[i].PaymentMethod {
			t.Errorf("odp[%v].PaymentMethod=%v, want: %v", i, odp[i].PaymentMethod, modp[i].PaymentMethod)
		}
		if odp[i].ParcelShape != modp[i].ParcelShape {
			t.Errorf("odp[%v].ParcelShape=%v, want: %v", i, odp[i].ParcelShape, modp[i].ParcelShape)
		}

		t.Logf("modp[%v] = %v\n", i, modp[i])
		t.Logf("odp[%v] = %v\n", i, odp[i])
	}
}

func TestRepository_Find(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)
	r, err := repo.Find("", "2021-02-20", "2021-02-20")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", r)
}

func TestRepository_FindOrderByIds(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)
	r, err := repo.FindOrderByIds("", []string{"f98905e6-0ced-4434-8480-00a749fe539d", "8c0c7458-c3f7-46a4-bf97-20c43a5fdf33"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", r)
}

func TestPostgres_UpdateOrder(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)
	repo.UpdateOrder(entity.OrderParcel{OrderParcelId: "31a1bcfe-8ba6-4ff3-b439-cba8833e6186", TrackingCode: "SHIPOPXXX1026"})
}
