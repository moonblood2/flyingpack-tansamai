if (process.env.JWT_SIGNING_KEY === undefined) {
    throw new Error("env: JWT_SIGNING_KEY undefined");
}
if (process.env.SERVICE_URL_SHIPPING === undefined) {
    throw new Error("env: SERVICE_URL_SHIPPING undefined");
}
if (process.env.SERVICE_URL_AURABLUE === undefined) {
    throw new Error("env: SERVICE_URL_AURABLUE undefined");
}


export default {
    JWT_SIGNING_KEY: process.env.JWT_SIGNING_KEY,
    SERVICE_URL_SHIPPING: process.env.SERVICE_URL_SHIPPING,
    SERVICE_URL_AURABLUE:process.env.SERVICE_URL_AURABLUE
}