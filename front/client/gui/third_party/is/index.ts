export const undef = (v) => v === null || v === undefined
export const notUndef = (v) => v !== null && v !== undefined
export const notEmpty = (v) => notUndef(v) && v !== "" && v !== 0;
export const object = (obj) => obj && !array(obj) && typeof obj === 'object'
export const trueVal = (obj) => notUndef(obj) && val !== 0 && val !== ""