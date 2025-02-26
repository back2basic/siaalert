package bench

type PriceTable struct {
	Uid                          string `json:"uid"`
	Validity                     int64  `json:"validity"`
	Hostblockheight              int64  `json:"hostblockheight"`
	Updatepricetablecost         string `json:"updatepricetablecost"`
	Accountbalancecost           string `json:"accountbalancecost"`
	Fundaccountcost              string `json:"fundaccountcost"`
	Latestrevisioncost           string `json:"latestrevisioncost"`
	Subscriptionmemorycost       string `json:"subscriptionmemorycost"`
	Subscriptionnotificationcost string `json:"subscriptionnotificationcost"`
	Initbasecost                 string `json:"initbasecost"`
	Memorytimecost               string `json:"memorytimecost"`
	Downloadbandwidthcost        string `json:"downloadbandwidthcost"`
	Uploadbandwidthcost          string `json:"uploadbandwidthcost"`
	Dropsectorsbasecost          string `json:"dropsectorsbasecost"`
	Dropsectorsunitcost          string `json:"dropsectorsunitcost"`
	Hassectorbasecost            string `json:"hassectorbasecost"`
	Readbasecost                 string `json:"readbasecost"`
	Readlengthcost               string `json:"readlengthcost"`
	Renewcontractcost            string `json:"renewcontractcost"`
	Revisionbasecost             string `json:"revisionbasecost"`
	Swapsectorcost               string `json:"swapsectorcost"`
	Writebasecost                string `json:"writebasecost"`
	Writelengthcost              string `json:"writelengthcost"`
	Writestorecost               string `json:"writestorecost"`
	Txnfeeminrecommended         string `json:"txnfeeminrecommended"`
	Txnfeemaxrecommended         string `json:"txnfeemaxrecommended"`
	Contractprice                string `json:"contractprice"`
	Collateralcost               string `json:"collateralcost"`
	Maxcollateral                string `json:"maxcollateral"`
	Maxduration                  int64  `json:"maxduration"`
	Windowsize                   int64  `json:"windowsize"`
	Registryentriesleft          int64  `json:"registryentriesleft"`
	Registryentriestotal         int64  `json:"registryentriestotal"`
	Expiry                       string `json:"expiry"`
}

type Settings struct {
	Acceptingcontracts         bool    `json:"acceptingcontracts"`
	Baserpcprice               string  `json:"baserpcprice"`
	Collateral                 string  `json:"collateral"`
	Contractprice              string  `json:"contractprice"`
	Downloadbandwidthprice     string  `json:"downloadbandwidthprice"`
	Ephemeralaccountexpiry     int64   `json:"ephemeralaccountexpiry"`
	Maxcollateral              string  `json:"maxcollateral"`
	Maxdownloadbatchsize       int64   `json:"maxdownloadbatchsize"`
	Maxduration                int64   `json:"maxduration"`
	Maxephemeralaccountbalance string  `json:"maxephemeralaccountbalance"`
	Maxrevisebatchsize         int64   `json:"maxrevisebatchsize"`
	Netaddress                 string  `json:"netaddress"`
	Release                    string  `json:"release"`
	Remainingstorage           float64 `json:"remainingstorage"`
	Revisionnumber             int64   `json:"revisionnumber"`
	Sectoraccessprice          string  `json:"sectoraccessprice"`
	Sectorsize                 int64   `json:"sectorsize"`
	Siamuxport                 string  `json:"siamuxport"`
	Storageprice               string  `json:"storageprice"`
	Totalstorage               float64 `json:"totalstorage"`
	Unlockhash                 string  `json:"unlockhash"`
	Uploadbandwidthprice       string  `json:"uploadbandwidthprice"`
	Version                    string  `json:"version"`
	Windowsize                 int64   `json:"windowsize"`
}

type Interactions struct {
	TotalScans              int64  `json:"totalScans"`
	LastScan                string `json:"lastScan"`
	LastScanSuccess         bool   `json:"lastScanSuccess"`
	LostSectors             int64  `json:"lostSectors"`
	SecondToLastScanSuccess bool   `json:"secondToLastScanSuccess"`
	Uptime                  int64  `json:"uptime"`
	Downtime                int64  `json:"downtime"`
	SuccessfulInteractions  int64  `json:"successfulInteractions"`
	FailedInteractions      int64  `json:"failedInteractions"`
}

type Checks struct {
	Scanned  bool `json:"scanned"`
	Blocked  bool `json:"blocked"`
	Resolved bool `json:"resolved"`
}

type HostResponse struct {
	KnownSince       string   `json:"knownSince"`
	LastAnnouncement string   `json:"lastAnnouncement"`
	PublicKey        string   `json:"publicKey"`
	NetAddress       string   `json:"netAddress"`
	Scanned          bool     `json:"scanned"`
	Blocked          bool     `json:"blocked"`
	StoredData       int64    `json:"storedData"`
	ResolvedAddress  string   `json:"resolvedAddress"`
	Subnets          []string `json:"subnets"`
	PriceTable       PriceTable
	Settings         Settings
	Interactions     Interactions
	Checks           Checks
}

type Scan struct {
	Settings   Settings
	PriceTable PriceTable
}

type Bench struct {
	Sectors      int64  `json:"sectors"`
	Handshake    int64  `json:"handshake"`
	AppendP99    int64  `json:"appendP99"`
	ReadP99      int64  `json:"readP99"`
	Upload       int64  `json:"upload"`
	Download     int64  `json:"download"`
	UploadCost   string `json:"uploadCost"`
	DownloadCost string `json:"downloadCost"`
}

type ChainIndex struct {
	Height uint64 `json:"height"`
	ID     string `json:"id"`
}

type Consensus struct {
	Synced     bool       `json:"synced"`
	ChainIndex ChainIndex `json:"chainIndex"`
}

type Peer struct {
	Address string `json:"address"`
}

type Peers struct {
	Address string `json:"address"`
	Version string `json:"version"`
}
