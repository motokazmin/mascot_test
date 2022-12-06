package mascot

func MakeGetBalanceResponce() interface{} {
	return GetBalanceResponce{
		Version: "2.0",
		Id:      0,
		BalanceResult: BalanceResult{
			Balance: 113,
		},
	}
}

func MakeGetRollbackTransactionResponce() interface{} {
	return RollbackTransactionResponce{
		Version: "2.0",
		Id:      0,
	}
}

func MakeWithdrawAndDepositResponce() interface{} {
	return WithdrawAndDepositResponce{
		Version: "2.0",
		Id:      0,
		WithdrawAndDepositResult: WithdrawAndDepositResult{
			NewBalance:    1101,
			TransactionId: "bla-bla-bla",
		},
	}
}
