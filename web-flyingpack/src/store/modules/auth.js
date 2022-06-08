/* eslint-disable promise/param-names */
import axios from 'axios'
import {AUTH_ERROR, AUTH_LOGOUT, AUTH_REQUEST, AUTH_SUCCESS} from "../actions/auth";
import env from '@/constants/env'
import {USER_SET} from "@/store/actions/user";

const state = {
    status: "",
    error: {
        details: null
    },
    token: localStorage.getItem("user-token") || "",
};

const getters = {
    authStatus: state => state.status,
    authLoading: state => state.status === "loading",
    authValidation: state => {
        return {
            valid: state.status !== "error",
            details: state.error.details
        }
    },
    isAuthenticated: state => {
        return !!state.token
    }
};

const actions = {
    [AUTH_REQUEST]: ({commit}, user) => {
        return new Promise((resolve, reject) => {
            commit(AUTH_REQUEST);
            axios({
                url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/auth/login/`,
                method: "POST",
                data: {
                    email: user.email,
                    password: user.password
                }
            })
                .then(response => {
                    const {data} = response;
                    if (data) {
                        let profile = {
                            id: data.id,
                            email: data.email,
                            name: data.name,
                            role: data.role,
                            roleString: data.roleString,
                        };
                        if (data.contact) {
                            profile.contact = {
                                id: data.contact['id'],
                                userId: data.contact['user_id'],
                                name: data.contact['name'],
                                phoneNumber: data.contact['phone_number'],
                                address: data.contact['address'],
                                district: data.contact['district'],
                                state: data.contact['state'],
                                province: data.contact['province'],
                                postcode: data.contact['postcode'],
                            }
                        }

                        const token = `Bearer ${data.token}`;
                        localStorage.setItem("user-token", token);

                        commit(AUTH_SUCCESS, token);
                        commit(USER_SET, profile);

                        resolve(response);
                    } else {
                        reject(response);
                    }
                })
                .catch(error => {
                    if (error.response) {
                        commit(AUTH_ERROR, error.response.data);
                    } else {
                        commit(AUTH_ERROR, {details: "ไม่สำเร็จ"});
                    }
                    localStorage.removeItem("user-token");
                    reject(error);
                });
        });
    },
    [AUTH_LOGOUT]: ({commit}) => {
        return new Promise(resolve => {
            commit(AUTH_LOGOUT);
            localStorage.removeItem("user-token");
            resolve();
        });
    }
};

const mutations = {
    [AUTH_REQUEST]: state => {
        state.status = "loading";
    },
    [AUTH_SUCCESS]: (state, token) => {
        state.status = "success";
        state.token = token
    },
    [AUTH_ERROR]: (state, error) => {
        state.status = "error";
        state.error = error;
    },
    [AUTH_LOGOUT]: state => {
        state.token = "";
    }
};

export default {
    state,
    getters,
    actions,
    mutations
};
