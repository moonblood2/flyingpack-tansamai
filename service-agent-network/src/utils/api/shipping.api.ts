import axios from "axios";

import serviceConfig from "../../config/service.config";
import apiKeyConfig from "../../config/apikey.config";
import debug from "debug";

const log: debug.IDebugger = debug("app:utils:api");

export const getOrderParcelByIds = (
  token: string,
  userId: string,
  ids: Array<string>
) => {
  return new Promise((resolve, reject) => {
    axios({
      url: `${serviceConfig.SERVICE_URL_SHIPPING}/orders/parcels/ids/`,
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
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
  });
};

export const makeOrder = (data:any) => {
//   console.log("bodyObject xx : ");
//   console.log(JSON.stringify({ data: anOrders }));
//   console.log("URL : " + `${serviceConfig.SERVICE_URL_AURABLUE}`);

    return new Promise((resolve, reject) => {
      axios({
        url: `${serviceConfig.SERVICE_URL_AURABLUE}`,
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          "api-key":`${apiKeyConfig.API_KEY_AURABULE}`,
        },
        data:data
      })
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
            console.log(error);
          reject(error);
        });
    });
};
