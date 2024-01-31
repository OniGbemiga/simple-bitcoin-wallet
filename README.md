# simple-bitcoin-wallet

## Prerequisites
 - Have an understanding of bitcoin
 - Have GO set up on your development environment
 - Have Bitcoin core set up and running on your development environment
 - Have access to Postman

 ## Getting Started
 Clone this GitHub repo to your development environment and run the commands below
 in the root directory of the folder in sequence to 
 start the project:

````
    go mod download
    go run cmd/main.go
````

When you execute the command you should see an url provided in the CLI. Copy the url and head over 
to your postman

## Actions(Endpoints)
Three endpoints are included in the project:
1. `/bitcoin/generate-key` this is a `POST` request used to generate keys, and it takes in a body
like below;
```
Request:

{
    "environment": "testnet"
}
```
```
Response
{
    "data": {
        "privateKey": "privateKey",
        "publicKey": "publicKey"
    },
    "message": "keys generated"
}
```

> Note:
> 
> `environment` can either be any of these:
> 
> 1. mainnet
> 2. testnet
> 3. regtest

2. `/bitcoin/generate-address` this is a `POST` request used to generate address from the keys,
and it takes in a body
   like below;
```
Request:

{
    "environment": "testnet",
    "publicKey": "use the public key from the response you got above"
}
```
```
Response:
{
    "data": {
        "address": "address"
    },
    "message": "address generated"
}
```
3. `/bitcoin/send-coin` this is a `POST` request used to send coin with your keys,
   and it takes in a body
   like below;
```
Request:

{
    "environment": "testnet",
    "publicKey": "publicKey",
    "privateKey": "privateKey",
    "amount": amount,
    "recipientAddress": "recipientAddress",
    "senderAddress": "senderAddress"
}
```

> Note: You can also monitor these responses in your CLI
