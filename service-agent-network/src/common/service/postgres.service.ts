import type, {Client} from "pg";
import postgresConfig from "../../config/postgres.config";

import debug from "debug";

const log: debug.IDebugger = debug("app:postgres-service")

class PostgresService {
    private static instance: PostgresService;
    private client: Client = new Client({
        // user: postgresConfig.POSTGRES_USER,
        // host: postgresConfig.POSTGRES_HOST,
        // database: postgresConfig.POSTGRES_DATABASE,
        // password: postgresConfig.POSTGRES_PASSWORD,
        // port: postgresConfig.POSTGRES_PORT,
        connectionString: postgresConfig.POSTGRES_URI,
        ssl: postgresConfig.POSTGRES_ENABLE_SSL ? {rejectUnauthorized: false} : false,
    });

    private constructor() {
        type.types.setTypeParser(1114, str => str);
        (async () => {
            await this.connect();
        })();
    }

    static getInstance(): PostgresService {
        if (!PostgresService.instance) {
            PostgresService.instance = new PostgresService();
        }
        return PostgresService.instance;
    }

    getClient() {
        return this.client;
    }

    private async connect() {
        try {
            log("trying to connect...")
            await this.client.connect();
            log("connected")
        } catch (e) {
            log("disconnected")
            log(e);
            await this.client.end();
            throw new Error(e);
        }
    }
}

export default PostgresService.getInstance();