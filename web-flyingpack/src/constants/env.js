//envMap, key = imported env key, value = exported key
const envMap = {
    "VUE_APP_SERVICE_SHIPPING_URL": "VUE_APP_SERVICE_SHIPPING_URL",
    "VUE_APP_SERVICE_AGENT_NETWORK_URL": "VUE_APP_SERVICE_AGENT_NETWORK_URL",
    "VUE_APP_SERVICE_ROBOT": "VUE_APP_SERVICE_ROBOT",
    "VUE_APP_JNA_NAME": "JNA_NAME",
    "VUE_APP_JNA_PHONE_NUMBER": "JNA_PHONE_NUMBER",
    "VUE_APP_JNA_ADDRESS": "JNA_ADDRESS",
    "VUE_APP_JNA_DISTRICT": "JNA_DISTRICT",
    "VUE_APP_JNA_STATE": "JNA_STATE",
    "VUE_APP_JNA_PROVINCE": "JNA_PROVINCE",
    "VUE_APP_JNA_POSTCODE": "JNA_POSTCODE",
    "VUE_APP_JNA_ID": "JNA_ID",
    "VUE_APP_JNA_BIRTH_DATE": "JNA_BIRTH_DATE",
    "VUE_APP_MODE": "MODE",
}

const exportedEnv = {};

//Check if env do not exist.
for (const importedKey in envMap) {
    if (!process.env[importedKey]) {
        throw new Error(`${importedKey} undefined`);
    }
    exportedEnv[envMap[importedKey]] = process.env[importedKey];
}

export default exportedEnv;