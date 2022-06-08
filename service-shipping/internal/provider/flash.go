package provider

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/ordparflash"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type flashProvider struct {
	proxyURI                   string
	apiUrl                     string
	apiKey                     string
	mchId                      string
	orderParcelFlashRepository ordparflash.Repository
}

type (
	responseData struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	createOrderResponseData struct {
		Pno                string `json:"pno"`                //	Flash Express tracking number
		MchID              string `json:"mchId"`              //	merchant no
		OutTradeNo         string `json:"outTradeNo"`         //	merchant order no
		SortCode           string `json:"sortCode"`           //	Flash Express sorting code
		DstStoreName       string `json:"dstStoreName"`       //	Flash Store Name to deliver the parcel
		SortingLineCode    string `json:"sortingLineCode"`    //	Flash Express Line code
		EarlyFlightEnabled bool   `json:"earlyFlightEnabled"` //	whether it apply early flight service
		PackEnabled        bool   `json:"packEnabled"`        //	whether it need bagging
	}
	createOrderResponse struct {
		responseData
		Data createOrderResponseData `json:"data"`
	}
)

func NewFlashProvider(proxyURI, apiUrl, apiKey, mchId string, orderParcelFlashRepository ordparflash.Repository) *flashProvider {
	return &flashProvider{
		proxyURI:                   proxyURI,
		apiUrl:                     apiUrl,
		apiKey:                     apiKey,
		mchId:                      mchId,
		orderParcelFlashRepository: orderParcelFlashRepository,
	}
}

func (p *flashProvider) BookingOrder(parcel entity.Parcel) (BookingOrderResult, error) {
	data := map[string]interface{}{}
	nonce, err := p.randomNonce(10)
	if err != nil {
		return BookingOrderResult{}, err
	}
	book := BookingOrderResult{}

	data["mchId"] = p.mchId
	data["nonceStr"] = nonce
	data["expressCategory"] = 1 //Standard delivery let's see at http://open-docs.flashexpress.com/#express_category
	data["srcName"] = parcel.Origin.Name
	data["srcPhone"] = parcel.Origin.PhoneNumber
	data["srcProvinceName"] = parcel.Origin.Province
	data["srcCityName"] = parcel.Origin.State
	data["srcDistrictName"] = parcel.Origin.District
	data["srcPostalCode"] = parcel.Origin.Postcode
	data["srcDetailAddress"] = parcel.Origin.Address
	data["dstName"] = parcel.Destination.Name
	data["dstPhone"] = parcel.Destination.PhoneNumber
	data["dstHomePhone"] = parcel.Destination.PhoneNumber
	data["dstProvinceName"] = parcel.Destination.Province
	data["dstCityName"] = parcel.Destination.State
	data["dstDistrictName"] = parcel.Destination.District
	data["dstPostalCode"] = parcel.Destination.Postcode
	data["dstDetailAddress"] = parcel.Destination.Address
	data["articleCategory"] = 99 //Other let's see at http://open-docs.flashexpress.com/#article_category
	data["weight"] = parcel.ParcelShape.Weight
	data["width"] = parcel.ParcelShape.Weight
	data["length"] = parcel.ParcelShape.Length
	data["height"] = parcel.ParcelShape.Height
	data["insured"] = 0 //No
	if parcel.EnableCOD {
		data["codEnabled"] = 1 //Yes
		data["codAmount"] = parcel.CODAmount
		data["remark"] = "เก็บเงินปลายทาง"
		book.ResultFlash.CODAmount = parcel.CODAmount
	} else {
		data["codEnabled"] = 0 //No
	}
	reqParam := p.buildRequestParam(data)
	res := createOrderResponse{}
	endpointUrl := fmt.Sprintf("%v/orders", p.apiUrl)
	err = p.postRequest(endpointUrl, reqParam, &res)
	if err != nil {
		return BookingOrderResult{}, err
	}
	book.ProviderCode = parcel.ProviderCode
	//If success
	if res.Code == 1 {
		book.Status = true
		book.Message = res.Message
		book.ResultFlash.Pno = res.Data.Pno
		book.ResultFlash.MchID = res.Data.MchID
		book.ResultFlash.OutTradeNo = res.Data.OutTradeNo
		book.ResultFlash.SortCode = res.Data.SortCode
		book.ResultFlash.DstStoreName = res.Data.DstStoreName
		book.ResultFlash.SortingLineCode = res.Data.SortingLineCode
		book.ResultFlash.EarlyFlightEnabled = res.Data.EarlyFlightEnabled
		book.ResultFlash.PackEnabled = res.Data.PackEnabled
	} else {
		book.Status = false
		book.Message = res.Message
	}
	return book, nil
}

func (p *flashProvider) BookingOrders(parcels []entity.Parcel) ([]BookingOrderResult, error) {
	return nil, nil
}

func (p *flashProvider) ConfirmOrder(result *BookingOrderResult) error {
	return nil
}

func (p *flashProvider) ConfirmOrders(results []BookingOrderResult) error {
	return nil
}

func (p *flashProvider) CancelOrder(result BookingOrderResult) error {
	nonce, err := p.randomNonce(10)
	if err != nil {
		return err
	}
	endpointUrl := fmt.Sprintf("%v/orders/%v/cancel", p.apiUrl, result.ResultFlash.Pno) //%v = pno
	data := map[string]interface{}{}
	data["mchId"] = p.mchId
	data["nonceStr"] = nonce
	reqParam := p.buildRequestParam(data)
	res := responseData{}
	err = p.postRequest(fmt.Sprintf(endpointUrl, result.ResultFlash.Pno), reqParam, &res)
	if err != nil {
		return err
	}
	return nil
}

func (p *flashProvider) GetPrice(parcel entity.Parcel) (GetPriceResult, error) {
	res := GetPriceResult{}
	res.Status = true
	return res, nil
}

func (p *flashProvider) InsertOrder(input []InsertOrderInput) error {
	orders := make([]entity.OrderParcelFlash, len(input))
	for i, v := range input {
		orders[i].OrderParcelId = v.OrderParcelId
		orders[i].Pno = v.ResultFlash.Pno
		orders[i].CODAmount = v.ResultFlash.CODAmount
	}
	orders, err := p.orderParcelFlashRepository.Insert(orders)
	if err != nil {
		return err
	}
	return nil
}

func (p *flashProvider) UpdateOrderStatus(trackingCode, status, codStatus string) (orderParcelId string, statusCompletedDate, codTransferredDate *time.Time, err error) {
	return "", &time.Time{}, &time.Time{}, nil
}

func (p *flashProvider) randomNonce(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	//1 Byte = 2 Char, take only 10 chars heading.
	return fmt.Sprintf("%X", b)[:n], nil
}

func (p *flashProvider) buildRequestParam(dataMap map[string]interface{}) url.Values {
	dataLen := len(dataMap)
	//keys contains keys in dataMap. It is parameter's name.
	keys := make([]string, dataLen)
	i := 0
	for k, _ := range dataMap {
		keys[i] = k
		i++
	}
	//Sort by the parameter's name with its ASCII code ASC.
	sort.Strings(keys)
	//Format the parameters with the format of key=value.
	var stringSignTemp strings.Builder
	for i := 0; i < dataLen; i++ {
		stringSignTemp.WriteString(keys[i])
		stringSignTemp.WriteString("=")
		stringSignTemp.WriteString(fmt.Sprintf("%v", dataMap[keys[i]]))
		stringSignTemp.WriteString("&")
	}
	//Join API Key.
	stringSignTemp.WriteString("key=")
	stringSignTemp.WriteString(p.apiKey)
	signByte := sha256.Sum256([]byte(stringSignTemp.String()))
	sign := fmt.Sprintf("%X", signByte)
	//Add sign.
	dataMap["sign"] = sign
	//Build request param.
	params := url.Values{}
	for k, v := range dataMap {
		params.Add(k, fmt.Sprintf("%v", v))
	}
	return params
}

func (p *flashProvider) postRequest(endpointUrl string, requestParams url.Values, responsePtr interface{}) error {
	client := &http.Client{}
	if p.proxyURI != "" {
		proxyURL, err := url.Parse(p.proxyURI)
		if err != nil {
			return err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	req, err := http.NewRequest(http.MethodPost, endpointUrl, strings.NewReader(requestParams.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(responsePtr)
	if err != nil {
		return err
	}
	return nil
}
