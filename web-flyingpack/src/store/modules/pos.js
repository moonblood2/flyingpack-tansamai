import {
    POS_SET_AN_ORDER_ID,
    POS_SET_COD_AMOUNT,
    POS_SET_COURIER_CODE,
    POS_SET_DESTINATION,
    POS_SET_TYPE,
    POS_SET_AN_ORDER_ITEMS,
    POS_SET_AN_PARCEL,
} from "@/store/actions/pos";
import {AnParcel, ContactInfo} from "@/entities";


const state = {
    type: 1,
    destination: new ContactInfo({}),
    courierCode: "",
    codAmount: 0,
    anOrderId: "",
    anOrderItems: [],
    anParcel: new AnParcel({}),
}

const mutations = {
    [POS_SET_TYPE]: (state, payload) => {
        state.type = payload;
    },
    [POS_SET_DESTINATION]: (state, payload) => {
        state.destination = payload;
    },
    [POS_SET_COURIER_CODE]: (state, payload) => {
        state.courierCode = payload;
    },
    [POS_SET_COD_AMOUNT]: (state, payload) => {
        state.codAmount = payload;
    },
    [POS_SET_AN_ORDER_ID]: (state, payload) => {
        state.anOrderId = payload;
    },
    [POS_SET_AN_ORDER_ITEMS]: (state, payload) => {
        state.anOrderItems = payload;
    },
    [POS_SET_AN_PARCEL]: (state, payload) => {
        state.anParcel = payload;
    },
}

export default {
    state,
    mutations,
}