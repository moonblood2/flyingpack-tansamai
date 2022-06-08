import cryptoRandomString from "crypto-random-string";
import jsSHA from "jssha";

export interface Key {
    apiKey: string;
    prefixApiKey?: string;
    hash: string;
}

class CommonApikey {
    private static instance: CommonApikey;

    static getInstance(): CommonApikey {
        if (!CommonApikey.instance) {
            CommonApikey.instance = new CommonApikey();
        }
        return CommonApikey.instance;
    }

    //Generate new api key, return api key, prefix and hash of whole api key.
    //Use crypto random for generate new key, and SHA3-512 for hashing.
    generateKey(apiKeyLength: number, prefixLength: number): Key {
        const apiKey = cryptoRandomString({length: apiKeyLength, type: 'alphanumeric'});
        const prefix = this.retrieveApiKeyPrefix(apiKey, prefixLength);
        const shaObj = new jsSHA("SHA3-512", "TEXT", { encoding: "UTF8" });
        shaObj.update(apiKey);
        //SHA3-512 have: 512 bits = 64 Bytes, 128 HEX length
        const hash = shaObj.getHash("HEX");

        return {
            apiKey: apiKey,
            prefixApiKey: prefix,
            hash: hash,
        }
    }

    retrieveApiKeyPrefix(apiKey: string, prefixLength: number): string {
        return apiKey.substr(0, prefixLength);
    }

    verifyKey(key: Key): boolean {
        const shaObj = new jsSHA("SHA3-512", "TEXT", { encoding: "UTF8" });
        shaObj.update(key.apiKey);
        const hash = shaObj.getHash("HEX");
        return key.hash === hash;
    }
}

export default CommonApikey.getInstance();