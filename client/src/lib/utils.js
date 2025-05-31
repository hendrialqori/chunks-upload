export const sleep = (time) => new Promise((r) => setTimeout(r, time))

export const toMB = (size) => (size / (1 * 1024 * 1024)).toFixed(2)