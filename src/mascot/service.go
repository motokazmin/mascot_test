// Бизнес логика mascot сервиса. Пакет предоставляет обработчик
// HTTP запросов.
package mascot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mascot/src/config"
	"mascot/src/db"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

type Service struct {
	Name   string
	config *config.Config
	Db     *db.Service
	ctx    context.Context
	wg     *sync.WaitGroup
}

func New(ctx context.Context, wg *sync.WaitGroup,
	config *config.Config, db *db.Service) *Service {
	return &Service{
		Name:   "MascotService",
		config: config,
		Db:     db,
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
		responce, err = s.getBalance(body)
		break
	case "rollbackTransaction":
		responce, err = s.rollbackTransaction(body)
		break
	case "withdrawAndDeposit":
		responce, err = s.withdrawAndDeposit(body)
		break
	default:
		break
	}

	if err != nil {
		s.writeResponse(w, err.Error())
		return
	}

	resp, err := json.Marshal(responce)
	if err != nil {
		log.Printf("%s", err.Error())
		s.writeResponse(w, setError(err.Error()))
		return
	}
	s.writeResponse(w, string(resp))
}

func (s *Service) getBalance(body []byte) (interface{}, error) {
	var balanceRequest GetBalanceBody
	err := json.Unmarshal(body, &balanceRequest)
	if err != nil {
		return nil, err
	}

	balance, err := s.Db.GetBalance(balanceRequest.PlayerName)
	if err != nil {
		return nil, err
	}
	return MakeGetBalanceResponce(balance,
		balanceRequest.Version, balanceRequest.Id), nil
}

func (s *Service) rollbackTransaction(body []byte) (interface{}, error) {
	var rollbackTransactionRequest RollbackTransactionBody
	err := json.Unmarshal(body, &rollbackTransactionRequest)
	if err != nil {
		return nil, err
	}

	if err = s.rollback(&rollbackTransactionRequest); err != nil {
		return nil, err
	}

	return MakeGetRollbackTransactionResponce(
		rollbackTransactionRequest.Version, rollbackTransactionRequest.Id, nil), nil
}

func (s *Service) withdrawAndDeposit(body []byte) (interface{}, error) {
	var withdrawAndDepositRequest WithdrawAndDepositBody
	err := json.Unmarshal(body, &withdrawAndDepositRequest)
	if err != nil {
		return nil, err
	}

	playerName := withdrawAndDepositRequest.PlayerName
	balance, err := s.Db.GetBalance(playerName)
	if err != nil {
		return nil, err
	}

	newBalance := balance - withdrawAndDepositRequest.Withdraw
	if newBalance < 0 {
		return nil, fmt.Errorf("%s: not enough funds", playerName)
	}

	transactionId := genUlid()
	newBalance += withdrawAndDepositRequest.Deposit

	if err = s.insertTransaction(&withdrawAndDepositRequest); err != nil {
		return nil, err
	}

	s.Db.UpdateBalance(playerName, newBalance)

	return MakeWithdrawAndDepositResponce(newBalance,
		withdrawAndDepositRequest.Version, withdrawAndDepositRequest.Id,
		transactionId), nil
}

func (s *Service) insertTransaction(r *WithdrawAndDepositBody) error {
	return s.Db.InsertTransaction(r.TransactionRef, r.PlayerName, r.GameId,
		r.SessionId, r.GameRoundRef, r.Currency, r.Deposit, r.Id,
		r.Withdraw, r.SpinDetails.BetType, r.SpinDetails.WinType, r.Reason)
}

func (s *Service) rollback(r *RollbackTransactionBody) error {
	return s.Db.RollbackTransaction(r.PlayerName, r.TransactionRef)
}

func genUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
