package internals

import (
	"github.com/OniGbemiga/simple-bitcoin-wallet/pkg"
	"github.com/gofiber/fiber/v2"
)

func RegisterHttpHandlers() *fiber.App {
	app := fiber.New()

	app.Post("/generate-key", GenerateKey)
	app.Post("/generate-address", GenerateAddress)
	app.Post("/send-coin", SendCoin)

	return app
}

func GenerateKey(ctx *fiber.Ctx) error {

	type requestData struct {
		Environment string `json:"environment" validate:"required"`
	}

	input := new(requestData)

	if err := ctx.BodyParser(input); err != nil {
		return pkg.BadRequest(ctx, err.Error())
	}

	envVariable := pkg.BitcoinCoreEnv(input.Environment)
	if !envVariable.IsValid() {
		return pkg.BadRequest(ctx, "provide a proper environment")
	}

	walletStruct := BasicWalletStruct{
		TestENV: envVariable,
	}

	keys, err := walletStruct.GenerateKey()
	if err != nil {
		return pkg.BadRequest(ctx, err.Error())
	}

	return pkg.Success(ctx, "keys generated", map[string]interface{}{
		"privateKey": keys.PrivateKey,
		"publicKey":  keys.PublicKey,
	})
}

func GenerateAddress(ctx *fiber.Ctx) error {
	type requestData struct {
		Environment string `json:"environment" validate:"required"`
		PublicKey   string `json:"publicKey" validate:"required"`
	}

	input := new(requestData)

	if err := ctx.BodyParser(input); err != nil {
		return pkg.BadRequest(ctx, err.Error())
	}

	envVariable := pkg.BitcoinCoreEnv(input.Environment)
	if !envVariable.IsValid() {
		return pkg.BadRequest(ctx, "provide a proper environment")
	}

	walletStruct := BasicWalletStruct{
		TestENV:   envVariable,
		PublicKey: input.PublicKey,
	}

	address, err := walletStruct.CreateAddress()
	if err != nil {
		return pkg.BadRequest(ctx, err.Error())
	}

	return pkg.Success(ctx, "address generated", map[string]interface{}{
		"address": address,
	})
}

func SendCoin(ctx *fiber.Ctx) error {
	type requestData struct {
		Environment string `json:"environment" validate:"required"`
		PublicKey   string `json:"publicKey" validate:"required"`
		PrivateKey  string `json:"privateKey" validate:"required"`
		//TxIndex          uint32 `json:"txIndex" validate:"required"`
		//OutputNumber     int    `json:"OutputNumber" validate:"required"`
		Amount           float64 `json:"amount" validate:"required"`
		RecipientAddress string  `json:"recipientAddress" validate:"required"`
		SenderAddress    string  `json:"senderAddress" validate:"required"`
	}

	input := new(requestData)

	if err := ctx.BodyParser(input); err != nil {
		return pkg.BadRequest(ctx, err.Error())
	}

	envVariable := pkg.BitcoinCoreEnv(input.Environment)
	if !envVariable.IsValid() {
		return pkg.BadRequest(ctx, "provide a proper environment")
	}

	walletStruct := BasicWalletStruct{
		TestENV:          envVariable,
		PrivateKey:       input.PrivateKey,
		PublicKey:        input.PublicKey,
		Amount:           int64(input.Amount),
		RecipientAddress: input.RecipientAddress,
		SenderAddress:    input.SenderAddress,
	}

	transaction, err := walletStruct.ProcessTransaction()
	if err != nil {
		return pkg.BadRequest(ctx, err.Error())
	}

	return pkg.Success(ctx, "coin sent", map[string]interface{}{
		"hash": transaction,
	})
}
