import axios from "axios";
import env from "@/constants/env";
import "./middleware";

// objectToRoboticString, parse object to Robotic's string format eg. {\"SKU4\": 1, \"SKU1\": 3}
// const objectToRoboticString = (data) => {
//     let pair = [];//Key and value in format like this: \"key\": value
//     for (const key in data) {
//         pair.push(`\\"${key}\\": ${data[key]}`);
//     }
//     return `{${pair.join(", ")}}`;
// }

export const takePhoto = (referenceNo) => {
    return new Promise((resolve, reject) => {
        axios({
            responseType: "blob",
            url: `${env.VUE_APP_SERVICE_ROBOT}/packing_img?packing_no=${referenceNo}`,
            method: "GET",
            headers: {"Content-Type": "application/json"},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    });
}

// addOrder, tell robotic to start packing.
export const addOrder = (referenceNo, order) => {
    return new Promise((resolve, reject) => {
        axios({
            url: `${env.VUE_APP_SERVICE_ROBOT}/add_data`,
            method: "POST",
            headers: {"Content-Type": "application/json"},
            data: {
                'packing_no': referenceNo,
                'packing_detail': order,
            }
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    })
}

export const getPackingData = (createdAt, packingStatus) => {
    let url = `${env.VUE_APP_SERVICE_ROBOT}/get_packing_data?created_at=${createdAt}`;
    if (packingStatus) {
        url += `&packing_status=${packingStatus}`
    }

    return new Promise((resolve, reject) => {
        axios({
            url: url,
            method: "GET",
            headers: {"Content-Type": "application/json"},
        })
            .then(response => {
                resolve(response);
            })
            .catch(error => {
                reject(error);
            });
    })
}
