export const AnCourier = {
    FLASH: {
        text: "Flash Express",
        code: 1,
    },
    KERRY: {
        text: "Kerry Express",
        code: 2,
    },
    SCG: {
        text: "SCG Yamato Express",
        code: 3,
    },
    EMS: {
        text: "EMS",
        code: 4,
    },
    EMS_WORLD: {
        text: "EMS World",
        code: 5,
    },
    DHL_WORLD: {
        text: "DHL World",
        code: 6,
    },
    MESSENGER: {
        text: "Messenger",
        code: 101,
    },
}

let _courierCodeToCourierName = {};
for (const k in AnCourier) {
    let e = AnCourier[k];
    _courierCodeToCourierName[e.code] = e.text;
}

export const CourierCodeToCourierName = (code) => {
    if (code in _courierCodeToCourierName) {
        return _courierCodeToCourierName[code];
    } else {
        return "No courier_code";
    }
};