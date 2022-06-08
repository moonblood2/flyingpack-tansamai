import {Parcel, parcelType} from "@/entities/Parcel";
import {AnBankAccount} from "@/entities/AnBankAccount";

//AnParcel is an order from Agent Network
export class AnParcel extends Parcel {
    constructor({
                    courierName,
                    price,
                    providerCode,
                    courierCode,
                    anCourierCode,
                    enableCOD,
                    codAmount,
                    trackingCode,
                    trackingCodes,
                    origin,
                    destination,
                    weight,
                    width,
                    length,
                    height,
                    fulfillmentStatus,
                    fulfillmentStatusString,
                    shippingStatus,
                    codStatus,
                    referenceNo,
                    spOrderParcelId,
                    anOrderId,
                    anOrderItemQuantitySum,
                    items,
                    anBankAccount,
                    createdAt,
                    statusCompletedDate,
                    codTransferredDate,
                    jnaCodTransferredDate,
                    spOrderParcelShippop,
                    spOrderParcelShippopFlash,
                }) {
        super({
            courierName: courierName,
            price: price,
            providerCode: providerCode,
            courierCode: courierCode,
            enableCOD: enableCOD,
            codAmount: codAmount,
            trackingCode: trackingCode,
            origin: origin,
            destination: destination,
            weight: weight,
            width: width,
            length: length,
            height: height,
            spOrderParcelShippop: spOrderParcelShippop,
            spOrderParcelShippopFlash: spOrderParcelShippopFlash,
        });
        this.type = parcelType.AN_PARCEL;
        this.anCourierCode = anCourierCode || 0;
        this.referenceNo = referenceNo || "";

        //Destination for agent-network
        this.desName = this.destination.name;
        this.desPhoneNumber = this.destination.phoneNumber;
        this.desAddress = this.destination.address;
        this.desSubdistrict = this.destination.district;
        this.desDistrict = this.destination.state;
        this.desProvince = this.destination.province;
        this.desPostcode = this.destination.postcode;

        this.fulfillmentStatus = fulfillmentStatus || 0;
        this.fulfillmentStatusString = fulfillmentStatusString || "";
        this.shippingStatus = shippingStatus || "";
        this.codStatus = codStatus || "";
        this.anOrderId = anOrderId || ""
        this.anOrderItemQuantitySum = anOrderItemQuantitySum || 0;
        this.items = items || [];
        this.createdAt = createdAt || "";
        this.statusCompletedDate = statusCompletedDate || "";
        this.codTransferredDate = codTransferredDate || "";
        this.jnaCodTransferredDate = jnaCodTransferredDate || "";
        this.spOrderParcelId = spOrderParcelId || null;
        this.anBankAccount = anBankAccount || new AnBankAccount({});

        this.trackingCodes = trackingCodes || [];
    }

    //clone
    clone() {
        return new AnParcel({
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
            height: this.height,
            anOrderId: this.anOrderId,
            anOrderItems: [...this.anOrderItems],
            createdAt: this.createdAt,
        })
    }

    //validate
    validate() {
        return super.validate();
    }
}