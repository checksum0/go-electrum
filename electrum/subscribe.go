package electrum

import "encoding/json"

/* TODO:
 * - blockchain.scripthash.subscribe
 * - masternode.subscribe
 */

// SubscribeHeadersResp ...
type SubscribeHeadersResp struct {
	Height int32  `json:"height"`
	Hex    string `json:"hex"`
}

// SubscribeHeaders ...
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
