const getByID = (id: string, dflt: any): any => document.getElementById(id) ? document.getElementById(id) : dflt;

export default getByID
