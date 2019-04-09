package electrum

import "encoding/json"

// SubscribeHeadersResp represent the response to SubscribeHeaders().
type SubscribeHeadersResp struct {
	Height int32  `json:"height"`
	Hex    string `json:"hex"`
}

// SubscribeHeaders subscribes to receive block headers notifications when new blocks are found.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-headers-subscribe
func (s *Server) SubscribeHeaders() (<-chan *SubscribeHeadersResp, error) {
	resp := &struct {
		Result *SubscribeHeadersResp `json:"result"`
	}{}

	err := s.request("blockchain.headers.subscribe", []interface{}{}, resp)
	if err != nil {
		return nil, err
	}

	respChan := make(chan *SubscribeHeadersResp, 1)
	respChan <- resp.Result

	go func() {
		for msg := range s.listenPush("blockchain.headers.subscribe") {
			if msg.err != nil {
				return
			}

			resp := &struct {
				Params []*SubscribeHeadersResp `json:"params"`
			}{}

			err := json.Unmarshal(msg.content, resp)
			if err != nil {
				return
			}

			for _, param := range resp.Params {
				respChan <- param
			}
		}
	}()

	return respChan, nil
}

// SubscribeScripthash subscribes to receive notifications when new transactions are received
// for that scripthash.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-headers-subscribe
func (s *Server) SubscribeScripthash(scripthash string) (<-chan string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.scripthash.subscribe", []interface{}{scripthash}, resp)
	if err != nil {
		return nil, err
	}

	respChan := make(chan string, 1)
	if len(resp.Result) > 0 {
		respChan <- resp.Result
	}

	go func() {
		for msg := range s.listenPush("blockchain.scripthash.subscribe") {
			if msg.err != nil {
				return
			}

			resp := &struct {
				Params []string `json:"params"`
			}{}

			err := json.Unmarshal(msg.content, resp)
			if err != nil {
				return
			}

			if resp.Params[0] == scripthash {
				respChan <- resp.Params[1]
			}
		}
	}()

	return respChan, nil
}

// SubscribeMasternode subscribes to receive notifications when a masternode status changes.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-headers-subscribe
func (s *Server) SubscribeMasternode(collateral string) (<-chan string, error) {
	resp := &basicResp{}
	err := s.request("blockchain.masternode.subscribe", []interface{}{collateral}, resp)
	if err != nil {
		return nil, err
	}

	respChan := make(chan string, 1)
	if len(resp.Result) > 0 {
		respChan <- resp.Result
	}

	go func() {
		for msg := range s.listenPush("blockchain.masternode.subscribe") {
			if msg.err != nil {
				return
			}

			resp := &struct {
				Params []string `json:"params"`
			}{}

			err := json.Unmarshal(msg.content, resp)
			if err != nil {
				return
			}

			for _, param := range resp.Params {
				respChan <- param
			}
		}
	}()

	return respChan, nil
}
