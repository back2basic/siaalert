import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function convertSectorsToBytes(sectors: number) {
  return sectors * 4 * 1024 * 1024;
}

export function formatStorage(bytes: number) {
  const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"]; // Units from bytes to yottabytes
  let i = 0;

  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024;
    i++;
  }

  return `${bytes.toFixed(2)} ${units[i]}`; // Two decimal places for readability
}

// https://github.com/mike76-dev/hostscore/blob/master/web/src/api/helpers.ts
export const convertPrice = (value: string) => {
  if (value.length < 13) return value + " H";
  if (value.length < 16) return value.slice(0, value.length - 12) + " pS";
  if (value.length < 19) return value.slice(0, value.length - 15) + " nS";
  if (value.length < 22) return value.slice(0, value.length - 18) + " uS";
  if (value.length < 25) return value.slice(0, value.length - 21) + " mS";
  if (value.length < 28) {
    let result = value.slice(0, value.length - 24);
    if (value[value.length - 24] !== "0")
      result += "." + value[value.length - 24];
    return result + " SC";
  }
  if (value.length < 31) {
    let result = value.slice(0, value.length - 27);
    if (value[value.length - 27] !== "0")
      result += "." + value[value.length - 27];
    return result + " KS";
  }
  let result = value.slice(0, value.length - 30);
  if (value[value.length - 30] !== "0")
    result += "." + value[value.length - 30];
  return result + " MS";
};

export const convertPricePerBlock = (value: string) => {
  if (value.length > 12) {
    value =
      value.slice(0, value.length - 12) + "." + value.slice(value.length - 12);
  } else {
    value = "0." + "0".repeat(12 - value.length) + value;
  }
  const result = parseFloat(value) * 144 * 30;
  if (result < 1e-12) return (result * 1e24).toFixed(0) + " H/TB/month";
  if (result < 1e-9) return (result * 1e24).toFixed(0) + " pS/TB/month";
  if (result < 1e-6) return (result * 1e24).toFixed(0) + " uS/TB/month";
  if (result < 1e-3) return (result * 1e24).toFixed(0) + " mS/TB/month";
  if (result < 1) return result.toFixed(1) + " SC/TB/month";
  if (result < 1e3) return result.toFixed(0) + " SC/TB/month";
  if (result < 1e6) return (result / 1e3).toFixed(1) + " KS/TB/month";
  return (result / 1e6).toFixed(1) + " MS/TB/month";
};

export const convertPriceRaw = (value: string) => {
  if (value.length > 12) {
    value =
      value.slice(0, value.length - 12) + "." + value.slice(value.length - 12);
  } else {
    value = "0." + "0".repeat(12 - value.length) + value;
  }
  return Number.parseFloat(value);
};
export const toSia = (value: string, perBlock: boolean) => {
  let price = convertPriceRaw(value);
  if (perBlock) price *= 144 * 30;
  if (price < 1e-12) return "0 H";
  if (price < 1e-9) return (price * 1000).toFixed(0) + " pS";
  if (price < 1e-6) return (price * 1000).toFixed(0) + " nS";
  if (price < 1e-3) return (price * 1000).toFixed(0) + " uS";
  if (price < 1) return (price * 1000).toFixed(0) + " mS";
  if (price < 10) return price.toFixed(1) + " SC";
  if (price < 1e3) return price.toFixed(0) + " SC";
  if (price < 1e4) return (price / 1000).toFixed(1) + " KS";
  if (price < 1e6) return (price / 1000).toFixed(0) + " KS";
  if (price < 1e7) return (price / 1e6).toFixed(1) + " MS";
  if (price < 1e9) return (price / 1e6).toFixed(0) + " MS";
  if (price < 1e10) return (price / 1e9).toFixed(1) + " GS";
  return (price / 1e9).toFixed(0) + " GS";
};
