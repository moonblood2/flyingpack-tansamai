export enum DateTypeEnum {
    CREATED_DATE = 1,
    STATUS_COMPLETED_DATE = 2,
    COD_TRANSFERRED_DATE = 3,
}

export const parseDateTypeEnum = (t: string) => {
    switch (t) {
        case "1":
            return DateTypeEnum.CREATED_DATE
        case "2":
            return DateTypeEnum.STATUS_COMPLETED_DATE
        case "3":
            return DateTypeEnum.COD_TRANSFERRED_DATE
        default:
            return DateTypeEnum.CREATED_DATE
    }
}