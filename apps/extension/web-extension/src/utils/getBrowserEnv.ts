export const getBrowserEnv = (): typeof chrome | undefined => {
    if (typeof chrome !== "undefined") {
        return chrome;
    } else if (typeof browser !== "undefined") {
        return browser;
    }
}