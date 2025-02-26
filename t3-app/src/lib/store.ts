import { create } from "zustand";

interface VersionState {
  version: string;
  setVersion: (version: string) => void;
}

export const UseVersionStore = create<VersionState>()((set) => ({
  version: "",
  setVersion: (version: string) => set({ version }),
}));
