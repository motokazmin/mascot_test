package mascot

func MakeGetBalanceResponce(balance int, version string, id int) interface{} {
	return GetBalanceResponce{
		Version: version,
		Id:      id,
		BalanceResult: BalanceResult{
			Balance: balance,
		},
	}
}

func MakeGetRollbackTransactionResponce(version string, id int, result *string) interface{} {
	return RollbackTransactionResponce{
		Version: version,
		Id:      id,
		Result:  result,
	}
}

func MakeWithdrawAndDepositResponce(newBalance int, version string,
	id int, transactionId string) interface{} {
	return WithdrawAndDepositResponce{
		Version: version,
		Id:      id,
		WithdrawAndDepositResult: WithdrawAndDepositResult{
			NewBalance:    newBalance,
			TransactionId: transactionId,
		},
	}
}
