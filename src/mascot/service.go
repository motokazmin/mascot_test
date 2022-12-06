package mascot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mascot/src/config"
	"net/http"
	"sync"
)

type Service struct {
	Name   string
	config *config.Config
	ctx    context.Context
	wg     *sync.WaitGroup
}

func New(ctx context.Context, wg *sync.WaitGroup, config *config.Config) *Service {
	return &Service{
		Name:   "MascotService",
		config: config,
		ctx:    ctx,
		wg:     wg,
	}
}

func (s *Service) GetMascotRouter() (path string, handler http.HandlerFunc, method string) {
	return "/mascot/seamless", s.handleMascotRequest, "POST"
}

func (s *Service) writeResponse(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, data)
}

func setError(data string) string {
	return fmt.Sprintf(`{"error": "%s"}`, data)
}

func (s *Service) handleMascotRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s", err.Error())
		s.writeResponse(w, setError(err.Error()))
		return
	}

	var responce interface{}
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	switch result["method"] {
	case "getBalance":
		responce, err = getBalance(body)
		break
	case "rollbackTransaction":
		responce, err = rollbackTransaction(body)
		break
	case "withdrawAndDeposit":
		responce, err = withdrawAndDeposit(body)
		break
	default:
		break
	}

	resp, err := json.Marshal(responce)
	if err != nil {
		log.Printf("%s", err.Error())
		s.writeResponse(w, setError(err.Error()))
		return
	}
	s.writeResponse(w, string(resp))
}

func getBalance(body []byte) (interface{}, error) {
	var balanceRequest GetBalanceBody
	err := json.Unmarshal(body, &balanceRequest)
	if err != nil {
		return nil, err
	}
	return MakeGetBalanceResponce(), nil
}

func rollbackTransaction(body []byte) (interface{}, error) {
	var rollbackTransactionRequest RollbackTransactionBody
	err := json.Unmarshal(body, &rollbackTransactionRequest)
	if err != nil {
		return nil, err
	}
	return MakeGetRollbackTransactionResponce(), nil
}

func withdrawAndDeposit(body []byte) (interface{}, error) {
	var withdrawAndDepositRequest WithdrawAndDepositBody
	err := json.Unmarshal(body, &withdrawAndDepositRequest)
	if err != nil {
		return nil, err
	}
	return MakeWithdrawAndDepositResponce(), nil
}
