package sdk

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"

	"github.com/back2basic/siadata/siaalert/config"
	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/mail"
)

var (
	HostCache = make(map[string]HostDocument)
	Mutex     sync.RWMutex
)

var (
	instance *AppwriteDatabaseService
	once     sync.Once
)

type AppwriteDatabaseService struct {
	Client *databases.Databases
}
type DocumentData map[string]interface{}

// DatabaseService defines the interface for Appwrite database operations.
type DatabaseService interface {
	ListDatabases() (*models.DatabaseList, error)
	ListCollections(databaseID string) (*models.CollectionList, error)
	ListDocuments(databaseID, collectionID string, queries []string) (*models.DocumentList, error)
	GetDocument(databaseID, collectionID, documentID string) (*models.Document, error)
	CreateHostDocument(documentID string, data Host) (*models.Document, error)
	UpdateHostDocument(documentID string, data Host) (*models.Document, error)
	CreateCheckDocument(documentID string, data Check) (*models.Document, error)
	UpdateCheckDocument(documentID string, data Check) (*models.Document, error)
	CreateStatusDocument(documentID string, data Status) (*models.Document, error)
	UpdateStatusDocument(documentID string, data Status) (*models.Document, error)
	CreateRhp2Document(documentID string, data Rhp2) (*models.Document, error)
	UpdateRhp2Document(documentID string, data Rhp2) (*models.Document, error)
	CreateRhp3Document(documentID string, data Rhp3) (*models.Document, error)
	UpdateRhp3Document(documentID string, data Rhp3) (*models.Document, error)
}

// ListDatabases calls the Appwrite SDK's List method.
func (ads *AppwriteDatabaseService) ListDatabases() (*models.DatabaseList, error) {
	return ads.Client.List()
}

// ListCollections calls the Appwrite SDK's ListCollections method.
func (ads *AppwriteDatabaseService) ListCollections(databaseID string) (*models.CollectionList, error) {
	return ads.Client.ListCollections(databaseID)
}

// ListDocuments wraps the Appwrite list documents call.
func (ads *AppwriteDatabaseService) ListDocuments(databaseID, collectionID string, queries []string) (*models.DocumentList, error) {
	return ads.Client.ListDocuments(databaseID, collectionID, ads.Client.WithListDocumentsQueries(queries))
}

// GetDocument wraps the Appwrite get document call.
func (ads *AppwriteDatabaseService) GetDocument(databaseID, collectionID, documentID string) (*models.Document, error) {
	return ads.Client.GetDocument(databaseID, collectionID, documentID)
}

// CreateDocument wraps the Appwrite create document call.
func (ads *AppwriteDatabaseService) CreateHostDocument(documentID string, data Host) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.CreateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColHosts.Id, documentID, data)
}

// UpdateDocument wraps the Appwrite update document call.
func (ads *AppwriteDatabaseService) UpdateHostDocument(documentID string, data Host) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.UpdateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColHosts.Id, documentID, ads.Client.WithUpdateDocumentData(data))
}

// CreateDocument wraps the Appwrite create document call.
func (ads *AppwriteDatabaseService) CreateStatusDocument(documentID string, data Status) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.CreateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColStatus.Id, documentID, data)
}

// UpdateDocument wraps the Appwrite update document call.
func (ads *AppwriteDatabaseService) UpdateStatusDocument(documentID string, data Status) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.UpdateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColStatus.Id, documentID, ads.Client.WithUpdateDocumentData(data))
}

// CreateDocument wraps the Appwrite create document call.
func (ads *AppwriteDatabaseService) CreateCheckDocument(documentID string, data Check) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.CreateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColCheck.Id, documentID, data)
}

// UpdateDocument wraps the Appwrite update document call.
func (ads *AppwriteDatabaseService) UpdateCheckDocument(documentID string, data Check) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.UpdateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColCheck.Id, documentID, ads.Client.WithUpdateDocumentData(data))
}

func (ads *AppwriteDatabaseService) CreateRhp2Document(documentID string, data Rhp2) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.CreateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColRhp2.Id, documentID, data)
}

func (ads *AppwriteDatabaseService) UpdateRhp2Document(documentID string, data Rhp2) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.UpdateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColRhp2.Id, documentID, ads.Client.WithUpdateDocumentData(data))
}

func (ads *AppwriteDatabaseService) CreateRhp3Document(documentID string, data Rhp3) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.CreateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColRhp3.Id, documentID, data)
}

func (ads *AppwriteDatabaseService) UpdateRhp3Document(documentID string, data Rhp3) (*models.Document, error) {
	cfg := config.GetConfig()
	return ads.Client.UpdateDocument(cfg.Appwrite.Database.Id, cfg.Appwrite.ColRhp3.Id, documentID, ads.Client.WithUpdateDocumentData(data))
}

func GetAppwriteDatabaseService() *AppwriteDatabaseService {
	once.Do(func() {
		instance = &AppwriteDatabaseService{}
	})
	return instance
}

func WriteTohostCache(key string, value HostDocument) {
	Mutex.Lock() // Lock for writing
	HostCache[key] = value
	Mutex.Unlock()
}

func ReadFromHostCache(key string) (HostDocument, bool) {
	Mutex.RLock() // Lock for reading
	value, ok := HostCache[key]
	Mutex.RUnlock()
	return value, ok
}

func PrepareAppwrite(cfg *config.Config) DatabaseService {
	// Initialize the Appwrite SDK client.
	SdkClient := appwrite.NewClient(
		appwrite.WithEndpoint(cfg.Appwrite.Endpoint),
		appwrite.WithProject(cfg.Appwrite.Project),
		appwrite.WithKey(cfg.Appwrite.Key),
	)

	// Create a new Databases instance from the client.
	sdkDB := appwrite.NewDatabases(SdkClient)

	// Wrap the concrete type in our interface implementation.
	dbSvc := &AppwriteDatabaseService{Client: sdkDB}
	return dbSvc
}

func ListDatabases(dbSvc DatabaseService) *models.DatabaseList {
	databases, err := dbSvc.ListDatabases()
	if err != nil {
		panic(err)
	}
	fmt.Println("Databases loaded:", len(databases.Databases))
	return databases
}

func PrepareDatabase(cfg *config.Config, dbSvc DatabaseService) *models.Database {
	dbs := ListDatabases(dbSvc)
	var foundDB *models.Database

	for i := 0; i < dbs.Total; i++ {
		if dbs.Databases[i].Name == cfg.Network.Name {
			foundDB = &dbs.Databases[i]
			break
		}
	}

	if foundDB == nil {
		panic("DB Not Found")
	}
	fmt.Println("Database loaded:", foundDB.Name, foundDB.Id)
	return foundDB
}

func PrepareCollection(dbSvc DatabaseService, databaseID string) (hosts, status, alert, check, rhp2, rhp3 *models.Collection) {
	allCollections, err := dbSvc.ListCollections(databaseID)
	if err != nil {
		panic(err)
	}

	for _, c := range allCollections.Collections {
		switch c.Name {
		case "hosts":
			hosts = &c
			fmt.Println("Collection loaded:", hosts.Name, hosts.Id)
		case "status":
			status = &c
			fmt.Println("Collection loaded:", status.Name, status.Id)
		case "alert":
			alert = &c
			fmt.Println("Collection loaded:", alert.Name, alert.Id)
		case "check":
			check = &c
			fmt.Println("Collection loaded:", check.Name, check.Id)
		case "rhp2":
			rhp2 = &c
			fmt.Println("Collection loaded:", rhp2.Name, rhp2.Id)
		case "rhp3":
			rhp3 = &c
			fmt.Println("Collection loaded:", rhp3.Name, rhp3.Id)
		}
	}
	return
}

// GetHostByPublicKey retrieves a host document by its public key.
func GetHostByPublicKey(databaseID, collectionID, publicKey string) (*models.DocumentList, error) {
	db := GetAppwriteDatabaseService()
	queries := []string{
		query.Equal("publicKey", publicKey),
	}

	response, err := db.ListDocuments(databaseID, collectionID, queries)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func CheckHost(host explored.Host) (HostDocument, error) {
	cfg := config.GetConfig()
	// fmt.Println("Run check", host.NetAddress)
	result, err := GetHostByPublicKey(cfg.Appwrite.Database.Id, cfg.Appwrite.ColHosts.Id, host.PublicKey)

	if err != nil {
		fmt.Println(err)
		return HostDocument{}, err
	}
	// fmt.Println("Found" + fmt.Sprint(len(result.Documents)))
	if len(result.Documents) == 0 {
		// fmt.Println("Create new host")
		return CreateHost(host)
	}
	var foundHost HostList
	result.Decode(&foundHost)
	if len(result.Documents) == 1 {
		// fmt.Println("Update host")
		UpdateHost(host, result.Documents[0].Id, foundHost.Documents[0])
	}
	if len(result.Documents) > 1 {
		panic("multiple publickeys found for host" + host.PublicKey)
	}
	return foundHost.Documents[0], nil
}

func CreateHost(host explored.Host) (HostDocument, error) {
	db := GetAppwriteDatabaseService()
	// fmt.Println("Creating host", host)
	sdkHost := Host{
		PublicKey:              host.PublicKey,
		V2:                     host.V2,
		NetAddress:             host.NetAddress,
		V2NetAddresses:         host.V2NetAddresses.Address,
		V2NetAddressesProto:    host.V2NetAddresses.Protocol,
		CountryCode:            host.CountryCode,
		KnownSince:             host.KnownSince.Format(time.RFC3339),
		LastScan:               host.LastScan.Format(time.RFC3339),
		LastScanSuccessful:     host.LastScanSuccessful,
		LastAnnouncement:       host.LastAnnouncement.Format(time.RFC3339),
		TotalScans:             host.TotalScans,
		SuccessfulInteractions: host.SuccessfulInteractions,
		FailedInteractions:     host.FailedInteractions,
		Online:                 true,
		OnlineSince:            host.LastAnnouncement.Format(time.RFC3339),
		OfflineSince:           "",
		Error:                  "",
	}

	doc, err := db.CreateHostDocument(id.Unique(), sdkHost)
	var resp HostDocument
	doc.Decode(&resp)
	if err != nil {
		// fmt.Println(host.PublicKey)
		return HostDocument{}, err
	}
	// fmt.Println("Host created", resp.Id)
	return resp, nil
}

func UpdateHost(host explored.Host, hostId string, foundHost HostDocument) {
	db := GetAppwriteDatabaseService()
	// fmt.Println("Updating host")
	sdkHost := Host{
		PublicKey:              host.PublicKey,
		V2:                     host.V2,
		NetAddress:             host.NetAddress,
		V2NetAddresses:         host.V2NetAddresses.Address,
		V2NetAddressesProto:    host.V2NetAddresses.Protocol,
		CountryCode:            host.CountryCode,
		KnownSince:             host.KnownSince.Format(time.RFC3339),
		LastScan:               host.LastScan.Format(time.RFC3339),
		LastScanSuccessful:     host.LastScanSuccessful,
		LastAnnouncement:       host.LastAnnouncement.Format(time.RFC3339),
		TotalScans:             host.TotalScans,
		SuccessfulInteractions: host.SuccessfulInteractions,
		FailedInteractions:     host.FailedInteractions,
		Online:                 foundHost.Online,
		OnlineSince:            foundHost.OnlineSince,
		OfflineSince:           foundHost.OfflineSince,
		Error:                  foundHost.Error,
	}

	_, err := db.UpdateHostDocument(hostId, sdkHost)

	if err != nil {
		fmt.Println(host.PublicKey)
		fmt.Println(err)
	}
	// fmt.Println("Host updated", resp.Id)
}

func UpdateNetAddress(host HostDocument) {
	db := GetAppwriteDatabaseService()

	docData := Host{
		PublicKey:              host.PublicKey,
		NetAddress:             host.NetAddress,
		V2:                     host.V2,
		V2NetAddresses:         host.V2NetAddresses,
		V2NetAddressesProto:    host.V2NetAddressesProto,
		CountryCode:            host.CountryCode,
		KnownSince:             host.KnownSince,
		LastScan:               host.LastScan,
		LastScanSuccessful:     host.LastScanSuccessful,
		LastAnnouncement:       host.LastAnnouncement,
		TotalScans:             host.TotalScans,
		SuccessfulInteractions: host.SuccessfulInteractions,
		FailedInteractions:     host.FailedInteractions,
		Error:                  host.Error,
		Online:                 host.Online,
		OnlineSince:            host.OnlineSince,
		OfflineSince:           host.OfflineSince,
	}

	_, err := db.UpdateHostDocument(host.Id, docData)
	if err != nil {
		fmt.Println("UpdateNetAddress:", err)
	}
}

func CheckUpdateStatus(hostId, netAddress, error string, online bool) {
	db := GetAppwriteDatabaseService()
	cfg := config.GetConfig()
	found, err := db.GetDocument(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColHosts.Id,
		hostId,
	)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting status, Failed to check update status")
		return
	}

	if found == nil {
		fmt.Println("Error getting document, Failed to check update status")
		return
	}
	var hostDoc HostDocument
	found.Decode(&hostDoc)

	// Update document
	if !online {
		if hostDoc.Online {
			// !! HOST WAS ONLINE AND IS NOW OFFLINE - More actions here for future plugins.
			UpdateHostStatus(error, "", time.Now().Format(time.RFC3339), false, hostDoc)
			PrepareEmails(hostId, netAddress, "Offline")
			fmt.Printf("Host %s is offline since %s\n", netAddress, time.Now().Format(time.RFC3339))
		} else {
			UpdateHostStatus(error, "", hostDoc.OfflineSince, false, hostDoc)
		}
	} else {
		if !hostDoc.Online {
			// !! HOST WAS OFFLINE AND IS NOW ONLINE - More actions here for future plugins.
			UpdateHostStatus(error, time.Now().Format(time.RFC3339), "", true, hostDoc)
			PrepareEmails(hostId, netAddress, "Online")
			fmt.Printf("Host %s is online since %s\n", netAddress, time.Now().Format(time.RFC3339))
		}
	}
}

func UpdateHostStatus(error, onlineSince, offlineSince string, online bool, host HostDocument) {
	db := GetAppwriteDatabaseService()

	docData := Host{
		PublicKey:              host.PublicKey,
		V2:                     host.V2,
		NetAddress:             host.NetAddress,
		V2NetAddresses:         host.V2NetAddresses,
		V2NetAddressesProto:    host.V2NetAddressesProto,
		CountryCode:            host.CountryCode,
		KnownSince:             host.KnownSince,
		LastScan:               host.LastScan,
		LastScanSuccessful:     host.LastScanSuccessful,
		LastAnnouncement:       host.LastAnnouncement,
		TotalScans:             host.TotalScans,
		SuccessfulInteractions: host.SuccessfulInteractions,
		FailedInteractions:     host.FailedInteractions,
		Error:                  error,
		Online:                 online,
		OnlineSince:            onlineSince,
		OfflineSince:           offlineSince,
	}

	doc, err := db.UpdateHostDocument(
		host.Id,
		docData,
	)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to update status")
		return
	}
	var resp HostDocument
	doc.Decode(&resp)
	Mutex.Lock()
	HostCache[host.PublicKey] = resp
	Mutex.Unlock()
}

func PrepareEmails(hostId, netAddress, status string) {
	subscribers := GetSubscribers(hostId)
	if len(subscribers) == 0 {
		return
	}
	for _, subscriber := range subscribers {
		mail.SendMail(subscriber, netAddress, status)
	}
}

func GetSubscribers(hostId string) []string {
	db := GetAppwriteDatabaseService()
	cfg := config.GetConfig()
	found, err := db.ListDocuments(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColAlert.Id,
		[]string{
			query.Equal("hostId", hostId),
		},
	)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	var alertList AlertList
	found.Decode(&alertList)

	// loop t and find email subscribers
	var subscribers []string
	for _, doc := range alertList.Documents {
		if doc.Type != "email" {
			continue
		}
		subscribers = append(subscribers, doc.Sender)
	}

	return subscribers
}

func UpdateStatus(state explored.Consensus) {
	db := GetAppwriteDatabaseService()
	cfg := config.GetConfig()
	found, err := db.ListDocuments(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColStatus.Id,
		[]string{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	var statusList StatusList
	found.Decode(&statusList)
	status := Status{
		Height: state.Index.Height,
	}
	if statusList.Total == 0 {
		_, err := db.CreateStatusDocument(id.Unique(), status)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := db.UpdateStatusDocument(statusList.Documents[0].Id, status)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GetCheck(hostId string) (CheckDocument, error) {
	db := GetAppwriteDatabaseService()
	cfg := config.GetConfig()
	found, err := db.ListDocuments(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColCheck.Id,
		[]string{
			query.Equal("hostId", hostId),
		},
	)
	if err != nil {
		fmt.Println(err)
		return CheckDocument{}, err
	}
	var checkList CheckList
	found.Decode(&checkList)
	if len(checkList.Documents) == 0 {
		return CheckDocument{}, errors.New("no check found for host " + hostId)
	}
	return checkList.Documents[0], nil
}

func UpdateCheck(params CheckParams, wg *sync.WaitGroup, task chan TaskCheckDoc) {
	db := GetAppwriteDatabaseService()
	cfg := config.GetConfig()
	found, err := db.ListDocuments(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColCheck.Id,
		[]string{
			query.Equal("hostId", params.HostId),
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	var release string
	var checkList CheckList
	found.Decode(&checkList)
	if len(checkList.Documents) == 0 {
		release = ""
	} else {
		release = checkList.Documents[0].Release
	}
	check := Check{
		HostId:             params.HostId,
		V4Addr:             params.V4,
		V6Addr:             params.V6,
		Rhp2Port:           params.Rhp2Port,
		Rhp2V4Delay:        float64(params.Rhp2v4Delay.Milliseconds()),
		Rhp2V6Delay:        float64(params.Rhp2v6Delay.Milliseconds()),
		Rhp2V4:             params.Rhp2v4,
		Rhp2V6:             params.Rhp2v6,
		Rhp3Port:           params.Rhp3Port,
		Rhp3V4Delay:        float64(params.Rhp3v4Delay.Milliseconds()),
		Rhp3V6Delay:        float64(params.Rhp3v6Delay.Milliseconds()),
		Rhp3V4:             params.Rhp3v4,
		Rhp3V6:             params.Rhp3v6,
		Rhp4Port:           params.Rhp4Port,
		Rhp4V4Delay:        float64(params.Rhp4v4Delay.Milliseconds()),
		Rhp4V6Delay:        float64(params.Rhp4v6Delay.Milliseconds()),
		Rhp4V4:             params.Rhp4v4,
		Rhp4V6:             params.Rhp4v6,
		AcceptingContracts: params.AcceptingContracts,
		Release:            release,
	}

	if checkList.Total == 0 {
		wg.Add(1)
		task <- TaskCheckDoc{ID: 1, Job: "createCheck", CheckID: id.Unique(), Check: check}
		// _, err := db.CreateCheckDocument(id.Unique(), check)
		// if err != nil {
		// 	fmt.Println(err)
		// }
	} else {
		// TODO: check if value have changed and send mail with array that have changed
		wg.Add(1)
		task <- TaskCheckDoc{ID: 1, Job: "updateCheck", CheckID: checkList.Documents[0].Id, Check: check}
		// _, err := db.UpdateCheckDocument(checkList.Documents[0].Id, check)
		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
}

func UpdateRelease(id string, check Check) {
	db := GetAppwriteDatabaseService()
	_, err := db.UpdateCheckDocument(id, check)
	if err != nil {
		fmt.Println(err)
	}
}

func UpdateRhp(host explored.Host) {
	db := GetAppwriteDatabaseService()
	cfg := config.GetConfig()
	found, err := GetHostByPublicKey(cfg.Appwrite.Database.Id, cfg.Appwrite.ColHosts.Id, host.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	var hostList HostList
	found.Decode(&hostList)

	foundRhp2, err := db.ListDocuments(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColRhp2.Id,
		[]string{
			query.Equal("hostId", hostList.Documents[0].Id),
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	foundRhp3, err := db.ListDocuments(
		cfg.Appwrite.Database.Id,
		cfg.Appwrite.ColRhp3.Id,
		[]string{
			query.Equal("hostId", hostList.Documents[0].Id),
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	params2 := Rhp2{
		AcceptingContracts:   host.Settings.AcceptingContracts,
		MaxDownloadBatchSize: host.Settings.MaxDownloadBatchSize,
		MaxDuration:          host.Settings.MaxDuration,
		MaxReviseBatchSize:   host.Settings.MaxReviseBatchSize,
		RemainingStorage:     host.Settings.RemainingStorage,
		TotalStorage:         host.Settings.TotalStorage,
		RevisionNumber:       host.Settings.RevisionNumber,
		Version:              host.Settings.Version,
		Release:              host.Settings.Release,
		SiaMuxPort:           host.Settings.SiaMuxPort,
		HostId:               hostList.Documents[0].Id,
	}
	if foundRhp2.Total == 0 {
		_, err = db.CreateRhp2Document(id.Unique(), params2)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = db.UpdateRhp2Document(foundRhp2.Documents[0].Id, params2)
	}
	if err != nil {
		fmt.Println(err)
	}

	params3 := Rhp3{
		HostBlockHeight: host.PriceTable.HostBlockHeight,
		HostId:          hostList.Documents[0].Id,
	}
	if foundRhp3.Total == 0 {
		_, err = db.CreateRhp3Document(id.Unique(), params3)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = db.UpdateRhp3Document(foundRhp3.Documents[0].Id, params3)
	}
	if err != nil {
		fmt.Println(err)
	}

}
