import {SLIP_SET} from "@/store/actions/slip";
import {AUTH_LOGOUT} from "@/store/actions/auth";

const slip = {
    hasSet: false,
    userId: "",
    message: "",
    createdAt: "",
    deletedAt: "",
}

const state = {
    slip: localStorage.getItem('slip') !== null ? JSON.parse(localStorage.getItem('slip')): slip,
};

const mutations = {
    [SLIP_SET]: (state, payload) => {
        state.slip = {hasSet: true, ...payload}
        localStorage.setItem('slip', JSON.stringify(state.slip))
    },
    [AUTH_LOGOUT]: state => {
        state.slip = {}
        localStorage.removeItem('slip')
    }
};

export default {
    state,
    mutations
};
