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
