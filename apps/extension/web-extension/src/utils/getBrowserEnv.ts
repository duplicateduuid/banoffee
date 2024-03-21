export const getBrowserEnv = (): typeof chrome => {
    if (typeof chrome !== "undefined") {
        return chrome;
    } else if (typeof browser !== "undefined") {
        return browser;
    }

    throw new Error("Not supported browser");
}