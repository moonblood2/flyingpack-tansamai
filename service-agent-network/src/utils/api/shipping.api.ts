import axios from "axios";

import serviceConfig from "../../config/service.config";

import debug from "debug";

const log: debug.IDebugger = debug("app:utils:api");

export const getOrderParcelByIds = (token: string, userId: string, ids: Array<string>) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${serviceConfig.SERVICE_URL_SHIPPING}/orders/parcels/ids/`,
            method: "POST",
            headers: {"Authorization": `Bearer ${token}`},
            data: {
                user_id: userId,
                ids: ids,
            },
        })
            .then((response) => {
                resolve(response);
            })
            .catch((error) => {
                reject(error);
            });
    })
}

export const makeOrder = (token: string, orders: Array<any>) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${serviceConfig.SERVICE_URL_SHIPPING}/orders/`,
            method: "POST",
            headers: {"Authorization": `Bearer ${token}`},
            data: {
                an_parcels: orders,
            },
        })
            .then((response) => {
                resolve(response);
            })
            .catch((error) => {
                reject(error);
            });
    })
}