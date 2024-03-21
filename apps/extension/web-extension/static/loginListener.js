// TODO: test it on firefox
document.addEventListener("loginAttempt", (event) => {
    if (typeof chrome !== "undefined") {
        chrome.storage.local.set({ 'sessionId': event.detail.sessionId }).then(() => {});
    } else if (typeof browser !== "undefined") {
        browser.storage.local.set({ 'sessionId': event.detail.sessionId }).then(() => {});
    }
});