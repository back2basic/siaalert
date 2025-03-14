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
  createdAt: string;
  publicKey: string;
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

export type HostScan = {
  publicKey: string;
  v2: boolean;
  v2NetAddresses: V2NetAddress[];
  netAddress: string;
  success: boolean;
  timestamp: string;
  nextScan: string;
  acceptingContracts: boolean;
  error: string;
  onlineSince: string;
  offlineSince: string;
  totalStorage: number;
  remainingStorage: number;
  settings: Rhpv2Settings;
  priceTable: Rhpv3Settings;
  rhpV4Settings: Rhpv4Settings;
};

export type Rhpv2Settings = {
  maxdownloadbatchsize: number;
  acceptingcontracts: boolean;
  baserpcprice: string;
  collateral: string;
  contractprice: string;
  downloadbandwidthprice: string;
  ephemeralaccountexpiry: number;
  maxcollateral: string;
  maxduration: number;
  maxephemeralaccountbalance: string;
  maxrevisebatchsize: number;
  netaddress: string;
  release: string;
  remainingstorage: number;
  revisionnumber: number;
  sectoraccessprice: string;
  sectorsize: number;
  siamuxport: string;
  storageprice: string;
  totalstorage: number;
  unlockhash: string;
  uploadbandwidthprice: string;
  version: string;
  windowsize: number;
};

export type Rhpv3Settings = {
  uid: string;
  validity: number;
  hostblockhieight: number;
  updatepricetablecost: number;
  accountbalancecost: number;
  fundaccountcost: number;
  latestrevisioncost: number;
  subscriptionmemorycost: number;
  subscriptionnotificationcost: number;
  initbasecost: number;
  memorytimecost: number;
  downloadbandwidthcost: string;
  uploadbandwidthcost: string;
  dropsectorsbasecost: number;
  dropsectorsunitcost: number;
  hassectorbasecost: number;
  readbasecost: number;
  readlengthcost: number;
  renewcontractcost: number;
  revisionbasecost: number;
  swapsectorcost: number;
  writebasecost: number;
  writelengthcost: number;
  writestorecost: number;
  txnfeeminrecommended: number;
  txnfeemaxrecommended: number;
  contractprice: number;
  collateralcost: number;
  maxcollateral: string;
  maxduration: number;
  maxephemeralaccountbalance: number;
  windowsize: number;
  registryentriesleft: number;
  registryentriestotal: number;
};

export type Rhpv4Settings = {
  protocolVersion: number[];
  release: string;
  walletAddress: string;
  acceptingContracts: boolean;
  maxCollateral: string;
  maxContractDuration: number;
  remainingStorage: number;
  totalStorage: number;
  prices: V2Prices;
};

export type V2NetAddress = {
  protocol: string;
  address: string;
};

export type V2Prices = {
  contractPrice: number;
  collateral: number;
  storagePrice: string;
  ingressPrice: string;
  egressPrice: string;
  freeSectorPrice: number;
  tipHeight: number;
  validUntil: Date;
  signature: string;
};
