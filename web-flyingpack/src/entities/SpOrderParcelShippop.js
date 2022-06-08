export class SpOrderParcelShippop {
    constructor({
                    id,
                    purchaseId,
                    status,
                    courierCode,
                    courierTrackingCode,
                    trackingCode,
                    codAmount,
                }) {
        this.id = id || 0;
        this.purchaseId = purchaseId || 0;
        this.status = status || "";
        this.courierCode = courierCode || "";
        this.courierTrackingCode = courierTrackingCode || "";
        this.trackingCode = trackingCode || "";
        this.codAmount = codAmount || 0;
    }
}