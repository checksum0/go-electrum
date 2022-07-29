package electrum

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

// AddressToElectrumScriptHex converts valid bitcoin address to electrum scriptHex sha256 encoded and reversed
func AddressToElectrumScriptHex(addressStr string) (string, error) {
	address, err := btcutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		return "", err
	}

	hashSum := sha256.Sum256(script)

	for i, j := 0, len(hashSum)-1; i < j; i, j = i+1, j-1 {
		hashSum[i], hashSum[j] = hashSum[j], hashSum[i]
	}

	return hex.EncodeToString(hashSum[:]), nil
}

// func ElectrumScriptHexToSHA256Reversed(scriptHex []byte) string {
// 	hashSum := sha256.Sum256(scriptHex)

// 	for i, j := 0, len(hashSum)-1; i < j; i, j = i+1, j-1 {
// 		hashSum[i], hashSum[j] = hashSum[j], hashSum[i]
// 	}

// 	return hex.EncodeToString(hashSum[:])
// }
