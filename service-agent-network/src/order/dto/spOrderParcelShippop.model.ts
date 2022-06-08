export interface SpOrderParcelShippopModel {
    id: number,
    purchaseId: number;
    status: string;
    courierCode: string;
    courierTrackingCode: string;
    trackingCode: string;
    codAmount: number;
}