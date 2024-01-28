package pkg

import "github.com/btcsuite/btcd/chaincfg"

type BitcoinCoreEnv string

const (
	MAINNET BitcoinCoreEnv = "mainnet"
	REGTEST BitcoinCoreEnv = "regtest"
	TESTNET BitcoinCoreEnv = "testnet"
)

func (receiver BitcoinCoreEnv) GetAppropriateENVVariable() *chaincfg.Params {
	switch receiver {
	case MAINNET:
		return &chaincfg.MainNetParams
	case REGTEST:
		return &chaincfg.RegressionNetParams
	case TESTNET:
		return &chaincfg.TestNet3Params
	default:
		return &chaincfg.TestNet3Params
	}
}

func (receiver BitcoinCoreEnv) IsValid() bool {
	switch receiver {
	case MAINNET, REGTEST, TESTNET:
		return true
	default:
		return false
	}
}
