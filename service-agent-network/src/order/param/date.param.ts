import express from "express";

export interface DateParam {
    startDate: string;
    endDate: string;
    timeZone?: string;
}

export const parseDateParam = (req: express.Request): DateParam => {
    const timeZone: string = "Asia/Bangkok";
    const startDate: string = <string> req.query.startDate;
    const endDate: string = <string> req.query.endDate;

    return {
        startDate: startDate,
        endDate: endDate,
        timeZone: timeZone,
    }
}