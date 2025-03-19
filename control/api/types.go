package api

type Otp struct {
	publicKey string
	secret    string
	expire    string
	email     string
}

type Alert struct {
	publicKey string
	alert     string
	sender    string
}
