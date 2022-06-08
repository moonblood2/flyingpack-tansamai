package order

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/internal/ordparflash"
	"github.com/jna-distribution/service-shipping/internal/ordproduct"
	"github.com/jna-distribution/service-shipping/internal/product"
	"github.com/jna-distribution/service-shipping/pkg"
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/jna-distribution/service-shipping/internal/destination"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/ordparcel"
	"github.com/jna-distribution/service-shipping/internal/ordparshippop"
	"github.com/jna-distribution/service-shipping/internal/origin"
	"github.com/jna-distribution/service-shipping/internal/provider"
	"github.com/jna-distribution/service-shipping/internal/sender"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func mockProduct(count int) []Product {
	products := make([]Product, count)
	for i := 0; i < count; i++ {
		products[i] = Product{
			Id:       6,
			Quantity: int32(1 + rand.Intn(4-1)),
		}
	}
	return products
}

func TestService_MakeOrderInputValidation(t *testing.T) {
	//Case 1, should have error.
	input := MakeOrderInput{
		PaymentMethod: 2,
		Products:      mockProduct(0),
	}
	err := input.Validate()
	if err == nil {
		t.Errorf("NO: %v, %v", input, err)
	}
	fmt.Printf("OK: %v, %v\n", input, err)
}

func TestService_MakeOrder(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	env, err := pkg.LoadEnv()

	productRepository := product.NewPostgresRepository(db)
	senderRepository := sender.NewPostgresRepository(db)
	originRepository := origin.NewPostgresRepository(db)
	destinationRepository := destination.NewPostgresRepository(db)
	orderParcelRepository := ordparcel.NewPostgresRepository(db)
	orderProductRepository := ordproduct.NewPostgresRepository(db)
	orderShippopRepository := ordparshippop.NewPostgresRepository(db)
	orderFlashRepository := ordparflash.NewPostgresRepository(db)
	providers := provider.NewProviderList(
		map[entity.ProviderCode]provider.Provider{
			entity.ProviderCodeShippop: provider.NewShippopProvider(
				env["PROVIdER_HTTP_API_PROXY_URI"],
				env["PROVIdER_SHIPPOP_API_URL"],
				env["PROVIDER_SHIPPOP_API_KEY"],
				orderShippopRepository,
			),
			entity.ProviderCodeFlash: provider.NewFlashProvider(
				env["PROVIDER_HTTP_API_PROXY_URI"],
				env["PROVIDER_FLASH_API_URL"],
				env["PROVIDER_FLASH_API_KEY"],
				env["PROVIDER_FLASH_MCH_ID"],
				orderFlashRepository,
			),
		},
	)
	jnaAddressAsOrigin := env["JNA_ADDRESS_AS_ORIGIN"] == "true"
	//JNA address
	jnaAddress := entity.Origin{ContactInfo: entity.ContactInfo{
		Name:        env["JNA_NAME"],
		PhoneNumber: env["JNA_PHONE_NUMBER"],
		Address:     env["JNA_ADDRESS"],
		District:    env["JNA_DISTRICT"],
		State:       env["JNA_STATE"],
		Province:    env["JNA_PROVINCE"],
		Postcode:    env["JNA_POSTCODE"],
	}}
	anUrl := env["AN_URL"]
	orderService := NewService(
		productRepository,
		senderRepository,
		originRepository,
		destinationRepository,
		orderParcelRepository,
		orderProductRepository,
		providers,
		jnaAddress,
		jnaAddressAsOrigin,
		anUrl,
	)

	parcelCount := 5
	productCount := 0
	rand.Seed(time.Now().Unix())
	input := MakeOrderInput{
		Sender:        mock.Sender(1)[0],
		PaymentMethod: mock.PaymentMethod(1)[0],
		Parcels:       mock.Parcel(parcelCount),
		Products:      mockProduct(productCount),
	}
	output, err := orderService.MakeOrder(input, "ddc5c679-6802-41e4-aa2b-00c3bdd3742b")
	t.Logf("input.Sender=%v\n", input.Sender)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("output.Status: %v\n", output.Status)
	for i := 0; i < len(output.Parcels); i++ {
		t.Logf("output.Parcels[%v]=%v\n", i, output.Parcels[i])
	}
	for i := 0; i < len(output.Products); i++ {
		t.Logf("output.Products[%v]=%v\n", i, output.Products[i])
	}
}
