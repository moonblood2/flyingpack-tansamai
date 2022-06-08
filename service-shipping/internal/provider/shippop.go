package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/ordparshippop"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

type shippopProvider struct {
	proxyURI           string
	apiUrl             string
	apiKey             string
	orderParcelShippop ordparshippop.Repository
}

func NewShippopProvider(proxyURI, apiUrl string, apiKey string, orderParcelShippop ordparshippop.Repository) *shippopProvider {
	return &shippopProvider{
		proxyURI:           proxyURI,
		apiUrl:             apiUrl,
		apiKey:             apiKey,
		orderParcelShippop: orderParcelShippop,
	}
}

type (
	contactInfo struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"tel"`
		Address     string `json:"address"`
		District    string `json:"district"`
		State       string `json:"state"`
		Province    string `json:"province"`
		Postcode    string `json:"postcode"`
	}
	origin struct {
		contactInfo
	}
	destination struct {
		contactInfo
	}
	parcelShape struct {
		Weight float32 `json:"weight"`
		Width  float32 `json:"width"`
		Length float32 `json:"length"`
		Height float32 `json:"height"`
	}
	parcel struct {
		Origin      origin      `json:"from"`
		Destination destination `json:"to"`
		ParcelShape parcelShape `json:"parcel"`
		CourierCode string      `json:"courier_code"`
	}
)

//Get price: getPriceDataObj, courierDataObj, requestGetPriceData and responseGetPriceData use for GET PRICE API.
type (
	//request
	getPriceDataObj struct {
		parcel
		ShowAll int `json:"showall"`
	}
	requestGetPriceData struct {
		APIKey string                  `json:"api_key"`
		Data   map[int]getPriceDataObj `json:"data"`
	}
	//response
	courierDataObj struct {
		Price        string `json:"price"` //Price use type string, document tell integer but actually return string.
		EstimateTime string `json:"estimate_time"`
		Available    bool   `json:"available"`
		Remark       string `json:"remark"`
	}
	responseGetPriceData struct {
		Status bool                              `json:"status"`
		Code   int                               `json:"code"`
		Data   map[int]map[string]courierDataObj `json:"data"` //int as key, string as courier_code.
	}
)

//Booking order: bookingDataObj, bookingResponseObj, requestBookingOrderData and responseBookingOrderData use for BOOKING ORDER API.
type (
	//request
	bookingDataObj struct {
		parcel
		Remark    string `json:"remark,omitempty"`
		CODAmount int    `json:"cod_amount,omitempty"`
	}
	requestBookingOrderData struct {
		APIKey string                 `json:"api_key"`
		Email  string                 `json:"email"`
		Data   map[int]bookingDataObj `json:"data"` //Data Use key as int.
	}
	//response, check status first and then populate to true or false following by status.
	responseBookingOrderStatus struct {
		Status bool `json:"status"`
	}
	responseBookingOrderStatusTrueData struct {
		parcel
		CourierTrackingCode string  `json:"courier_tracking_code"`
		TrackingCode        string  `json:"tracking_code"`
		CODAmount           float64 `json:"cod_amount"`
	}
	responseBookingOrderStatusTrue struct {
		Status     bool                                       `json:"status"`
		Code       int                                        `json:"code"`
		Data       map[int]responseBookingOrderStatusTrueData `json:"data"`
		PurchaseID int64                                      `json:"purchase_id"`
		PaymentURL string                                     `json:"payment_url"`
		TotalPrice float64                                    `json:"total_price"`
	}
	responseBookingOrderStatusFalseData struct {
		Remark  string `json:"remark"`
		ErrCode string `json:"err_code"`
	}
	responseBookingOrderStatusFalse struct {
		Status bool                                                   `json:"status"`
		Data   map[int]map[string]responseBookingOrderStatusFalseData `json:"data"`
	}
	//use for return value.
	resultBookingOrderTrue struct {
		PurchaseID          int64   `json:"purchase_id"`
		CourierCode         string  `json:"courier_code"`
		CourierTrackingCode string  `json:"courier_tracking_code"`
		TrackingCode        string  `json:"tracking_code"`
		CODAmount           float64 `json:"cod_amount"`
	}
	resultBookingOrderFalse struct {
		Remark  string `json:"remark"`
		ErrCode string `json:"err_code"`
	}
	resultBookingOrder struct {
		Status   bool `json:"status"`
		ResTrue  resultBookingOrderTrue
		ResFalse resultBookingOrderFalse
	}
)

//Confirm purchase: requestConfirmPurchaseData and responseConfirmPurchaseData use for CONFIRM PURCHASE API.
type (
	requestConfirmPurchaseData struct {
		APIKey     string `json:"api_key"`
		PurchaseID int64  `json:"purchase_id"`
	}
	responseConfirmPurchaseData struct {
		Status  bool   `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

//Get Flash sorting code
type (
	//getFlashSortingResult key=trackingCode, value=sortingCode
	getFlashSortingResult map[string]entity.ShippopFlashSortingCode
)

//Get label
type (
	getLabelResponse struct {
		Status bool `json:"status"`
		Json   struct {
			Logo   string `json:"logo"`
			Labels []struct {
				TrackingCode        string `json:"trackingCode"`
				CourierTrackingCode string `json:"courierTrackingCode"`
				SortingCode         struct {
					SortCode        string `json:"sortCode"`
					DstCode         string `json:"dstCode"`
					SortingLineCode string `json:"sortingLineCode"`
				} `json:"sortingCode"`
			} `json:"labels"`
		} `json:"json"`
	}
)

// convertEntityParcel convert entity.Parcel to parcel, in some field use different JSON keys.
func convertEntityParcel(ep entity.Parcel) parcel {
	return parcel{
		Origin: origin{
			contactInfo: contactInfo{
				Name:        ep.Origin.Name,
				PhoneNumber: ep.Origin.PhoneNumber,
				Address:     ep.Origin.Address,
				District:    ep.Origin.District,
				State:       ep.Origin.State,
				Province:    ep.Origin.Province,
				Postcode:    ep.Origin.Postcode,
			},
		},
		Destination: destination{
			contactInfo: contactInfo{
				Name:        ep.Destination.Name,
				PhoneNumber: ep.Destination.PhoneNumber,
				Address:     ep.Destination.Address,
				District:    ep.Destination.District,
				State:       ep.Destination.State,
				Province:    ep.Destination.Province,
				Postcode:    ep.Destination.Postcode,
			},
		},
		ParcelShape: parcelShape{
			Weight: ep.ParcelShape.Weight,
			Width:  ep.ParcelShape.Width,
			Height: ep.ParcelShape.Height,
			Length: ep.ParcelShape.Length,
		},
		CourierCode: ep.CourierCode,
	}
}

//BookingOrder book only one order.
func (p *shippopProvider) BookingOrder(parcel entity.Parcel) (BookingOrderResult, error) {
	//Convert parcel to bookingDataObj.
	b := bookingDataObj{
		parcel: convertEntityParcel(parcel),
	}
	//Check request for COD.
	if parcel.EnableCOD {
		b.Remark = "เก็บเงินปลายทาง"
		b.CODAmount = int(parcel.CODAmount)
	}
	result := BookingOrderResult{}
	result.ProviderCode = parcel.ProviderCode
	//Send request to Booking Order API, retrieve purchase_id.
	resBooking, err := p.bookingOrder(b)
	if err != nil {
		return result, err
	}
	//Check if booking success then set resultShippop else set status to false and error message.
	if resBooking.Status == true {
		result.Status = true
		result.TrackingCode = resBooking.ResTrue.CourierTrackingCode
		result.ResultShippop = resultShippop{
			PurchaseID:          resBooking.ResTrue.PurchaseID,
			CourierCode:         resBooking.ResTrue.CourierCode,
			CourierTrackingCode: resBooking.ResTrue.CourierTrackingCode,
			TrackingCode:        resBooking.ResTrue.TrackingCode,
		}
	} else {
		result.Status = false
		result.Message = fmt.Sprintf("%s, %s", resBooking.ResFalse.ErrCode, resBooking.ResFalse.Remark)
	}
	return result, nil
}

//BookingOrders book many orders.
func (p *shippopProvider) BookingOrders(parcels []entity.Parcel) ([]BookingOrderResult, error) {
	count := len(parcels)
	bookingDataObjs := make([]bookingDataObj, count)
	//Convert parcels to bookingDataObjs.
	for i := 0; i < count; i++ {
		b := bookingDataObj{
			parcel: convertEntityParcel(parcels[i]),
		}
		//Check request for COD.
		if parcels[i].EnableCOD {
			b.Remark = "เก็บเงินปลายทาง"
			b.CODAmount = int(parcels[i].CODAmount)
		}
		bookingDataObjs[i] = b
	}
	resBookingOrders, err := p.bookingOrders(bookingDataObjs)
	if err != nil {
		return nil, err
	}
	results := make([]BookingOrderResult, count)
	for i := 0; i < count; i++ {
		if resBookingOrders[i].Status == true {
			results[i].Status = true
			results[i].ProviderCode = parcels[i].ProviderCode
			results[i].TrackingCode = resBookingOrders[i].ResTrue.CourierTrackingCode
			results[i].ResultShippop = resultShippop{
				PurchaseID:          resBookingOrders[i].ResTrue.PurchaseID,
				CourierCode:         resBookingOrders[i].ResTrue.CourierCode,
				CourierTrackingCode: resBookingOrders[i].ResTrue.CourierTrackingCode,
				TrackingCode:        resBookingOrders[i].ResTrue.TrackingCode,
			}
		} else {
			results[i].Status = false
			results[i].Message = fmt.Sprintf("%v: %v", resBookingOrders[i].ResFalse.ErrCode, resBookingOrders[i].ResFalse.Remark)
		}
		//Save key
		results[i].ParcelKey = parcels[i].Key
	}
	return results, nil
}

// ConfirmOrder confirm order, if courier is Flash, populate the Flash Sorting Code to BookingOrderResult
func (p *shippopProvider) ConfirmOrder(result *BookingOrderResult) error {
	if err := p.confirmPurchase(result.ResultShippop.PurchaseID); err != nil {
		return err
	}
	if result.ResultShippop.CourierCode == "FLE" {
		res, err := p.getFlashSortingCode(result.ResultShippop.PurchaseID, []string{result.ResultShippop.TrackingCode})
		if err != nil {
			return err
		}
		//populate the Flash Sorting Code to BookingOrderResult
		result.ResultShippop.FlashSortingCode = res[result.ResultShippop.TrackingCode]
	}
	return nil
}

// ConfirmOrders confirm orders, if courier is Flash, populate the Flash Sorting Code to BookingOrderResult
func (p *shippopProvider) ConfirmOrders(results []BookingOrderResult) error {
	purchaseIDSet := make(map[int64]bool)
	for _, v := range results {
		if _, ok := purchaseIDSet[v.ResultShippop.PurchaseID]; !ok {
			purchaseIDSet[v.ResultShippop.PurchaseID] = true
		}
	}
	for purchaseId := range purchaseIDSet {
		if err := p.confirmPurchase(purchaseId); err != nil {
			return err
		}
	}

	//Aggregate Flash, get Shippop's Tracking Code, get Flash sorting code, call label API.
	//A purchaseId for many tracking codes.

	//mapPurchaseID key=purchaseId, value=list of Flash's trackingCode.
	mapPurchaseIdTrackingCode := make(map[int64][]string)
	//memResultIndex key=trackingCode, value=index of results[]
	mapTrackingCodeResultIndex := make(map[string]int)
	for i, v := range results {
		if v.ResultShippop.CourierCode == "FLE" {
			mapPurchaseIdTrackingCode[v.ResultShippop.PurchaseID] = append(mapPurchaseIdTrackingCode[v.ResultShippop.PurchaseID], v.ResultShippop.TrackingCode)
			mapTrackingCodeResultIndex[v.ResultShippop.TrackingCode] = i
		}
	}

	for purchaseId, trackingCodes := range mapPurchaseIdTrackingCode {
		res, err := p.getFlashSortingCode(purchaseId, trackingCodes)
		if err != nil {
			return err
		}
		for trackingCode, flashSortingCode := range res {
			results[mapTrackingCodeResultIndex[trackingCode]].ResultShippop.FlashSortingCode = flashSortingCode
		}
	}

	return nil
}

//CancelOrder is empty, Shippop has no cancel order.
func (p *shippopProvider) CancelOrder(result BookingOrderResult) error {
	return nil
}

func (p *shippopProvider) GetPrice(parcel entity.Parcel) (GetPriceResult, error) {
	//Convert parcels to getPriceDataObj.
	data := make(map[int]getPriceDataObj, 1)
	objKey := 0
	data[objKey] = getPriceDataObj{
		parcel:  convertEntityParcel(parcel),
		ShowAll: 0, //ShowAll = 0, show only specific courier_code.
	}
	reqGetPriceData := requestGetPriceData{Data: data}
	resGetPriceData, err := p.getPrice(reqGetPriceData)
	if err != nil {
		return GetPriceResult{}, err
	}
	result := GetPriceResult{}
	//Check status from Booking Order API.
	if resGetPriceData.Status == true {
		result.Status = true
	} else {
		result.Status = false
	}
	o := resGetPriceData.Data[objKey][parcel.CourierCode]
	result.Status = result.Status && o.Available
	if o.Available {
		price, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			return GetPriceResult{}, err
		}
		result.Price = price
	} else {
		result.Message = o.Remark
	}
	return result, nil
}

//InsertOrder save order data to order_parcel_shippop (order details) and order_parcel_shippp_flash (Flash sorting code).
func (p *shippopProvider) InsertOrder(input []InsertOrderInput) error {
	orders := make([]entity.OrderParcelShippop, len(input))
	for i, v := range input {
		orders[i].OrderParcelId = v.OrderParcelId
		orders[i].PurchaseId = v.ResultShippop.PurchaseID
		orders[i].Status = ""
		orders[i].CourierCode = v.ResultShippop.CourierCode
		orders[i].CourierTrackingCode = v.ResultShippop.CourierTrackingCode
		orders[i].TrackingCode = v.ResultShippop.TrackingCode
		orders[i].CODAmount = v.ResultShippop.CODAmount
		orders[i].OrderParcelShippopFlash = entity.OrderParcelShippopFlash{
			ShippopFlashSortingCode: v.ResultShippop.FlashSortingCode,
		}
	}
	orders, err := p.orderParcelShippop.Insert(orders)
	if err != nil {
		return err
	}
	return nil
}

//UpdateOrderStatus update order status in order_parcel and order_parcel_shippop.
func (p *shippopProvider) UpdateOrderStatus(trackingCode, status, codStatus string) (orderParcelId string, statusCompletedDate, codTransferredDate *time.Time, err error) {
	return p.orderParcelShippop.UpdateOrderStatus(trackingCode, status, codStatus)
}

// getPrice request get price API, return prices in slice, order by the same as input.
func (p *shippopProvider) getPrice(reqData requestGetPriceData) (responseGetPriceData, error) {
	//Check has API key.
	if reqData.APIKey == "" {
		reqData.APIKey = p.apiKey
	}
	//Prepare url, method and payload.
	url := fmt.Sprintf("%v/pricelist/", p.apiUrl)
	method := http.MethodPost
	payload := &bytes.Buffer{}
	//Encode payload as JSON.
	err := json.NewEncoder(payload).Encode(&reqData)
	if err != nil {
		return responseGetPriceData{}, err
	}
	//Initialize client, if have proxy use it.
	client := &http.Client{}
	if p.proxyURI != "" {
		proxyUrl, err := neturl.Parse(p.proxyURI)
		if err != nil {
			return responseGetPriceData{}, err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}
	//Initial Request.
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return responseGetPriceData{}, err
	}
	//Set header, content-type = json
	req.Header.Add("Content-Type", "application/json")
	//Do request, receive response.
	res, err := client.Do(req)
	if err != nil {
		return responseGetPriceData{}, err
	}
	//Read response.
	defer res.Body.Close()
	resData := responseGetPriceData{}
	//Decode response JSON to data struct.
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return responseGetPriceData{}, err
	}
	return resData, nil
}

// bookingOrder request booking order API, return response data.
func (p *shippopProvider) bookingOrder(b bookingDataObj) (resultBookingOrder, error) {
	objKey := 0
	reqData := requestBookingOrderData{
		APIKey: p.apiKey,
		Email:  "acc.jna@gmail.com",
		Data:   map[int]bookingDataObj{objKey: b}, //Key is 0.
	}
	//Prepare url, method and payload.
	url := fmt.Sprintf("%v/booking/", p.apiUrl)
	method := http.MethodPost
	payload := &bytes.Buffer{}
	//Encode payload as JSON.
	err := json.NewEncoder(payload).Encode(&reqData)
	if err != nil {
		return resultBookingOrder{}, err
	}
	//Initialize client, if have proxy use it.
	client := &http.Client{}
	if p.proxyURI != "" {
		proxyUrl, err := neturl.Parse(p.proxyURI)
		if err != nil {
			return resultBookingOrder{}, err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}
	//Initial Request.
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return resultBookingOrder{}, err
	}
	//Set header, content-type = json
	req.Header.Add("Content-Type", "application/json")
	//Do request, receive response.
	res, err := client.Do(req)
	if err != nil {
		return resultBookingOrder{}, err
	}

	defer res.Body.Close()

	//Copy body to byte slice.
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return resultBookingOrder{}, err
	}
	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	body2 := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	defer body.Close()
	defer body2.Close()

	//Initialize response and result.
	resStatus := responseBookingOrderStatus{}
	resTrue := responseBookingOrderStatusTrue{}
	resFalse := responseBookingOrderStatusFalse{}
	resultBooking := resultBookingOrder{}
	//Decode status.
	if err := json.NewDecoder(body).Decode(&resStatus); err != nil {
		return resultBookingOrder{}, err
	}
	//Check status from Booking Order API.
	if resStatus.Status == true {
		//Decode True.
		if err := json.NewDecoder(body2).Decode(&resTrue); err != nil {
			return resultBookingOrder{}, err
		}
		resultBooking.Status = true
		resultBooking.ResTrue = resultBookingOrderTrue{
			PurchaseID:          resTrue.PurchaseID,
			CourierCode:         resTrue.Data[objKey].CourierCode,
			CourierTrackingCode: resTrue.Data[objKey].CourierTrackingCode,
			TrackingCode:        resTrue.Data[objKey].TrackingCode,
			//CODAmount:           resTrue.Data[objKey].CODAmount,
		}
	} else {
		//Decode False.
		if err := json.NewDecoder(body2).Decode(&resFalse); err != nil {
			return resultBookingOrder{}, err
		}
		resultBooking.Status = false
		resultBooking.ResFalse = resultBookingOrderFalse{
			Remark:  resFalse.Data[objKey][b.CourierCode].Remark,
			ErrCode: resFalse.Data[objKey][b.CourierCode].ErrCode,
		}
	}
	return resultBooking, nil
}

func (p *shippopProvider) bookingOrders(objs []bookingDataObj) ([]resultBookingOrder, error) {
	data := make(map[int]bookingDataObj)
	for i, v := range objs {
		data[i] = v
	}
	reqData := requestBookingOrderData{
		APIKey: p.apiKey,
		Email:  "acc.jna@gmail.com",
		Data:   data,
	}
	//Prepare url, method and payload.
	url := fmt.Sprintf("%v/booking/", p.apiUrl)
	method := http.MethodPost
	payload := &bytes.Buffer{}
	//Encode payload as JSON.
	err := json.NewEncoder(payload).Encode(&reqData)
	if err != nil {
		return nil, err
	}
	//Initialize client, if have proxy use it.
	client := &http.Client{}
	if p.proxyURI != "" {
		proxyUrl, err := neturl.Parse(p.proxyURI)
		if err != nil {
			return nil, err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}
	//Initial Request.
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	//Set header, content-type = json
	req.Header.Add("Content-Type", "application/json")
	//Do request, receive response.
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	//Check status
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(res.Body)
	resStatus := responseBookingOrderStatus{}
	if err := json.Unmarshal(bodyBytes, &resStatus); err != nil {
		return nil, err
	}

	resultBookingOrders := make([]resultBookingOrder, len(objs))

	if resStatus.Status {
		resTrue := responseBookingOrderStatusTrue{}
		if err := json.Unmarshal(bodyBytes, &resTrue); err != nil {
			return nil, err
		}
		for i := range objs {
			resultBookingOrders[i].Status = true
			resultBookingOrders[i].ResTrue = resultBookingOrderTrue{
				PurchaseID:          resTrue.PurchaseID,
				CourierCode:         resTrue.Data[i].CourierCode,
				CourierTrackingCode: resTrue.Data[i].CourierTrackingCode,
				TrackingCode:        resTrue.Data[i].TrackingCode,
			}
		}
	} else {
		resFalse := responseBookingOrderStatusFalse{}
		if err := json.Unmarshal(bodyBytes, &resFalse); err != nil {
			return nil, err
		}
		for i := range objs {
			resultBookingOrders[i].Status = false
			resultBookingOrders[i].ResFalse = resultBookingOrderFalse{
				ErrCode: resFalse.Data[i][objs[i].CourierCode].ErrCode,
				Remark:  resFalse.Data[i][objs[i].CourierCode].Remark,
			}
		}
	}
	return resultBookingOrders, nil
}

// confirmPurchase request confirm purchase API, return response data.
func (p *shippopProvider) confirmPurchase(purchaseID int64) error {
	reqData := requestConfirmPurchaseData{
		APIKey:     p.apiKey,
		PurchaseID: purchaseID,
	}
	//Prepare url, method and payload.
	url := fmt.Sprintf("%v/confirm/", p.apiUrl)
	method := http.MethodPost
	payload := bytes.Buffer{}
	//Encode reqData to JSON payload.
	if err := json.NewEncoder(&payload).Encode(reqData); err != nil {
		return err
	}
	//Initialize client, if have proxy use it.
	client := &http.Client{}
	if p.proxyURI != "" {
		proxyUrl, err := neturl.Parse(p.proxyURI)
		if err != nil {
			return err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}
	//Initial Request.
	req, err := http.NewRequest(method, url, &payload)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	//Decode JSON response to data struct.
	resData := responseConfirmPurchaseData{}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return err
	}
	//Check status from confirm purchase API, if failed.
	if resData.Status == false {
		return fmt.Errorf("confirm purchase failed: %+v", resData)
	}
	return nil
}

// getFlashSortingCode call getLabel, and retrieve only flash sorting code.
func (p *shippopProvider) getFlashSortingCode(purchaseId int64, trackingCodes []string) (getFlashSortingResult, error) {
	res, err := p.getLabel(purchaseId, trackingCodes, "json")
	if err != nil {
		return nil, err
	}
	result := getFlashSortingResult{}
	for _, v := range res.Json.Labels {
		result[v.TrackingCode] = entity.ShippopFlashSortingCode{
			SortCode:        v.SortingCode.SortCode,
			DstCode:         v.SortingCode.DstCode,
			SortingLineCode: v.SortingCode.SortingLineCode,
		}
	}
	return result, nil
}

// getLabel request label API, return response data.
func (p *shippopProvider) getLabel(purchaseId int64, trackingCodes []string, labelType string) (getLabelResponse, error) {
	url := fmt.Sprintf("%v/label/", p.apiUrl)
	method := http.MethodPost
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(struct {
		ApiKey       string `json:"api_key"`
		PurchaseId   int64  `json:"purchase_id"`
		TrackingCode string `json:"tracking_code"`
		Size         string `json:"size"`
		Type         string `json:"type"`
	}{
		ApiKey:       p.apiKey,
		PurchaseId:   purchaseId,
		TrackingCode: strings.Join(trackingCodes, ","),
		Type:         labelType,
	}); err != nil {
		return getLabelResponse{}, nil
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, &payload)
	if err != nil {
		return getLabelResponse{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return getLabelResponse{}, err
	}
	//Close body
	defer res.Body.Close()
	//Store in slice byte
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(res.Body)
	//Read status
	resStatus := struct {
		Status bool `json:"status"`
	}{}
	if err := json.Unmarshal(bodyBytes, &resStatus); err != nil {
		return getLabelResponse{}, err
	}

	resData := getLabelResponse{}
	//Check status
	if resStatus.Status {
		if err := json.Unmarshal(bodyBytes, &resData); err != nil {
			return getLabelResponse{}, nil
		}
	} else {
		return getLabelResponse{}, errors.New("Fail to get label. ")
	}

	return resData, nil
}
