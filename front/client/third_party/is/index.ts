export const undef = (v) => v === null || v === undefined
export const notUndef = (v) => v !== null && v !== undefined
export const notEmpty = (v) => notUndef(v) && v !== "" && v !== 0;
export const empty = (v) => undef(v) || v == "" || v == 0;
export const object = (v) => v && !array(v) && typeof v === 'object'
export const trueVal = (v) => notUndef(v) && v !== 0 && v !== ""
export const func = (v) => typeof v === "function"
export const string = (v) => typeof v === "string"
export const number = (v) => typeof n === "number"