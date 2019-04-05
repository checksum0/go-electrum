package electrum

// EstimateFee ...
func (s *Server) EstimateFee(target uint32) (float32, error) {
	resp := &struct {
		Result float32 `json:"result"`
	}{}

	err := s.request("blockchain.estimatefee", []interface{}{target}, resp)
	if err != nil {
		return -1, err
	}

	return resp.Result, err
}

// RelayFee ...
func (s *Server) RelayFee() (float32, error) {
	resp := &struct {
		Result float32 `json:"result"`
	}{}

	err := s.request("blockchain.relayfee", []interface{}{}, resp)
	if err != nil {
		return -1, err
	}

	return resp.Result, err
}

// FeeHistogram ...
func (s *Server) FeeHistogram() (map[uint32]uint32, error) {
	resp := &struct {
		Result map[uint32]uint32 `json:"result"`
	}{}

	err := s.request("mempool.get_fee_histogram", []interface{}{}, resp)
	if err != nil {
		return nil, err
	}

	return resp.Result, err
}
