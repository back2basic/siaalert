package explored

import "github.com/back2basic/siaalert/scanner/types"

type Host types.Host

type Consensus struct {
	Height int    `json:"height"`
	Id     string `json:"id"`
}
