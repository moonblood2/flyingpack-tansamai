export class AnBankAccount {
    constructor({
                    bank,
                    accountName,
                    accountNo,
                    email,
                    fiCode,
                    fiName,
                }) {
        this.bank = bank || "";
        this.accountName = accountName || "";
        this.accountNo = accountNo || "";
        this.email = email || "";
        this.fiCode = fiCode || "";
        this.fiName = fiName || "";
    }
}