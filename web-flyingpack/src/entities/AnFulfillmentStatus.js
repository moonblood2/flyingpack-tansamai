export const anFulfillmentStatus = {
    NOT_SET: 0,
    PACKED: 1,
    NOT_PACKED: 2,
    CANCEL: 3,
}

export const anFulfillmentStatusToString = {
    [anFulfillmentStatus.NOT_SET]: "-",
    [anFulfillmentStatus.PACKED]: "จัดส่งแล้ว",
    [anFulfillmentStatus.NOT_PACKED]: "ยังไม่จัดส่ง",
    [anFulfillmentStatus.CANCEL]: "ยกเลิก",
}
