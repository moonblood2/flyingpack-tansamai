import axios from "axios";
import env from "@/constants/env";
import store from "@/store";
import qs from "qs";
import "./middleware";

// getOrderFulfillment use by each AgentNetwork Member.
const getOrderFulfillment = (startDate, endDate, page = 1, perPage = 100, keyWord = "", fulfillmentStatus = "-1", courierCode = 0, productIds = [], quantity = "") => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            params: {
                startDate: startDate,
                endDate: endDate,
                page: page,
                perPage: perPage,
                keyWord: keyWord,
                fulfillmentStatus: fulfillmentStatus,
                courierCode: courierCode,
                productIds: productIds,
                quantity: quantity,
            },
            paramsSerializer: function (params) {
                return qs.stringify(params, {arrayFormat: 'brackets'})
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

// updateOrderFulfillment
// orders: Array of order
// anUserId: optional for called by admin.
const updateOrderFulfillment = (orders, agentUserId) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order`,
            method: "PUT",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
            data: {
                orders: orders,
                agentUserId: agentUserId,
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

const updateOrderFulfillmentSerialNumber = (anOrderId, items) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order`,
            method: "PUT",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
            data: {
                anOrderId: anOrderId,
                items: items,
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

// getOrderFulfillmentPrice use by Accounting.
const getOrderFulfillmentPrice = (startDate, endDate, page = 1, perPage = 999999, keyWord = "", fulfillmentStatus = "-1", courierCode = 0, productIds = [], userId = "") => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/price/user-id/${userId}`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            params: {
                startDate: startDate,
                endDate: endDate,
                page: page,
                perPage: perPage,
                keyWord: keyWord,
                fulfillmentStatus: fulfillmentStatus,
                courierCode: courierCode,
                productIds: productIds,
            },
            paramsSerializer: function (params) {
                return qs.stringify(params, {arrayFormat: 'brackets'})
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

const getOrderCod = (userId, startDate, endDate, dateType, keyWord) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/cod/user-id/${userId}`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            params: {
                startDate: startDate,
                endDate: endDate,
                dateType: dateType,
                keyWord: keyWord,
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

// getAnProducts use by each AgentNetwork Member.
const getAnProducts = () => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/product`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

// getAnProductsByUserId use by Accounting.
const getAnProductsByUserId = (userId) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/product/user-id/${userId}`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

//getOrderFulfillmentSummary get summary of order fulfillment of today.
const getOrderFulfillmentSummary = (startDate, endDate) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/summary`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            params: {
                startDate: startDate,
                endDate: endDate,
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

//getOrderByReferenceNo get order and product by a referenceNo.
const getOrderByReferenceNo = (referenceNo) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/reference-no`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            params: {
                referenceNo: referenceNo,
            }
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

//cancelAnOrder request for cancel order of AgentNetwork
const cancelAnOrder = (anOrderId) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/cancel/an-order-id/${anOrderId}`,
            method: "PUT",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

const createOrderByWebhook = (origin, orders) => {

    console.log("orders");
    console.log(orders);
    for (let i = 0; i < orders.length; i++) {
        orders[i].codAmount = parseFloat(orders[i].codAmount);
    }
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/shipping/order`,
            method: "POST",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            data: {
                origin: origin,
                orders: orders,
            },
        })
            .then(response => {
                
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

const createSerialNumbers = (anOrderId, items) => {
    return new Promise((resolve, reject) => {
       axios({
           url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/order-product/serial-number`,
           method: "POST",
           headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
           data: {
               anOrderId: anOrderId,
               items: items,
           },
       })
           .then(response => {
               resolve(response);
           })
           .catch(error => {
               reject(error);
           })
    });
}

const updateSerialNumbers = (anOrderId, serialNumbers) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/order/order-product/serial-number`,
            method: "PUT",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
            data: {
                anOrderId: anOrderId,
                serialNumbers: serialNumbers
            },
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            })
    });
}

const getSlip = () => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_AGENT_NETWORK_URL}/closed-api/slip`,
            method: "GET",
            headers: {"Content-Type": "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            })
    });
}

export {
    getOrderFulfillment,
    getOrderCod,
    updateOrderFulfillment,
    updateOrderFulfillmentSerialNumber,
    getOrderFulfillmentPrice,
    getAnProducts,
    getAnProductsByUserId,
    getOrderFulfillmentSummary,
    getOrderByReferenceNo,
    cancelAnOrder,
    createOrderByWebhook,
    createSerialNumbers,
    updateSerialNumbers,

    getSlip,
};
