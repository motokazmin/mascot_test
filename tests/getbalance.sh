curl -X POST 127.0.0.1:8080/mascot/seamless -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0","method":"getBalance","params":{"callerId":1,"playerName":"player1","currency":"EUR","gameId":"riot"},"id":0}'
echo
