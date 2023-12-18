export const getClass = (classSuffix: string) => `content-blocks__${classSuffix}`;
export const getClasses = (classSuffixes: string[]) => classSuffixes.map(getClass).join(' ');

export const getPageClass = (classSuffix: string) => `content-page__${classSuffix}`;
export const getPageClasses = (classSuffixes: string[]) => classSuffixes.map(getPageClass).join(' ');
