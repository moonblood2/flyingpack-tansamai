package provider

import (
	"database/sql"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/internal/ordparshippop"
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

func createMockParcels(normalCount, codCount int) []entity.Parcel {
	contactOri := entity.ContactInfo{
		Name:        "สวนกล้วยไม้ทันหิกรณ์",
		PhoneNumber: "022700499",
		Address:     "101-120 อาคารสยามออร์คิด เซนเตอร์ ตลาด อตก Kamphaeng Phet Rd",
		District:    "คลองต้นไทร",
		State:       "คลองสาน",
		Province:    "กรุงเทพมหานคร",
		Postcode:    "10600",
	}
	contactDes := entity.ContactInfo{
		Name:        "Niyom’s Orchid Farm",
		PhoneNumber: "0944785260",
		Address:     "1 109/1 Soi Kamnan",
		District:    "สมเด็จเจ้าพระยา",
		State:       "คลองสาน",
		Province:    "กรุงเทพมหานคร",
		Postcode:    "10600",
	}
	info := entity.ParcelShape{
		Weight: 10,
		Width:  10,
		Height: 10,
		Length: 10,
	}
	ori := entity.Origin{
		ContactInfo: contactOri,
	}
	des := entity.Destination{
		ContactInfo: contactDes,
	}
	var parcels []entity.Parcel
	for i := 0; i < normalCount; i++ {
		parcels = append(parcels, entity.Parcel{
			CourierCode: "FLE",
			ParcelShape: info,
			Origin:      ori,
			Destination: des,
		})
	}
	for i := 0; i < codCount; i++ {
		parcels = append(parcels, entity.Parcel{
			ProviderCode: entity.ProviderCodeShippop,
			CourierCode:  "FLE",
			ParcelShape:  info,
			Origin:       ori,
			Destination:  des,
			EnableCOD:    true,
			CODAmount:    100,
		})
	}

	return parcels
}

func createMockFailParcels(normalCount, codCount int) []entity.Parcel {
	contactOri := entity.ContactInfo{
		Name:        "สวนกล้วยไม้ทันหิกรณ์",
		PhoneNumber: "022700499",
		Address:     "101-120 อาคารสยามออร์คิด เซนเตอร์ ตลาด อตก Kamphaeng Phet Rd",
		District:    "คลองต้นไทร",
		State:       "คลองสาน",
		Province:    "กรุงเทพมหานคร",
		Postcode:    "1060", //set postcode to 4 digit for fail.
	}
	contactDes := entity.ContactInfo{
		Name:        "Niyom’s Orchid Farm",
		PhoneNumber: "0944785260",
		Address:     "1 109/1 Soi Kamnan",
		District:    "สมเด็จเจ้าพระยา",
		State:       "คลองสาน",
		Province:    "กรุงเทพมหานคร",
		Postcode:    "10600",
	}
	info := entity.ParcelShape{
		Weight: 10,
		Width:  10,
		Height: 10,
		Length: 10,
	}
	ori := entity.Origin{
		ContactInfo: contactOri,
	}
	des := entity.Destination{
		ContactInfo: contactDes,
	}
	var parcels []entity.Parcel
	for i := 0; i < normalCount; i++ {
		parcels = append(parcels, entity.Parcel{
			CourierCode: "FLE",
			ParcelShape: info,
			Origin:      ori,
			Destination: des,
		})
	}
	for i := 0; i < codCount; i++ {
		parcels = append(parcels, entity.Parcel{
			CourierCode: "FLE",
			ParcelShape: info,
			Origin:      ori,
			Destination: des,
			EnableCOD:   true,
			CODAmount:   100,
		})
	}

	return parcels
}

func TestShippopProvider_BookingOrder(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	//Case 1
	parcel := createMockParcels(1, 0)[0]
	result, err := p.BookingOrder(parcel)
	if err != nil {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	if result.Status != true {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	if result.Message != "" {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	fmt.Printf("OK: %+v\n", result)

	//Case 2
	parcel = createMockFailParcels(1, 0)[0]
	result, err = p.BookingOrder(parcel)
	if err != nil {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	if result.Status != false {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	if result.Message == "" {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	fmt.Printf("OK: %+v\n", result)
}

func TestShippopProvider_BookingOrder_ConfirmOrder(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	//Case 1
	parcel := createMockParcels(1, 0)[0]
	result, err := p.BookingOrder(parcel)
	if err != nil {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	if result.Status != true {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	if result.Message != "" {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	t.Logf("OK: %+v\n", result)
	err = p.ConfirmOrder(&result)
	if err != nil {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	t.Logf("OK: %+v\n", result)
}

func TestShippopProvider_BookingOrders(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	//Case 1
	count := 100
	parcels := createMockParcels(count, 0)
	result, err := p.BookingOrders(parcels)
	if err != nil {
		t.Errorf("NO: %+v, %v", result, err)
		return
	}
	for i := 0; i < count; i++ {
		log.Printf("result[%v] = %#v", i, result[i])
	}
}

func TestShippopProvider_BookingOrders_ConfirmOrders(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	normalCount := 10
	parcels := createMockParcels(normalCount, 0)
	results, err := p.BookingOrders(parcels)
	if err != nil {
		t.Errorf("NO: %+v, %v", results, err)
		return
	}
	if err := p.ConfirmOrders(results); err != nil {
		t.Errorf("NO: %+v, %v", results, err)
		return
	}
	for i := 0; i < normalCount; i++ {
		if results[i].ProviderCode == entity.ProviderCodeShippop {
			if results[i].ResultShippop.CourierCode == "FLE" {
				if results[i].ResultShippop.FlashSortingCode.SortingLineCode == "" {
					t.Errorf("NO: FlashSortingCode.SortingLineCode: %v\n", results[i].ResultShippop.FlashSortingCode.SortingLineCode)
				}
				if results[i].ResultShippop.FlashSortingCode.DstCode == "" {
					t.Errorf("NO: FlashSortingCode.DstCode: %v\n", results[i].ResultShippop.FlashSortingCode.DstCode)
				}
				if results[i].ResultShippop.FlashSortingCode.SortCode == "" {
					t.Errorf("NO: FlashSortingCode.SortCode: %v\n", results[i].ResultShippop.FlashSortingCode.SortCode)
				}
			}
		}
		log.Printf("result[%v] = %#v", i, results[i])
	}
}

func TestShippopProvider_InsertOrder(t *testing.T) {
	db, err := pkg.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	count := 5
	ord := mock.OrderParcelShippop(count, []string{""})
	input := make([]InsertOrderInput, count)
	for i := 0; i < count; i++ {
		input[i] = InsertOrderInput{
			BookingOrderResult: BookingOrderResult{
				ResultShippop: resultShippop{
					PurchaseID:          ord[i].PurchaseId,
					CourierCode:         ord[i].CourierCode,
					CourierTrackingCode: ord[i].CourierTrackingCode,
					TrackingCode:        ord[i].CourierTrackingCode,
					CODAmount:           ord[i].CODAmount,
				},
			},
		}
	}
	if err := p.InsertOrder(input); err != nil {
		t.Fatal(err)
	}
}

func TestShippopProvider_GetPrice(t *testing.T) {
	parcels := createMockParcels(1, 0)
	data := make(map[int]getPriceDataObj, 1)
	objKey := 0
	data[objKey] = getPriceDataObj{
		parcel:  convertEntityParcel(parcels[0]),
		ShowAll: 0,
	}
	reqGetPriceData := requestGetPriceData{
		Data: data,
	}

	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	resGetPriceData, err := p.getPrice(reqGetPriceData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resGetPriceData: %+v\n", resGetPriceData)
}

func TestShippopProvider_bookingOrder(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	var parcel entity.Parcel

	//Case 1: should OK.
	parcel = createMockParcels(1, 0)[0]
	//Convert parcel to bookingDataObj.
	b := bookingDataObj{
		parcel: convertEntityParcel(parcel),
	}
	//Check request for COD.
	if parcel.EnableCOD {
		b.Remark = "เก็บเงินปลายทาง"
		b.CODAmount = int(parcel.CODAmount)
	}
	//Send request to Booking Order API, retrieve purchase_id.
	resBooking, err := p.bookingOrder(b)
	if err != nil {
		t.Errorf("NO: %+v %v\n", resBooking, err)
		return
	}
	if resBooking.Status != true {
		t.Errorf("NO: %+v %v\n", resBooking, err)
		return
	}
	fmt.Printf("OK: %+v\n", resBooking)

	//Case 1: should NO.
	parcel = createMockFailParcels(1, 0)[0]
	//Convert parcel to bookingDataObj.
	b = bookingDataObj{
		parcel: convertEntityParcel(parcel),
	}
	//Check request for COD.
	if parcel.EnableCOD {
		b.Remark = "เก็บเงินปลายทาง"
		b.CODAmount = int(parcel.CODAmount)
	}
	//Send request to Booking Order API, retrieve purchase_id.
	resBooking, err = p.bookingOrder(b)
	if err != nil {
		t.Errorf("NO: %+v %v\n", resBooking, err)
		return
	}
	if resBooking.Status != false {
		t.Errorf("NO: %+v %v\n", resBooking, err)
		return
	}
	fmt.Printf("OK: %+v\n", resBooking)
}

func TestShippopProvider_confirmPurchase(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	err = p.confirmPurchase(184711)
	if err != nil {
		t.Errorf("NO: %v\n", err)
	}
	fmt.Printf("OK: %v\n", err)
}

func TestShippopProvider_getLabel(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewShippopProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_SHIPPOP_API_URL"],
		env["PROVIDER_SHIPPOP_API_KEY"],
		ordparshippop.NewPostgresRepository(&sql.DB{}),
	)

	res, err := p.getLabel(193806, []string{"SP049070275", "SP049070280"}, "json")
	if err != nil {
		t.Errorf("NO: %v\n", err)
	}
	t.Logf("OK: %+v\n", res)
}
