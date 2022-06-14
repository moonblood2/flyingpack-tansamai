import {ContactInfo} from './ContactInfo'
import {validateObject, VError} from "@/utils/joi-validation";
import Joi from "joi";
import {SpOrderParcelShippop} from "@/entities/SpOrderParcelShippop";
import {SpOrderParcelShippopFlash} from "@/entities/SpOrderParcelShippopFlash";

export const parcelType = {
    PARCEL: 1,
    AN_PARCEL: 2,
}

export const typeOfParcel = (parcel) => {
    if (parcel.type) {
        if (parcel.type === parcelType.AN_PARCEL) {
            return parcelType.AN_PARCEL;
        } else if (parcel.type === parcelType.PARCEL) {
            return parcelType.PARCEL;
        } else {
            console.warn("This parcel has no parcel type.")
        }
    } else {
        console.warn("This parcel has no parcel type.")
    }
}

export class Parcel {
    constructor({
                    courierName,
                    price,
                    providerCode,
                    courierCode,
                    enableCOD,
                    codAmount,
                    trackingCode,
                    origin,
                    destination,
                    weight,
                    width,
                    length,
                    height,
                    spOrderParcelShippop,
                    spOrderParcelShippopFlash,
                    sortCode,
                    lineCode,
                    sortingLineCode,
                    dstStoreName,
                }) {
        this.type = parcelType.PARCEL;
        this.courierName = courierName || ""
        this.price = price || 0.0// float64
        this.providerCode = providerCode || 0 // ProviderCode
        this.courierCode = courierCode || ""  // string
        this.enableCOD = enableCOD || false
        this.codAmount = codAmount || 0.0 // float64
        this.trackingCode = trackingCode || ""
        this.sortCode = sortCode
        this.lineCode = lineCode
        this.sortingLineCode = sortingLineCode
        this.dstStoreName = dstStoreName
        //Create new one for deep clone instance.
        this.origin = origin ? new ContactInfo({
            name: origin.name,
            phoneNumber: origin.phoneNumber,
            address: origin.address,
            district: origin.district,
            state: origin.state,
            province: origin.province,
            postcode: origin.postcode,
        }) : new ContactInfo({});
        this.destination = destination ? new ContactInfo({
            name: destination.name,
            phoneNumber: destination.phoneNumber,
            address: destination.address,
            district: destination.district,
            state: destination.state,
            province: destination.province,
            postcode: destination.postcode,
        }) : new ContactInfo({});
        this.weight = weight || 0.0
        this.width = width || 0.0
        this.length = length || 0.0
        this.height = height || 0.0

        this.spOrderParcelShippop = spOrderParcelShippop || new SpOrderParcelShippop({});
        this.spOrderParcelShippopFlash = spOrderParcelShippopFlash || new SpOrderParcelShippopFlash({});

        this.valid = false
        this.error = {
            name: new VError(),
            price: new VError(),
            providerCode: new VError(),
            courierCode: new VError(),
            enableCOD: new VError(),
            codAmount: new VError(),
            trackingCode: new VError(),

            weight: new VError(),
            width: new VError(),
            length: new VError(),
            height: new VError(),
        }
    }

    //clone
    clone() {
        return new Parcel({
            name: this.name,
            price: this.price,
            providerCode: this.providerCode,
            courierCode: this.courierCode,
            enableCOD: this.enableCOD,
            codAmount: this.codAmount,
            origin: this.origin.clone(),
            destination: this.destination.clone(),
            weight: this.weight,
            width: this.width,
            length: this.length,
            height: this.height
        })
    }

    //validate
    validate() {
        const {valid, error} = validateObject({
            price: Joi.number().optional(),
            providerCode: Joi.number().min(1).required(),
            courierCode: Joi.string().optional(),
            codAmount: Joi.number().optional(),
            weight: Joi.number().min(1).required(),
            width: Joi.number().min(1).required(),
            length: Joi.number().min(1).required(),
            height: Joi.number().min(1).required(),
        }, {
            price: this.price,
            providerCode: this.providerCode,
            courierCode: this.courierCode,
            codAmount: this.codAmount,
            weight: this.weight,
            width: this.width,
            length: this.length,
            height: this.height,
        })
        const originValidate = this.origin.validate();
        const destinationValidate = this.destination.validate();
        this.valid = valid && originValidate.valid && destinationValidate.valid
        this.error = {
            ...this.error,
            ...error,
        }
        return {valid: this.valid, error: this.error}
    }

    //toPayload
    toPayload() {
        return {
            "provider_code": parseInt(this.providerCode),
            "courier_code": this.courierCode,
            "enable_cod": this.enableCOD,
            "cod_amount": parseFloat(this.codAmount),
            "origin": this.origin.toPayload(),
            "destination": this.destination.toPayload(),
            "parcel_shape": {
                "weight": parseFloat(this.weight),
                "width": parseFloat(this.width),
                "length": parseFloat(this.length),
                "height": parseFloat(this.height)
            },
        }
    }
}