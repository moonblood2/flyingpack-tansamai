import {AnParcel} from "@/entities/AnParcel";
import {BankAccount} from "@/entities/BankAccount";
import {CourierCodeToCourierName} from "@/entities/AnCourier";
import {SpOrderParcelShippop} from "@/entities/SpOrderParcelShippop";
import {SpOrderParcelShippopFlash} from "@/entities/SpOrderParcelShippopFlash";
import provider from "@/entities/Provider";
import {ContactInfo} from "@/entities/ContactInfo";

export const parseAnParcel = (e) => {
    let anParcel = new AnParcel({
        courierName: CourierCodeToCourierName(e["courierCode"]),
        price: 0,
        providerCode: provider.PROVIDER_SHIPPOP, //Shippop
        courierCode: e["courierCode"],
        enableCOD: parseInt(e["codAmount"]) > 0,
        codAmount: parseInt(e["codAmount"]),
        trackingCode: e["trackingCode"],
        sortCode: e["sortCode"],
        lineCode: e["lineCode"],
        sortingLineCode: e["sortingLineCode"],
        dstStoreName: e["dstStoreName"],
        // trackingCodes: e["trackingCodes"],
        origin: new ContactInfo({}),
        destination: new ContactInfo({
            name: e["desName"],
            phoneNumber: e["desPhoneNumber"],
            address: e["desAddress"],
            district: e["desSubdistrict"],
            state: e["desDistrict"],
            province: e["desProvince"],
            postcode: e["desPostcode"],
        }),
        referenceNo: e["referenceNo"],
        spOrderParcelId: e["spOrderParcelId"],
        anOrderId: e["anOrderId"],
        items: e["items"],
        anOrderItemQuantitySum: e["itemQuantitySum"],
        createdAt: e["createdAt"],
        statusCompletedDate: e["statusCompletedDate"],
        codTransferredDate: e["codTransferredDate"],
        jnaCodTransferredDate: e["jnaCodTransferredDate"],
        fulfillmentStatus: e["fulfillmentStatus"],
        fulfillmentStatusString: e["fulfillmentStatusString"],
        shippingStatus: e["shippingStatus"],
        codStatus: e["codStatus"],
    })
    if (e["spOrderParcel"]) {
        //Populate dimension and weight.
        anParcel.weight = e["spOrderParcel"]["weight"];
        anParcel.width = e["spOrderParcel"]["width"];
        anParcel.length = e["spOrderParcel"]["length"];
        anParcel.height = e["spOrderParcel"]["height"];
    }
    if (e["spOrderParcelShippop"]) {
        anParcel.spOrderParcelShippop = new SpOrderParcelShippop({
            id: e["spOrderParcelShippop"]["id"],
            purchaseId: e["spOrderParcelShippop"]["purchaseId"],
            status: e["spOrderParcelShippop"]["status"],
            courierCode: e["spOrderParcelShippop"]["courierCode"],
            courierTrackingCode: e["spOrderParcelShippop"]["courierTrackingCode"],
            trackingCode: e["spOrderParcelShippop"]["trackingCode"],
            codAmount: e["spOrderParcelShippop"]["codAmount"],
        });
    }
    if (e["spOrderParcelShippopFlash"]) {
        anParcel.spOrderParcelShippopFlash = new SpOrderParcelShippopFlash({
            sortCode: e["spOrderParcelShippopFlash"]["sortCode"],
            dstCode: e["spOrderParcelShippopFlash"]["dstCode"],
            sortingLineCode: e["spOrderParcelShippopFlash"]["sortingLineCode"],
        });
    }
    if (e["anOrderPrice"]) {
        anParcel.anOrderPrice = {
            parcelCost: e["anOrderPrice"]["parcelCost"],
            parcelRemoteAreaPrice: e["anOrderPrice"]["parcelRemoteAreaPrice"],
            parcelSellingPrice: e["anOrderPrice"]["parcelSellingPrice"],
            fulfillmentServiceCharge: e["anOrderPrice"]["fulfillmentServiceCharge"]
        }
    }
    if (e["bankAccount"]) {
        anParcel.bankAccount = new BankAccount({
            bank: e["bankAccount"]["bank"],
            accountName: e["bankAccount"]["accountName"],
            accountNo: e["bankAccount"]["accountNo"],
            email: e["bankAccount"]["email"],
            fiCode: e["bankAccount"]["fiCode"],
            fiName: e["bankAccount"]["fiName"],
        });
    }

    return anParcel;
}