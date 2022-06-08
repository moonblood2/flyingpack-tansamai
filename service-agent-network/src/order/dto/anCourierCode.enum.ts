export enum AnCourierCodeEnum {
    FLASH = 1,
    KERRY = 2,
    SCG = 3,
    EMS = 4,
    EMS_WORLD = 5,
    DHL_WORLD = 6,
    MESSENGER = 101,
}

export const AnCourierCodeToString = {
    [AnCourierCodeEnum.FLASH]: "Flash Express",
    [AnCourierCodeEnum.KERRY]: "Kerry Express",
    [AnCourierCodeEnum.SCG]: "SCG",
    [AnCourierCodeEnum.EMS]: "EMS",
    [AnCourierCodeEnum.EMS_WORLD]: "EMS World",
    [AnCourierCodeEnum.DHL_WORLD]: "DHL World",
    [AnCourierCodeEnum.MESSENGER]: "Messenger",
}

export const anCourierCodeToSpCourierCode = {
    [AnCourierCodeEnum.FLASH]: "FLE",
    [AnCourierCodeEnum.KERRY]: "KRYX",//Use KRYX for temporary in the promotion season, will back to KRYP on 2021/08/31
    [AnCourierCodeEnum.SCG]: "SCG",
    [AnCourierCodeEnum.EMS]: "EMST",
    [AnCourierCodeEnum.EMS_WORLD]: "",
    [AnCourierCodeEnum.DHL_WORLD]: "",
    [AnCourierCodeEnum.MESSENGER]: "",
}