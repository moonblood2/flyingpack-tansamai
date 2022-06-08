package ordparflash

import (
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/pkg"
	"testing"
)

func TestPostgres_Insert(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	count := 1
	mock.Seed(-1)
	modp := mock.OrderParcelFlash(count, []string{"", ""})
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
		if odp[i].Pno != modp[i].Pno {
			t.Errorf("odp[%v].Pno=%v, want: %v", i, odp[i].Pno, modp[i].Pno)
		}
		if odp[i].State != modp[i].State {
			t.Errorf("odp[%v].State=%v, want: %v", i, odp[i].State, modp[i].State)
		}
		if odp[i].StateText != modp[i].StateText {
			t.Errorf("odp[%v].StateText=%v, want: %v", i, odp[i].StateText, modp[i].StateText)
		}
		if odp[i].CODAmount != modp[i].CODAmount {
			t.Errorf("odp[%v].CODAmount=%v, want: %v", i, odp[i].CODAmount, modp[i].CODAmount)
		}
		t.Logf("odp[%v] = %v\n", i, odp[i])
		t.Logf("modp[%v] = %v\n", i, modp[i])
	}
}
