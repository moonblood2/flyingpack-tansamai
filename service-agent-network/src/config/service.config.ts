if (process.env.JWT_SIGNING_KEY === undefined) {
    throw new Error("env: JWT_SIGNING_KEY undefined");
}
if (process.env.SERVICE_URL_SHIPPING === undefined) {
    throw new Error("env: SERVICE_URL_SHIPPING undefined");
}

export default {
    JWT_SIGNING_KEY: process.env.JWT_SIGNING_KEY,
    SERVICE_URL_SHIPPING: process.env.SERVICE_URL_SHIPPING,
}