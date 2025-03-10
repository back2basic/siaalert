import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatBytes(v2: boolean, bytes: number) {
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
