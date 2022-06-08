import {USER_SET} from "@/store/actions/user";
import {AUTH_LOGOUT} from "@/store/actions/auth";

const profile = {
    id: "",
    email: "",
    name: "",
    role: 0,
    roleString: "",
    contact: {
        id: 0,
        userId: "",
        name: "",
        phoneNumber: "",
        address: "",
        district: "",
        state: "",
        province: "",
        postcode: "",
    }
}
const state = {
    profile: localStorage.getItem('profile') !== null ? JSON.parse(localStorage.getItem('profile')): profile,
};

const mutations = {
    [USER_SET]: (state, payload) => {
        state.profile = payload
        localStorage.setItem('profile', JSON.stringify(payload))
    },
    [AUTH_LOGOUT]: state => {
        state.profile = {};
        localStorage.removeItem('profile')
    }
};

export default {
    state,
    mutations
};
