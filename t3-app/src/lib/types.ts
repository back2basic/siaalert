import type { Models } from "appwrite";

export type Host = {
  countryCode: string;
  failedInteractions: number;
  knownSince: string;
  lastAnnouncement: string;
  lastScan: string;
  lastScanSuccessful: boolean;
  netAddress: string;
  publicKey: string;
  successfulInteractions: number;
  totalScans: number;
  v2: boolean;
  v2NetAddresses: string;
  v2NetAddressesProto: string;
  error: string;
  online: boolean;
  onlineSince: string;
  offlineSince: string;
};
export type HostDoc = Models.Document & Host;
export type HostList = Models.DocumentList<HostDoc>;

export type Status = {
  height: number;
};

export type StatusDoc = Models.Document & Status;
export type StatusList = Models.DocumentList<StatusDoc>;

export type Network = "mainnet" | "zen";

export type Check = {
  hostId: string;
  v4Addr: string;
  v6Addr: string;
  rhp2Port: string;
  rhp2V4Delay: number;
  rhp2V6Delay: number;
  rhp2V4: boolean;
  rhp2V6: boolean;
  rhp3Port: string;
  rhp3V4: boolean;
  rhp3V6: boolean;
  rhp3V4Delay: number;
  rhp3V6Delay: number;
  rhp4Port: string;
  rhp4V4: boolean;
  rhp4V6: boolean;
  rhp4V4Delay: number;
  rhp4V6Delay: number;
  acceptingContracts: boolean;
  release: string;
};

export type CheckDoc = Models.Document & Check;
export type CheckList = Models.DocumentList<CheckDoc>;

export type GitHiubHostdRelease = {
  name: string;
  tag_name: string;
  html_url: string;
  body: string;
};
