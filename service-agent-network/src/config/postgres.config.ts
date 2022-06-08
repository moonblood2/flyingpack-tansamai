if (process.env.POSTGRES_USER === undefined) {
    throw new Error("env: POSTGRES_USER undefined");
}
if (process.env.POSTGRES_HOST === undefined) {
    throw new Error("env: POSTGRES_HOST undefined");
}
if (process.env.POSTGRES_DATABASE === undefined) {
    throw new Error("env: POSTGRES_DATABASE undefined");
}
if (process.env.POSTGRES_PASSWORD === undefined) {
    throw new Error("env: POSTGRES_PASSWORD undefined");
}
if (process.env.POSTGRES_PORT === undefined) {
    throw new Error("env: POSTGRES_PORT undefined");
}
if (process.env.POSTGRES_URI === undefined) {
    throw new Error("env: POSTGRES_URI undefined");
}
if (process.env.POSTGRES_ENABLE_SSL === undefined) {
    throw new Error("env: POSTGRES_ENABLE_SSL undefined");
}

export default {
    POSTGRES_USER: process.env.POSTGRES_USER,
    POSTGRES_HOST: process.env.POSTGRES_HOST,
    POSTGRES_DATABASE: process.env.POSTGRES_DATABASE,
    POSTGRES_PASSWORD: process.env.POSTGRES_PASSWORD,
    POSTGRES_PORT: parseInt(process.env.POSTGRES_PORT),
    POSTGRES_URI: process.env.POSTGRES_URI,
    POSTGRES_ENABLE_SSL: process.env.POSTGRES_ENABLE_SSL === "true",
}