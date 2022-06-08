import {PACKING_SLIP_ADD_ITEM} from "@/store/actions/packing-slip";
import {AUTH_LOGOUT} from "@/store/actions/auth";

const storageKey = 'packing-slip-item';
const state = {
    items: JSON.parse(localStorage.getItem(storageKey)),
}

const mutations = {
    [PACKING_SLIP_ADD_ITEM]: (state, payload) => {
        const {key, items} = payload;
        const itemOnLocalStorage = JSON.parse(localStorage.getItem(storageKey))
        const itemObject = ({...itemOnLocalStorage, [key]: items})
        state.items = items
        localStorage.setItem(storageKey, JSON.stringify(itemObject))
    },
    [AUTH_LOGOUT]: () => {
        state.parcel = {}
        localStorage.removeItem(storageKey)
    }
}

export default {
    state,
    mutations,
}