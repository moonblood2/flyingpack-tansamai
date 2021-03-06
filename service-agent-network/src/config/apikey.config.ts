if (process.env.API_KEY_LENGTH === undefined) {
    throw new Error("env: API_KEY_LENGTH undefined");
}
if (process.env.API_KEY_PREFIX_LENGTH === undefined) {
    throw new Error("env: API_KEY_PREFIX_LENGTH undefined");
}
if (process.env.API_KEY_AURABULE === undefined) {
    throw new Error("env: API_KEY_AURABULE undefined");
}

export default {
    API_KEY_LENGTH: parseInt(process.env.API_KEY_LENGTH),
    API_KEY_PREFIX_LENGTH: parseInt(process.env.API_KEY_PREFIX_LENGTH),
    API_KEY_AURABULE: process.env.API_KEY_AURABULE
}