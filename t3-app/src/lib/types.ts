export type Network = "mainnet" | "zen";

export type Host = {
  // "_id": string,
  publicKey: string;
  failedInteractions: number;
  knownSince: string;
  lastAnnouncement: string;
  lastScan: string;
  lastScanSuccessful: false;
  netAddress: string;
  successfulInteractions: number;
  totalScans: number;
  v2: boolean;
};

export type Rhp = {
  publicKey: string;
  netAddress: string;
  v2: boolean;
  v2Adresses: V2Address[];
  success: boolean;
  timestamp: string;
  nextScan: string;
  acceptingContracts: boolean;
  error: string;
  onlineSince: string;
  offlineSince: string;
  totalStorage: number;
  remainingStorage: number;
};

type V2Address = {
  address: string;
  protocol: string;
};

export type Scan = {
  hostId: string;
  v4addr: string;
  v6addr: string;
  rhp2port: string;
  rhp2v4delay: number;
  rhp2v6delay: number;
  rhp2v4: boolean;
  rhp2v6: boolean;
  rhp3port: string;
  rhp3v4: boolean;
  rhp3v6: boolean;
  rhp3v4delay: number;
  rhp3v6delay: number;
  rhp4port: string;
  rhp4v4: boolean;
  rhp4v6: boolean;
  rhp4v4delay: number;
  rhp4v6delay: number;
  // acceptingContracts: boolean;
  // release: string;
};

// export type CheckDoc = Models.Document & Check;
// export type CheckList = Models.DocumentList<CheckDoc>;

export type Consensus = {
  height: number;
  id: string;
};

export type GitHiubHostdRelease = {
  name: string;
  tag_name: string;
  html_url: string;
  body: string;
};

export type AuthOutput = {
  message: string;
  publicKey: string;
  email: string;
  address: string;
};
