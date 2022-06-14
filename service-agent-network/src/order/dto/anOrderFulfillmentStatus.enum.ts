export enum AnOrderFulfillmentStatusEnum {
    NOT_SET,
    PACKED,
    NOT_PACKED,
    CANCEL,
}

export const AnOrderFulfillmentStatusToString = {
    [AnOrderFulfillmentStatusEnum.NOT_SET]: "-",
    [AnOrderFulfillmentStatusEnum.PACKED]: "จัดส่งแล้ว",
    [AnOrderFulfillmentStatusEnum.NOT_PACKED]: "ยังไม่จัดส่ง",
    [AnOrderFulfillmentStatusEnum.CANCEL]: "ยกเลิก",
}

export const parseAnOrderFulfillmentStatusEnum = (status: string) => {
    switch (status) {
        case "0" : {
            return AnOrderFulfillmentStatusEnum.NOT_SET;
        }
        case "1": {
            return AnOrderFulfillmentStatusEnum.PACKED
        }
        case "2": {
            return AnOrderFulfillmentStatusEnum.NOT_PACKED;
        }
        case "3": {
            return AnOrderFulfillmentStatusEnum.CANCEL;
        }
        default: {
            return undefined;
        }
    }
}