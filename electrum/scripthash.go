package electrum

/* TODO:
 * - blockchain.scripthash.utxos (version 1.5)
 * - blockchain.scripthash.history (version 1.5)
 */

// BalanceResp represents the response to GetBalance().
type BalanceResp struct {
	Confirmed   float64 `json:"confirmed"`
	Unconfirmed float64 `json:"unconfirmed"`
}

// GetBalance returns the confirmed and unconfirmed balance for a scripthash.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-scripthash-get-balance
func (s *Server) GetBalance(scripthash string) (BalanceResp, error) {
	resp := &struct {
		Result BalanceResp `json:"result"`
	}{}
	err := s.request("blockchain.scripthash.get_balance", []interface{}{scripthash}, resp)
	if err != nil {
		return BalanceResp{}, err
	}

	return resp.Result, err
}

// MempoolResp represents the reponse to GetHistory() and GetMempool().
type MempoolResp struct {
	Hash   string `json:"tx_hash"`
	Height int32  `json:"height"`
	Fee    uint32 `json:"fee,omitempty"`
}

// GetHistory returns the confirmed and unconfirmed history for a scripthash.
func (s *Server) GetHistory(scripthash string) ([]*MempoolResp, error) {
	resp := &struct {
		Result []*MempoolResp `json:"result"`
	}{}
	err := s.request("blockchain.scripthash.get_history", []interface{}{scripthash}, resp)
	if err != nil {
		return []*MempoolResp{}, err
	}

	return resp.Result, err
}

// GetMempool returns the unconfirmed transacations of a scripthash.
func (s *Server) GetMempool(scripthash string) ([]*MempoolResp, error) {
	resp := &struct {
		Result []*MempoolResp `json:"result"`
	}{}
	err := s.request("blockchain.scripthash.get_mempool", []interface{}{scripthash}, resp)
	if err != nil {
		return []*MempoolResp{}, err
	}

	return resp.Result, err
}

// UnspentResp represents the response to ListUnspent()
type UnspentResp struct {
	Height   uint32 `json:"height"`
	Position uint32 `json:"tx_pos"`
	Hash     string `json:"tx_hash"`
	Value    uint64 `json:"value"`
}

// ListUnspent returns an ordered list of UTXOs for a scripthash.
func (s *Server) ListUnspent(scripthash string) ([]*UnspentResp, error) {
	resp := &struct {
		Result []*UnspentResp `json:"result"`
	}{}
	err := s.request("blockchain.scripthash.listunspent", []interface{}{scripthash}, resp)
	if err != nil {
		return []*UnspentResp{}, err
	}

	return resp.Result, err
}
