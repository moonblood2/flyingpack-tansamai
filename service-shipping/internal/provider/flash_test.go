package provider

import (
	"database/sql"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/mock"
	"github.com/jna-distribution/service-shipping/internal/ordparflash"
	"github.com/jna-distribution/service-shipping/pkg"
	"net/url"
	"testing"
	"time"
)

func TestFlashProvider_randomNonce(t *testing.T) {
	p := NewFlashProvider("", "", "", "", ordparflash.NewPostgresRepository(&sql.DB{}))
	nonceStr, err := p.randomNonce(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("p.randomNonce(10)=%v\n", nonceStr)
	nonceStr, err = p.randomNonce(5)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("p.randomNonce(5)=%v\n", nonceStr)
}

func TestFlashProvider_buildRequestParam(t *testing.T) {
	p := NewFlashProvider(
		"",
		"",
		"",
		"",
		ordparflash.NewPostgresRepository(&sql.DB{}),
	)
	//c := 2
	var dataTable []map[string]interface{}
	var want []url.Values
	dataTable = append(dataTable, map[string]interface{}{
		"mchId":    "5a7bdfd22593414adb72df5f",
		"nonceStr": "yyv6YJP436wCkdpNdghC",
		"body":     "test",
	})
	dataTable = append(dataTable, map[string]interface{}{
		"mchId":    "UBP18100020",
		"nonceStr": "1525314174723",
	})
	dataTable = append(dataTable, map[string]interface{}{
		"mchId":              p.mchId,
		"nonceStr":           "1536749552628",
		"srcName":            "หอมรวม  create order test name",
		"srcPhone":           "0630101454",
		"srcProvinceName":    "อุบลราชธานี",
		"srcCityName":        "เมืองอุบลราชธานี",
		"srcDistrictName":    "ในเมือง",
		"srcPostalCode":      "34000",
		"srcDetailAddress":   "68/5-6 ม.1 บ้านท่าบ่อ create order test address",
		"dstName":            "น้ำพริกแม่อำพร",
		"dstPhone":           "0970209976",
		"dstHomePhone":       "0970220220",
		"dstProvinceName":    "เชียงใหม่",
		"dstCityName":        "สันทราย",
		"dstDistrictName":    "สันพระเนตร",
		"dstPostalCode":      "50210",
		"dstDetailAddress":   "127 หมู่ 3 ต.หนองแหย่ง อ.สันทราย จ.เชียงใหม่ create order test address",
		"articleCategory":    "1",
		"expressCategory":    "1",
		"weight":             "1000",
		"insured":            "1",
		"insureDeclareValue": "10000",
		"codEnabled":         "1",
		"codAmount":          "10000",
		"remark":             "ขึ้นบันได",
	})

	want = append(want, url.Values{
		"mchId":    []string{"5a7bdfd22593414adb72df5f"},
		"nonceStr": []string{"yyv6YJP436wCkdpNdghC"},
		"body":     []string{"test"},
		"sign":     []string{"AC7B837F4AFE34EB23CC8E32FD5F1745F2A25A963CF5B2D0FB42ACD6B165C4C8"},
	})
	want = append(want, url.Values{
		"mchId":    []string{"UBP18100020"},
		"nonceStr": []string{"1525314174723"},
		"sign":     []string{"998A5030D4F0F5F35DFA634B0E22086A16BB0799B5CA7D031E06F68C10C820EC"},
	})
	/*
		for i := 0; i < c; i++ {
			params := p.buildRequestParam(dataTable[i])
			for k := range params {
				if params.Get(k) != want[i].Get(k) {
					t.Errorf("NO: p.makeSignature(dataTable[%v]).Get(%v)=%v, %v", i, k, params.Get(k), want[i].Get(k))
				}
			}
			t.Logf("OK, [%v]: %v\n", i, params.Encode())
		}
	*/
	params := p.buildRequestParam(dataTable[2])
	for k := range params {
		fmt.Printf("%v: %v\n", k, params.Get(k))
	}
}

func TestFlashProvider_postRequest(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewFlashProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_FLASH_API_URL"],
		env["PROVIDER_FLASH_API_KEY"],
		env["PROVIDER_FLASH_MCH_ID"],
		ordparflash.NewPostgresRepository(&sql.DB{}),
	)
	data := map[string]interface{}{
		"mchId":              p.mchId,
		"nonceStr":           fmt.Sprintf("%v", time.Now().Unix()),
		"srcName":            "หอมรวม  create order test name",
		"srcPhone":           "0630101454",
		"srcProvinceName":    "อุบลราชธานี",
		"srcCityName":        "เมืองอุบลราชธานี",
		"srcDistrictName":    "ในเมือง",
		"srcPostalCode":      "34000",
		"srcDetailAddress":   "68/5-6 ม.1 บ้านท่าบ่อ create order test address",
		"dstName":            "น้ำพริกแม่อำพร",
		"dstPhone":           "0970209976",
		"dstHomePhone":       "0970220220",
		"dstProvinceName":    "เชียงใหม่",
		"dstCityName":        "สันทราย",
		"dstDistrictName":    "สันพระเนตร",
		"dstPostalCode":      "50210",
		"dstDetailAddress":   "127 หมู่ 3 ต.หนองแหย่ง อ.สันทราย จ.เชียงใหม่ create order test address",
		"articleCategory":    "1",
		"expressCategory":    "1",
		"weight":             "1000",
		"insured":            "1",
		"insureDeclareValue": "10000",
		"codEnabled":         "1",
		"codAmount":          "10000",
		"remark":             "ขึ้นบันได",
	}
	params := p.buildRequestParam(data)
	res := createOrderResponse{}
	err = p.postRequest("http://api-training.flashexpress.com/open/v1/orders", params, &res)
	if err != nil {
		t.Errorf("NO: %v", err)
	}
	t.Logf("OK: %+v\n", res)
}

func TestFlashProvider_BookingOrder(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewFlashProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_FLASH_API_URL"],
		env["PROVIDER_FLASH_API_KEY"],
		env["PROVIDER_FLASH_MCH_ID"],
		ordparflash.NewPostgresRepository(&sql.DB{}),
	)
	parcel := entity.Parcel{
		ProviderCode: entity.ProviderCodeFlash,
		EnableCOD:    false,
		CODAmount:    0,
		Origin:       mock.Origin(1)[0],
		Destination:  mock.Destination(1)[0],
		ParcelShape:  mock.ParcelShape(1)[0],
	}
	result, err := p.BookingOrder(parcel)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result: %#v\n", result)
}

func TestFlashProvider_CancelOrder(t *testing.T) {
	env, err := pkg.LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	p := NewFlashProvider(
		env["PROVIDER_HTTP_API_PROXY_URI"],
		env["PROVIDER_FLASH_API_URL"],
		env["PROVIDER_FLASH_API_KEY"],
		env["PROVIDER_FLASH_MCH_ID"],
		ordparflash.NewPostgresRepository(&sql.DB{}),
	)
	book := BookingOrderResult{ResultFlash: resultFlash{
		Pno: "TH01293GQS4A",
	}}
	err = p.CancelOrder(book)
	if err != nil {
		t.Fatal(err)
	}
}
