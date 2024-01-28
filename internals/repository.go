package internals

import (
	"fmt"
	"github.com/OniGbemiga/simple-bitcoin-wallet/pkg"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

type BasicWalletStruct struct {
	TestENV          pkg.BitcoinCoreEnv
	PrivateKey       string
	PublicKey        string
	TxIndex          uint32
	OutputNumber     int
	Amount           int64
	Fee              int
	RecipientAddress string
	SenderAddress    string
	UTXOHash         string
}

type PrivatePublicKey struct {
	PrivateKey string
	PublicKey  []byte
}

type BasicWalletRepository interface {
	GenerateKey() (string, error)
	CreateAddress() (string, error)
	ProcessTransaction() (string, error)
}

func (s BasicWalletStruct) bitcoinEnv() *chaincfg.Params {
	return s.TestENV.GetAppropriateENVVariable()
}

func (s BasicWalletStruct) GenerateKey() (*PrivatePublicKey, error) {
	//generate a new private key
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	log.Println("---- private key ---- ", privateKey)

	//convert the private key to a wallet import format
	wif, err := btcutil.NewWIF(privateKey, s.bitcoinEnv(), true)
	if err != nil {
		return nil, err
	}

	log.Println("---- private key wif ---- ", wif)

	//derive the public key
	publicKey := privateKey.PubKey()

	log.Println("---- public key ---- ", publicKey)

	return &PrivatePublicKey{
		PrivateKey: wif.String(),
		PublicKey:  publicKey.SerializeCompressed(),
	}, nil
}

func (s BasicWalletStruct) CreateAddress() (string, error) {
	//generate a p2pkh address
	address, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160([]byte(s.PublicKey)), s.bitcoinEnv())
	if err != nil {
		return "", err
	}

	log.Println("---- address ---- ", address)

	return address.EncodeAddress(), err
}

func (s BasicWalletStruct) ProcessTransaction() (string, error) {
	// Connect to the local Bitcoin testnet node
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		//Host:         "localhost:8332",
		//User:         "user",
		//Pass:         "password",
	}, nil)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//fetch utxos
	utxos, err := s.filterUtxos(client)
	if err != nil {
		return "", err
	}

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add sinput referencing the UTXO
	for _, utxo := range utxos {
		outPoint := wire.NewOutPoint(s.convertTxID(utxo.TxID), utxo.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	// Decode the destination address
	destAddr, err := btcutil.DecodeAddress(s.RecipientAddress, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}

	// Add an output to the recipient
	pkScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return "", err
	}
	tx.AddTxOut(wire.NewTxOut(s.Amount, pkScript))

	// Set transaction lock time
	tx.LockTime = 0

	//sign the transaction
	signedTransaction, err := s.signTransaction(client, tx)
	if err != nil {
		return "", err
	}

	//send the transaction
	sendTransaction, err := s.sendTransaction(client, signedTransaction)
	if err != nil {
		return "", err
	}

	return sendTransaction.String(), nil
}

func (s BasicWalletStruct) signTransaction(client *rpcclient.Client, tx *wire.MsgTx) (*wire.MsgTx, error) {
	sourcePrivateKey, err := btcutil.DecodeWIF(s.PrivateKey)
	if err != nil {
		return nil, err
	}

	for i, txIn := range tx.TxIn {
		// Fetch the previous transaction output script
		prevTxHash := txIn.PreviousOutPoint.Hash
		prevTx, err := client.GetRawTransaction(&prevTxHash)
		if err != nil {
			return nil, err
		}

		prevTxOut := prevTx.MsgTx().TxOut[txIn.PreviousOutPoint.Index]
		script, err := txscript.SignatureScript(tx, i, prevTxOut.PkScript, txscript.SigHashAll, sourcePrivateKey.PrivKey, true)
		if err != nil {
			return nil, err
		}

		txIn.SignatureScript = script
	}

	log.Println("---- signed transaction ---- ", tx)

	return tx, nil
}

func (s BasicWalletStruct) sendTransaction(client *rpcclient.Client, tx *wire.MsgTx) (*chainhash.Hash, error) {
	txHash, err := client.SendRawTransaction(tx, true)
	if err != nil {
		return nil, err
	}

	log.Println("---- send transaction ---- ", txHash)

	return txHash, nil

}

func (s BasicWalletStruct) filterUtxos(client *rpcclient.Client) ([]btcjson.ListUnspentResult, error) {
	info, err := client.GetBlockChainInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Blockchain Info:", info)
	// Fetch unspent transaction outputs (UTXOs) for the source address
	unspentUtxos, err := client.ListUnspent()
	if err != nil {
		return nil, err
	}
	var filteredUtxos []btcjson.ListUnspentResult
	for _, utxo := range unspentUtxos {
		if utxo.Address == s.SenderAddress {
			filteredUtxos = append(filteredUtxos, utxo)
		}
	}

	log.Println("---- filtered utxos ---- ", filteredUtxos)

	return filteredUtxos, nil
}

func (s BasicWalletStruct) convertTxID(txID string) *chainhash.Hash {
	hash, _ := chainhash.NewHashFromStr(txID)
	return hash
}
