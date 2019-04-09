package electrum

import "encoding/json"

/* TODO:
 * - blockchain.scripthash.subscribe
 * - masternode.subscribe
 */

// SubscribeHeadersResp represent the response to SubscribeHeaders().
type SubscribeHeadersResp struct {
	Height int32  `json:"height"`
	Hex    string `json:"hex"`
}

// SubscribeHeaders subscribe to receive block headers notifications
// when new blocks are found.
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
