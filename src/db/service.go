// Сервис базы данных postgres
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mascot/src/config"
	"sync"

	_ "github.com/lib/pq"
)

type Service struct {
	Name   string
	config *config.PostgresService
	client *sql.DB
	ctx    context.Context
	wg     *sync.WaitGroup
}

func New(ctx context.Context, wg *sync.WaitGroup, config *config.Config) *Service {
	return &Service{
		Name:   "PostgresService",
		config: &config.PostgresService,
		ctx:    ctx,
		wg:     wg,
	}
}

func (s *Service) Run() {
	connStr := "user=" + s.config.User + " dbname=" + s.config.Dbname +
		" sslmode=" + s.config.SslMode + " password=" + s.config.Password

	client, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping()
	if err != nil {
		log.Fatalf("could not connect to postgres after ping: %v", err)
	}

	s.client = client
	s.wg.Add(1)
	log.Printf("%s started", s.Name)

	<-s.ctx.Done()
	s.Close()
	s.wg.Done()
	log.Printf("%s stopped", s.Name)
}

func (s *Service) Close() {
	s.client.Close()
}

func (s *Service) GetBalance(playerName string) (int, error) {
	var balance int
	row := s.client.QueryRow(
		"SELECT balance FROM "+
			s.config.PlayersTable+"  WHERE playername=$1", playerName)
	if err := row.Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: no such player", playerName)
		}
		return 0, fmt.Errorf("%s: %v", playerName, err)
	}
	return balance, nil
}

func (s *Service) UpdateBalance(playerName string, newBalance int) {
	s.client.Exec(
		"UPDATE "+s.config.PlayersTable+
			" SET balance=$1 WHERE playername=$2", newBalance, playerName)
}

func (s *Service) InsertTransaction(
	transactionRef string, playerName string, gameId *string, sessionId *string,
	gameRoundRef *string, currency string, deposit int, id int, withdraw int,
	betType *string, winType *string, reason *string) error {

	_, err := s.client.Exec(
		"INSERT INTO "+s.config.TransactionsTable+
			" (transactionref, playername, gameid, sessionid, gameroundref, "+
			"currency, deposit, id, withdraw, betType, winType, reason)"+
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		transactionRef, playerName, gameId, sessionId, gameRoundRef, currency,
		deposit, id, withdraw, betType, winType, reason)

	return err
}

func (s *Service) RollbackTransaction(playerName string, transactionRef string) error {
	var result sql.Result

	deposit, withdraw, err := s.getTransactionDetails(playerName, transactionRef)
	if err != nil {
		return err
	}

	balance, err := s.GetBalance(playerName)
	if err != nil {
		return err
	}

	result, err = s.client.Exec(
		"DELETE FROM "+s.config.TransactionsTable+
			" WHERE playername=$1 AND transactionref=$2 "+
			"RETURNING transactionref", playerName, transactionRef)

	if err == nil {
		count, _ := result.RowsAffected()
		if count == 0 {
			return fmt.Errorf(
				"unknown transaction: playerName %s,  transactionRef %s ", playerName, transactionRef)
		}
	}

	s.UpdateBalance(playerName, balance+withdraw-deposit)
	return nil
}

func (s *Service) getTransactionDetails(playerName string, transactionRef string) (int, int, error) {
	var deposit, withdraw int
	row := s.client.QueryRow("SELECT deposit, withdraw FROM "+s.config.TransactionsTable+
		" WHERE playername=$1 AND transactionref=$2", playerName, transactionRef)
	if err := row.Scan(&deposit, &withdraw); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, fmt.Errorf("%s: no such transaction", transactionRef)
		}
		return -1, -1, fmt.Errorf("%s: %v", playerName, err)
	}
	return deposit, withdraw, nil
}
