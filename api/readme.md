
## Client: ##

go run client/main.go

- start game
{"option":"game"}
- post round
{"option":"round","gameId":"7943","sequence":[{"color":"red","int":1},{"color":"green","int":3},{"color":"blue","int":2},{"color":"yellow","int":4}]}

## API Gateway Websockets: ##
- RouteSelectionExpression: $request.body.option
- Routes:
  - $connect,$disconnect,$default point to bullsandcows-connect lambda
  - game points to bullsandcows-game
  - round points to bullsandcows-round
