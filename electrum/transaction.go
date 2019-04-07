package electrum

// TransactionBroadcast ...
func (s *Server) TransactionBroadcast(rawTx string) (string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.transaction.broadcast", []interface{}{rawTx}, resp)
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}

// TransactionResp ...
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

// Vin ...
type Vin struct {
	Coinbase  string     `json:"coinbase"`
	ScriptSig *ScriptSig `json:"scriptsig"`
	Sequence  uint32     `json:"sequence"`
	TxID      string     `json:"txid"`
	Vout      uint32     `json:"vout"`
}

// ScriptSig ...
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// Vout ...
type Vout struct {
	N            uint32       `json:"n"`
	ScriptPubkey ScriptPubkey `json:"scriptpubkey"`
	Value        float64      `json:"value"`
}

// ScriptPubkey ...
type ScriptPubkey struct {
	Addresses []string `json:"addresses,omitempty"`
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   uint32   `json:"reqsigs,omitempty"`
	Type      string   `json:"type"`
}

// TransactionGet ...
func (s *Server) TransactionGet(txHash string) (TransactionResp, error) {
	resp := &struct {
		Result TransactionResp `json:"result"`
	}{}
	err := s.request("blockchain.transaction.get", []interface{}{txHash, true}, resp)
	if err != nil {
		return TransactionResp{}, err
	}

	return resp.Result, nil
}

// TransactionGetRaw ...
func (s *Server) TransactionGetRaw(txHash string) (string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.transaction.get", []interface{}{txHash, false}, resp)
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}

// MerkleResp ...
type MerkleResp struct {
	Merkle   []string `json:"merkle"`
	Height   uint32   `json:"block_height"`
	Position uint32   `json:"pos"`
}

// TransactionGetMerkle ...
func (s *Server) TransactionGetMerkle(txHash string, height uint32) (MerkleResp, error) {
	resp := &struct {
		Result MerkleResp `json:"result"`
	}{}
	err := s.request("blockchain.transaction.get_merkle", []interface{}{txHash, height}, resp)
	if err != nil {
		return MerkleResp{}, err
	}

	return resp.Result, err
}

// TransactionHashFromPosition ...
func (s *Server) TransactionHashFromPosition(height, position uint32) (string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.transaction.id_from_pos", []interface{}{height, position, false}, resp)
	if err != nil {
		return "", err
	}

	return resp.Result, err
}

// MerkleFromPosResp ...
type MerkleFromPosResp struct {
	Hash   string   `json:"tx_hash"`
	Merkle []string `json:"merkle"`
}

// TransactionMerkleFromPosition ...
func (s *Server) TransactionMerkleFromPosition(height, position uint32) (MerkleFromPosResp, error) {
	resp := struct {
		Result MerkleFromPosResp `json:"result"`
	}{}
	err := s.request("blockchain.transaction.id_from_pos", []interface{}{height, position, true}, resp)
	if err != nil {
		return MerkleFromPosResp{}, err
	}

	return resp.Result, err
}
