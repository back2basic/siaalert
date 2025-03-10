import { create } from "zustand";

interface VersionState {
  version: string;
  setVersion: (version: string) => void;
}

interface StorageState {
  v1total: number;
  v2total: number;
  v1used: number;
  v2used: number;
  v1hosts: number;
  v2hosts: number;

  setV1Total: (v1total: number) => void;
  setV2Total: (v2total: number) => void;
  setV1Used: (v1used: number) => void;
  setV2used: (v2used: number) => void;
  setV1Hosts: (v1hosts: number) => void;
  setV2Hosts: (v2hosts: number) => void;
}

export const UseVersionStore = create<VersionState>()((set) => ({
  version: "",
  setVersion: (version: string) => set({ version }),
}));

export const UseStorageStore = create<StorageState>()((set) => ({
  v1total: 0,
  v2total: 0,
  v1used: 0,
  v2used: 0,
  v1hosts: 0,
  v2hosts: 0,
  setV1Total: (v1total: number) => set({ v1total }),
  setV2Total: (v2total: number) => set({ v2total }),
  setV1Used: (v1used: number) => set({ v1used }),
  setV2used: (v2used: number) => set({ v2used }),
  setV1Hosts: (v1hosts: number) => set({ v1hosts }),
  setV2Hosts: (v2hosts: number) => set({ v2hosts }),
}));
