package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

type RequestBody struct {
	Network   string `json:"network"`
	PublicKey string `json:"publicKey"`
	Email     string `json:"email"`
	Secret    string `json:"secret"`
	Expire    string `json:"expire"`
}

type Response struct {
	Message   string `json:"message"`
	PublicKey string `json:"publicKey,omitempty"`
	Email     string `json:"email,omitempty"`
	Address   string `json:"address,omitempty"`
}

type Host struct {
	PublicKey              string `json:"publicKey"`
	V2                     bool   `json:"v2,omitempty"`
	NetAddress             string `json:"netAddress,omitempty"`
	V2NetAddresses         string `json:"v2NetAddresses,omitempty"`
	V2NetAddressesProto    string `json:"v2NetAddressesProto,omitempty"`
	CountryCode            string `json:"countryCode,omitempty"`
	KnownSince             string `json:"knownSince,omitempty"`
	LastScan               string `json:"lastScan,omitempty"`
	LastScanSuccessful     bool   `json:"lastScanSuccessful,omitempty"`
	LastAnnouncement       string `json:"lastAnnouncement,omitempty"`
	TotalScans             uint64 `json:"totalScans,omitempty"`
	SuccessfulInteractions uint64 `json:"successfulInteractions,omitempty"`
	FailedInteractions     uint64 `json:"failedInteractions,omitempty"`

	Error        string `json:"error"`
	Online       bool   `json:"online"`
	OnlineSince  string `json:"onlineSince"`
	OfflineSince string `json:"offlineSince"`
}

type HostDocument struct {
	models.Document
	Host
}

type HostList struct {
	Documents []HostDocument `json:"documents"`
	Total     uint64         `json:"total"`
}

type Alert struct {
	HostId string `json:"hostId"`
	Type   string `json:"type"`
	Sender string `json:"sender"`
}

type AlertDocument struct {
	models.Document
	Alert
}

type AlertList struct {
	Documents []AlertDocument `json:"documents"`
	Total     uint64          `json:"total"`
}

func sendMail(publicKey, email, otp, expire, network string) (string, error) {
	host := os.Getenv("SMTP_HOST")
	port := 587
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	from := "Sia Host Alert <info@euregiohosting.nl>"
	to := email
	subject := "Setup Alerts for your Sia Host"
	htmlBody := fmt.Sprintf("<b>Control Alerts for your Sia Host</b><br><a href=\"https://siaalert.euregiohosting.nl/auth?otp=%s&network=%s&email=%s&expire=%s&publicKey=%s\">https://siaalert.euregiohosting.nl/auth?otp=%s&network=%s&email=%s&expire=%s&publicKey=%s</a>", otp, network, url.QueryEscape(email), expire, publicKey, otp, network, url.QueryEscape(email), expire, publicKey)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, user, pass)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}
	return "Email sent successfully", nil
}

func getRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func Main(Context openruntimes.Context) openruntimes.Response {
	switch Context.Req.Method {
	case http.MethodPost:
		var reqBody RequestBody
		Context.Req.BodyJson(&reqBody)
		Context.Log("Post Parameters are " + fmt.Sprintf("%+v", reqBody))
		Context.Log("Post Parameters are OK")
		sec, err := getRandomString(32)
		if err != nil {
			return Context.Res.Json(Response{
				Message: "Error generating random string",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
		}
		exp := time.Now().Format(time.RFC3339)
		Context.Log("Secret is " + sec)
		Context.Log("Expire is " + exp)

		client := appwrite.NewClient(
			appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
			appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
			appwrite.WithKey(os.Getenv("APPWRITE_FUNCTION_API_KEY")),
		)
		database := appwrite.NewDatabases(client)

		databaseID := ""
		collectionID := ""
		if reqBody.Network == "mainnet" {
			databaseID = os.Getenv("APPWRITE_FUNCTION_DATABASE_MAIN")
			collectionID = os.Getenv("APPWRITE_FUNCTION_COLLECTION_OTP_MAIN")
		} else {
			databaseID = os.Getenv("APPWRITE_FUNCTION_DATABASE_ZEN")
			collectionID = os.Getenv("APPWRITE_FUNCTION_COLLECTION_OTP_ZEN")
		}

		_, err = database.CreateDocument(databaseID, collectionID, id.Unique(), map[string]interface{}{
			"publicKey": reqBody.PublicKey,
			"email":     reqBody.Email,
			"secret":    sec,
			"expire":    exp,
		})
		if err != nil {
			Context.Error(err)
			return Context.Res.Json(Response{
				Message: "Could not create OTP",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
		}

		info, err := sendMail(reqBody.PublicKey, reqBody.Email, sec, exp, reqBody.Network)
		if err != nil {
			Context.Error(err)
			return Context.Res.Json(Response{
				Message: "Could not send email",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
		}
		Context.Log(info)
		return Context.Res.Json(Response{
			Message: "Email sent",
		}, Context.Res.WithStatusCode(http.StatusCreated)) // 201 response

	case http.MethodPut:
		var reqBody RequestBody
		Context.Req.BodyJson(&reqBody)
		Context.Log("Put Parameters are " + fmt.Sprintf("%+v", reqBody))

		client := appwrite.NewClient(
			appwrite.WithEndpoint(os.Getenv("APPWRITE_FUNCTION_API_ENDPOINT")),
			appwrite.WithProject(os.Getenv("APPWRITE_FUNCTION_PROJECT_ID")),
			appwrite.WithKey(os.Getenv("APPWRITE_FUNCTION_API_KEY")),
		)
		database := appwrite.NewDatabases(client)

		databaseID := ""
		collectionID := ""
		if reqBody.Network == "mainnet" {
			databaseID = os.Getenv("APPWRITE_FUNCTION_DATABASE_MAIN")
			collectionID = os.Getenv("APPWRITE_FUNCTION_COLLECTION_OTP_MAIN")
		} else {
			databaseID = os.Getenv("APPWRITE_FUNCTION_DATABASE_ZEN")
			collectionID = os.Getenv("APPWRITE_FUNCTION_COLLECTION_OTP_ZEN")
		}

		documents, err := database.ListDocuments(
			databaseID,
			collectionID,
			database.WithListDocumentsQueries(
				[]string{
					query.Equal("publicKey", reqBody.PublicKey),
					query.Equal("email", reqBody.Email),
					query.Equal("secret", reqBody.Secret),
				},
			),
		)
		if err != nil || len(documents.Documents) == 0 {
			Context.Error(err)
			return Context.Res.Json(Response{
				Message: "Invalid OTP",
			}, Context.Res.WithStatusCode(http.StatusBadRequest)) // 400 response
		}
		docID := documents.Documents[0].Id

		// Check if is created iin the last 24 hours
		createdAt, err := time.Parse(time.RFC3339, documents.Documents[0].CreatedAt)
		if err != nil {
			Context.Error(err)
			return Context.Res.Json(Response{
				Message: "Error parsing created at",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
		}

		if time.Since(createdAt).Hours() > 24 {
			Context.Error("OTP has expired")
			return Context.Res.Json(Response{
				Message: "OTP has expired",
			}, Context.Res.WithStatusCode(http.StatusBadRequest)) // 400 response
		}

		_, err = database.DeleteDocument(databaseID, collectionID, docID)
		if err != nil {
			return Context.Res.Json(Response{
				Message: "Error deleting document",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
		}

		collectionIDHosts := os.Getenv("APPWRITE_FUNCTION_COLLECTION_HOSTS_MAIN")
		if reqBody.Network == "zen" {
			collectionIDHosts = os.Getenv("APPWRITE_FUNCTION_COLLECTION_HOSTS_ZEN")
		}

		hostDocs, err := database.ListDocuments(
			databaseID,
			collectionIDHosts,
			database.WithListDocumentsQueries(
				[]string{
					query.Equal("publicKey", reqBody.PublicKey),
				},
			),
		)
		if err != nil || len(hostDocs.Documents) == 0 {
			return Context.Res.Json(Response{
				Message: "Host not found",
			}, Context.Res.WithStatusCode(http.StatusBadRequest)) // 400 response
		}
		var hostList HostList
		hostDocs.Decode(&hostList)

		collectionIDAlerts := os.Getenv("APPWRITE_FUNCTION_COLLECTION_ALERT_MAIN")
		if reqBody.Network == "zen" {
			collectionIDAlerts = os.Getenv("APPWRITE_FUNCTION_COLLECTION_ALERT_ZEN")
		}

		alertDocs, err := database.ListDocuments(
			databaseID,
			collectionIDAlerts,
			database.WithListDocumentsQueries(
				[]string{
					query.Equal("hostId", hostList.Documents[0].Id),
					query.Equal("type", "email"),
					query.Equal("sender", reqBody.Email),
				},
			),
		)
		if err != nil {
			Context.Error(err)
			return Context.Res.Json(Response{
				Message: "Error listing documents",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
		}

		if len(alertDocs.Documents) == 0 {
			_, err := database.CreateDocument(databaseID, collectionIDAlerts, id.Unique(), Alert{
				HostId: hostDocs.Documents[0].Id,
				Type:   "email",
				Sender: reqBody.Email,
			})
			if err != nil {
				Context.Error(err)
				return Context.Res.Json(Response{
					Message: "Error creating document",
				}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
			}
			return Context.Res.Json(Response{
				Message:   "enabled",
				PublicKey: reqBody.PublicKey,
				Email:     reqBody.Email,
				Address:   hostList.Documents[0].NetAddress,
			}, Context.Res.WithStatusCode(http.StatusOK)) // 200 response
		} else {
			for _, doc := range alertDocs.Documents {
				_, err := database.DeleteDocument(databaseID, collectionIDAlerts, doc.Id)
				if err != nil {
					Context.Error(err)
					return Context.Res.Json(Response{
						Message: "Error deleting document",
					}, Context.Res.WithStatusCode(http.StatusInternalServerError)) // 500 response
				}
			}
			return Context.Res.Json(Response{
				Message:   "disabled",
				PublicKey: reqBody.PublicKey,
				Email:     reqBody.Email,
				Address:   hostList.Documents[0].NetAddress,
			}, Context.Res.WithStatusCode(http.StatusOK)) // 200 response
		}

	default:
		return Context.Res.Json(Response{
			Message: "OK",
		}, Context.Res.WithStatusCode(http.StatusOK)) // 200 response
	}
}
