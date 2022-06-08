import {LABEL_ADD_PARCEL} from "@/store/actions/label";
import {AUTH_LOGOUT} from "@/store/actions/auth";

const state = {
    parcel: JSON.parse(localStorage.getItem('parcel')),
}

const mutations = {
    [LABEL_ADD_PARCEL]: (state, payload) => {
        const {key, parcels} = payload;
        const parcelOnLocalStorage = JSON.parse(localStorage.getItem('parcel'))
        const parcelObject = ({...parcelOnLocalStorage, [key]: parcels})
        state.parcel = parcelObject
        localStorage.setItem('parcel', JSON.stringify(parcelObject))
    },
    [AUTH_LOGOUT]: () => {
        state.parcel = {}
        localStorage.removeItem('parcel')
    }
}

export default {
    state,
    mutations,
}