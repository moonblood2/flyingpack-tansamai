package provider

import (
	"database/sql"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/ordparflash"
	"testing"
)

func TestProviderList_Get(t *testing.T) {
	repo := ordparflash.NewPostgresRepository(&sql.DB{})
	providerList := NewProviderList(map[entity.ProviderCode]Provider{
		entity.ProviderCodeFlash: NewFlashProvider("", "", "", "", repo),
	})
	//Case 1: should success.
	_, err := providerList.Get(entity.ProviderCodeFlash)
	if err != nil {
		t.Errorf("NO: %v", err)
		return
	}
	fmt.Printf("OK: %v\n", err)
	//Case 2: should fail.
	_, err = providerList.Get(entity.ProviderCodeShippop)
	if err == nil {
		t.Errorf("NO: %v", err)
		return
	}
	fmt.Printf("OK: %v\n", err)
}
