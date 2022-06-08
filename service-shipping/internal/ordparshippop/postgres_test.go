package ordparshippop

import (
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
	modp := mock.OrderParcelShippop(count, []string{"f98905e6-0ced-4434-8480-00a749fe539d"})
	odp, err := repo.Insert(modp)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		if odp[i].Id == 0 {
			t.Errorf("odp[%v].Id=%v, want: other", i, odp[i].Id)
		}
		if odp[i].OrderParcelId != modp[i].OrderParcelId {
			t.Errorf("odp[%v].OrderParcelId=%v, want: %v", i, odp[i].OrderParcelId, modp[i].OrderParcelId)
		}
		if odp[i].PurchaseId != modp[i].PurchaseId {
			t.Errorf("odp[%v].PurchaseId=%v, want: %v", i, odp[i].PurchaseId, modp[i].PurchaseId)
		}
		if odp[i].Status != modp[i].Status {
			t.Errorf("odp[%v].Status=%v, want: %v", i, odp[i].Status, modp[i].Status)
		}
		if odp[i].CourierCode != modp[i].CourierCode {
			t.Errorf("odp[%v].CourierCode=%v, want: %v", i, odp[i].CourierCode, modp[i].CourierCode)
		}
		if odp[i].CourierTrackingCode != modp[i].CourierTrackingCode {
			t.Errorf("odp[%v].CourierTrackingCode=%v, want: %v", i, odp[i].CourierTrackingCode, modp[i].CourierTrackingCode)
		}
		if odp[i].TrackingCode != modp[i].TrackingCode {
			t.Errorf("odp[%v].TrackingCode=%v, want: %v", i, odp[i].TrackingCode, modp[i].TrackingCode)
		}
		if odp[i].CODAmount != modp[i].CODAmount {
			t.Errorf("odp[%v].CODAmount=%v, want: %v", i, odp[i].CODAmount, modp[i].CODAmount)
		}

		t.Logf("modp[%v] = %v\n", i, modp[i])
		t.Logf("odp[%v] = %v\n", i, odp[i])
	}
}
