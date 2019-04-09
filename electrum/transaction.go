package electrum

// TransactionBroadcast sends a raw transaction to the remote server to
// be broadcasted on the server network.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-transaction-broadcast
func (s *Server) TransactionBroadcast(rawTx string) (string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.transaction.broadcast", []interface{}{rawTx}, resp)
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}

// TransactionResp represents the response to TransactionGet().
type TransactionResp struct {
	Blockhash     string     `json:"blockhash"`
	Blocktime     uint64     `json:"blocktime"`
	Confirmations int32      `json:"confirmations"`
	Hash          string     `json:"hash"`
	Hex           string     `json:"hex"`
	Locktime      uint32     `json:"locktime"`
	Size          uint32     `json:"size"`
	Time          uint64     `json:"time"`
	Version       uint32     `json:"version"`
	Vin           []Vin      `json:"vin"`
	Vout          []Vout     `json:"vout"`
	Merkle        MerkleResp `json:"merkle,omitempty"` // For protocol v1.5 and up.
}

// Vin represents the input side of a transaction.
type Vin struct {
	Coinbase  string     `json:"coinbase"`
	ScriptSig *ScriptSig `json:"scriptsig"`
	Sequence  uint32     `json:"sequence"`
	TxID      string     `json:"txid"`
	Vout      uint32     `json:"vout"`
}

// ScriptSig represents the signature script for that transaction input.
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// Vout represents the output side of a transaction.
type Vout struct {
	N            uint32       `json:"n"`
	ScriptPubkey ScriptPubkey `json:"scriptpubkey"`
	Value        float64      `json:"value"`
}

// ScriptPubkey represents the script of that transaction output.
type ScriptPubkey struct {
	Addresses []string `json:"addresses,omitempty"`
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   uint32   `json:"reqsigs,omitempty"`
	Type      string   `json:"type"`
}

// TransactionGet gets the detailed information for a transaction.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-transaction-get
func (s *Server) TransactionGet(txHash string) (*TransactionResp, error) {
	resp := &struct {
		Result *TransactionResp `json:"result"`
	}{}
	err := s.request("blockchain.transaction.get", []interface{}{txHash, true}, resp)
	if err != nil {
		return &TransactionResp{}, err
	}

	return resp.Result, nil
}

// TransactionGetRaw gets a raw encoded transaction.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-transaction-get
func (s *Server) TransactionGetRaw(txHash string) (string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.transaction.get", []interface{}{txHash, false}, resp)
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}

// MerkleResp represents the response TransactionGetMerkle().
type MerkleResp struct {
	Merkle   []string `json:"merkle"`
	Height   uint32   `json:"block_height"`
	Position uint32   `json:"pos"`
}

// TransactionGetMerkle returns the merkle proof for a confirmed transaction.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-transaction-get-merkle
func (s *Server) TransactionGetMerkle(txHash string, height uint32) (*MerkleResp, error) {
	resp := &struct {
		Result *MerkleResp `json:"result"`
	}{}
	err := s.request("blockchain.transaction.get_merkle", []interface{}{txHash, height}, resp)
	if err != nil {
		return &MerkleResp{}, err
	}

	return resp.Result, err
}

// TransactionHashFromPosition returns the transaction hash for a specific position in a block.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-transaction-id-from-pos
func (s *Server) TransactionHashFromPosition(height, position uint32) (string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.transaction.id_from_pos", []interface{}{height, position, false}, resp)
	if err != nil {
		return "", err
	}

	return resp.Result, err
}

// MerkleFromPosResp represents the response to TransactionMerkleFromPosition().
type MerkleFromPosResp struct {
	Hash   string   `json:"tx_hash"`
	Merkle []string `json:"merkle"`
}

// TransactionMerkleFromPosition returns the merkle proof for a specific position in a block.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-transaction-id-from-pos
func (s *Server) TransactionMerkleFromPosition(height, position uint32) (*MerkleFromPosResp, error) {
	resp := struct {
		Result *MerkleFromPosResp `json:"result"`
	}{}
	err := s.request("blockchain.transaction.id_from_pos", []interface{}{height, position, true}, resp)
	if err != nil {
		return &MerkleFromPosResp{}, err
	}

	return resp.Result, err
}
