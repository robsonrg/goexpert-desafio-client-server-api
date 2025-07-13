# Desafio Pós Go Expert - Client Server API

> Este projeto contém a solução para o desafio sobre webserver http, contextos, banco de dados e manipulação de arquivos com Go da pós-graduação `Go Expert` da FullCycle.
 
(...)
Os requisitos para cumprir este desafio:
 
- O `client.go` deverá realizar uma requisição HTTP no `server.go` solicitando a cotação do dólar.
 
- O `server.go` deverá consumir a API contendo o câmbio de Dólar e Real no endereço: `https://economia.awesomeapi.com.br/json/last/USD-BRL` e em seguida deverá retornar no formato JSON o resultado para o cliente.
 
- Usando o package "context", o `server.go` deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
 
- O `client.go` precisará receber do `server.go` apenas o valor atual do câmbio (campo `bid` do JSON). Utilizando o package "context", o `client.go` terá um timeout máximo de 300ms para receber o resultado do `server.go`.
 
- Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
 
- O `client.go` terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
 
- O endpoint necessário gerado pelo `server.go` para este desafio será: `/cotacao` e a porta a ser utilizada pelo servidor HTTP será a `8080`.

---

# Como executar

1. Executar o aplicação `server`

```sh
cd server
go mod tidy
go run server.go
```

Essa aplicação expõe o endpoint `/cotacao`, na porta `:8080`, no qual retorna a cotação atual do dolar:

```json 
//request: GET http://localhost:8080/cotacao
//response
{
    "bid": "5.5569"
}
```

2. Executar a aplicação `client`

```sh
cd client
go run client.go
```