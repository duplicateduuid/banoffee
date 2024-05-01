export const getCookies = (cookieString: string): Record<string, unknown> => {
	if (!cookieString || cookieString === '') return {};
	return cookieString
		.split(';')
		.map((x) => x.trim().split(/(=)/))
		.reduce(
			(cookiesObject, currentArray) => ({
				...cookiesObject,
				[currentArray[0]]: decodeURIComponent(currentArray[2])
			}),
			{}
		);
};

export const isValidUrl = (urlString: string) => {
	try {
		return Boolean(new URL(urlString));
	} catch (e) {
		return false;
	}
};

export const removeQueryParams = (url: string): string => {
	const urlObject = new URL(url);
	urlObject.search = '';
	return urlObject.toString();
};

export const deleteCookie = (name: string) => {
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`;
};