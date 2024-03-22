// TODO: test it on firefox
document.addEventListener("loginAttempt", (event) => {
    if (typeof chrome !== "undefined") {
        chrome.storage.local.set({ 'session': { id: event.detail.sessionId, expiration: event.detail.expiration } }).then(() => {});
    } else if (typeof browser !== "undefined") {
        browser.storage.local.set({ 'session': { id: event.detail.sessionId, expiration: event.detail.expiration } }).then(() => {});
    }
});