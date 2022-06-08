import axios from "axios";
import env from "@/constants/env";
import store from "@/store";
import "./middleware";

const getUsers = () => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/users/`,
            method: "GET",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response)
            })
            .catch(error => {
                reject(error)
            });
    })
}

// getParcelPrice get parcel price from ShippingService.
const getParcelPrice = (parcels) => {
    if (!Array.isArray(parcels)) {
        parcels = [parcels]
    }
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/parcel-price/`,
            method: "POST",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
            data: parcels,
        })
            .then(response => {
                resolve(response)
            })
            .catch(error => {
                reject(error)
            });
    })
}

// createOrder send created order request to ShippingService.
const createOrder = ({type = 1, sender, parcels = [], anParcels = [], products = [], paymentMethod = 1}) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/orders/`,
            method: "POST",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
            data: {
                type: type,
                sender: sender,
                parcels: parcels,
                an_parcels: anParcels,
                products: products,
                payment_method: paymentMethod,
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

// getOrderProduct get order parcel of each Shop from ShippingService.
const getOrderParcel = (startDate, endDate) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/orders/parcels/start/${startDate}/end/${endDate}/`,
            method: "GET",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

//updateOrderParcel
const updateOrderParcel = (type, orderParcelId, trackingCode) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/orders/parcels/`,
            method: "PUT",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
            data: {
                parcel_type: type,
                order_parcel_id: orderParcelId,
                tracking_code: trackingCode,
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

// getOrderProduct get order product of each Shop from ShippingService.
const getOrderProduct = (startDate, endDate) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_SHIPPING_URL}/orders/products/start/${startDate}/end/${endDate}/`,
            method: "GET",
            headers: {'Content-Type': "application/json", 'Authorization': store.state.auth.token},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

export {
    getUsers,
    getParcelPrice,
    createOrder,
    getOrderParcel,
    updateOrderParcel,
    getOrderProduct,
}