import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatStorage(v2: boolean, bytes: number) {
  if (v2) {
    bytes = bytes * 4 * 1024 * 1024;
    // return;
  }
  const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"]; // Units from bytes to yottabytes
  let i = 0;

  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024;
    i++;
  }

  return `${bytes.toFixed(2)} ${units[i]}`; // Two decimal places for readability
}

// https://github.com/mike76-dev/hostscore/blob/master/web/src/api/helpers.ts#L81
export const convertPrice = (value: string) => {
	if (value.length < 13) return value + ' H'
	if (value.length < 16) return value.slice(0, value.length - 12) + ' pS'
	if (value.length < 19) return value.slice(0, value.length - 15) + ' nS'
	if (value.length < 22) return value.slice(0, value.length - 18) + ' uS'
	if (value.length < 25) return value.slice(0, value.length - 21) + ' mS'
	if (value.length < 28) {
		let result = value.slice(0, value.length - 24)
		if (value[value.length-24] !== '0') result += '.' + value[value.length-24]
		return result + ' SC'
	}
	if (value.length < 31) {
		let result = value.slice(0, value.length - 27)
		if (value[value.length-27] !== '0') result += '.' + value[value.length-27]
		return result + ' KS'
	}
	let result = value.slice(0, value.length - 30)
	if (value[value.length-30] !== '0') result += '.' + value[value.length-30]
	return result + ' MS'
}