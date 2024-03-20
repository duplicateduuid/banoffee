// TODO: make it work for firefox
document.addEventListener("loginAttempt", (event) => {
    if (typeof chrome !== "undefined") {
        chrome.storage.local.set({ 'sessionId': event.detail.sessionId }).then(() => {});
    }
});