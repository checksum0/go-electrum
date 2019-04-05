package electrum

import "errors"

var (
	// ErrCheckpointHeight ...
	ErrCheckpointHeight = errors.New("checkpoint height must be greater or equal than block height")
)

// BlockHeaderResp ...
type BlockHeaderResp struct {
	Branch []string `json:"branch"`
	Header string   `json:"header"`
	Root   string   `json:"root"`
}

// BlockHeader ...
func (s *Server) BlockHeader(height uint32, checkpointHeight ...uint32) (*BlockHeaderResp, error) {
	if checkpointHeight != nil {
		if height > checkpointHeight[0] {
			return nil, ErrCheckpointHeight
		}

		resp := &struct {
			Result *BlockHeaderResp `json:"result"`
		}{}
		err := s.request("blockchain.block.header", []interface{}{height, checkpointHeight}, resp)

		return resp.Result, err
	}

	resp := &basicResp{}
	err := s.request("blockchain.block.header", []interface{}{height}, resp)
	if err != nil {
		return nil, err
	}

	result := &BlockHeaderResp{
		Branch: nil,
		Header: resp.Result,
		Root:   "",
	}

	return result, err
}

// BlockHeadersResp ...
type BlockHeadersResp struct {
	Count   uint32   `json:"count"`
	Headers string   `json:"hex"`
	Max     uint32   `json:"max"`
	Branch  []string `json:"branch,omitempty"`
	Root    string   `json:"root,omitempty"`
}

// BlockHeaders ...
func (s *Server) BlockHeaders(startHeight, count uint32, checkpointHeight ...uint32) (*BlockHeadersResp, error) {
	resp := &struct {
		Result *BlockHeadersResp `json:"result"`
	}{}

	if checkpointHeight != nil {
		if (startHeight + (count - 1)) > checkpointHeight[0] {
			return nil, ErrCheckpointHeight
		}
		err := s.request("blockchain.block.headers", []interface{}{startHeight, count, checkpointHeight}, resp)
		if err != nil {
			return nil, err
		}

		return resp.Result, err
	}

	err := s.request("blockchain.block.headers", []interface{}{startHeight, count, checkpointHeight}, resp)
	if err != nil {
		return nil, err
	}

	return resp.Result, err
}
