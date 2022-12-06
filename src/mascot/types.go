package mascot

type GetBalanceBody struct {
	Version          string `json:"jsonrpc"`
	Method           string `json:"method"`
	GetBalanceParams `json:"params"`
	Id               int `json:"id"`
}

type WithdrawAndDepositBody struct {
	Version                  string `json:"jsonrpc"`
	Method                   string `json:"method"`
	WithdrawAndDepositParams `json:"params"`
	Id                       int `json:"id"`
}

type RollbackTransactionBody struct {
	Version                   string `json:"jsonrpc"`
	Method                    string `json:"method"`
	RollbackTransactionParams `json:"params"`
	Id                        int `json:"id"`
}

// Params
type GetBalanceParams struct {
	CallerId   int     `json:"callerId"`
	PlayerName string  `json:"playerName"`
	Currency   string  `json:"currency"`
	GameId     *string `json:"gameId"`
}

type WithdrawAndDepositParams struct {
	CallerId       int          `json:"callerId"`
	PlayerName     string       `json:"playerName"`
	Withdraw       int          `json:"withdraw"`
	Deposit        int          `json:"deposit"`
	Currency       string       `json:"currency"`
	TransactionRef string       `json:"transactionRef"`
	GameRoundRef   *string      `json:"gameRoundRef"`
	GameId         *string      `json:"gameId"`
	Reason         *string      `json:"reason"`
	SessionId      *string      `json:"sessionId"`
	SpinDetails    *SpinDetails `json:"spinDetails"`
}

type SpinDetails struct {
	BetType string `json:"betType"`
	WinType string `json:"winType"`
}

type RollbackTransactionParams struct {
	CallerId       int     `json:"callerId"`
	PlayerName     string  `json:"playerName"`
	TransactionRef string  `json:"transactionRef"`
	GameId         *string `json:"gameId"`
	SessionId      *string `json:"sessionId"`
	GameRoundRef   *string `json:"gameRoundRef"`
}

// Responces
type GetBalanceResponce struct {
	Version       string        `json:"jsonrpc"`
	Id            int           `json:"id"`
	BalanceResult BalanceResult `json:"result"`
}

type BalanceResult struct {
	Balance int `json:"balance"`
}

type WithdrawAndDepositResponce struct {
	Version                  string `json:"jsonrpc"`
	Id                       int    `json:"id"`
	WithdrawAndDepositResult `json:"result"`
}

type WithdrawAndDepositResult struct {
	NewBalance    int    `json:"newBalance"`
	TransactionId string `json:"transactionId"`
}

type RollbackTransactionResponce struct {
	Version string  `json:"jsonrpc"`
	Id      int     `json:"id"`
	Result  *string `json:"result"`
}
