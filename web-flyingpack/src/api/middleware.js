import axios from "axios";
import store from "@/store";
import {AUTH_LOGOUT} from "@/store/actions/auth";
import router from "@/router";

//Middleware, handle 401 and 403
axios.interceptors.response.use((response) => {
    return response;
}, (async error => {
    switch (error.response.status) {
        case 401:
            store.dispatch(AUTH_LOGOUT).then(async () => await router.push("/login"))
            break
        case 403:
            await router.push("/")
    }
    return Promise.reject(error);
}))