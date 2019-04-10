package electrum

import "errors"

var (
	// ErrCheckpointHeight is thrown if the checkpoint height is smaller than the block height.
	ErrCheckpointHeight = errors.New("checkpoint height must be greater or equal than block height")
)

// BlockHeaderResp represents the response to GetBlockHeader().
type BlockHeaderResp struct {
	Branch []string `json:"branch"`
	Header string   `json:"header"`
	Root   string   `json:"root"`
}

// GetBlockHeader returns the block header at a specific height.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-block-header
func (s *Server) GetBlockHeader(height uint32, checkpointHeight ...uint32) (*BlockHeaderResp, error) {
	if checkpointHeight != nil {
		if height > checkpointHeight[0] {
			return nil, ErrCheckpointHeight
		}

		resp := &struct {
			Result *BlockHeaderResp `json:"result"`
		}{}
		err := s.request("blockchain.block.header", []interface{}{height, checkpointHeight[0]}, resp)

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

// BlockHeadersResp represents the response to GetBlockHeaders().
type BlockHeadersResp struct {
	Count   uint32   `json:"count"`
	Headers string   `json:"hex"`
	Max     uint32   `json:"max"`
	Branch  []string `json:"branch,omitempty"`
	Root    string   `json:"root,omitempty"`
}

// GetBlockHeaders return a concatenated chunk of block headers.
// https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-block-headers
func (s *Server) GetBlockHeaders(startHeight, count uint32, checkpointHeight ...uint32) (*BlockHeadersResp, error) {
	resp := &struct {
		Result *BlockHeadersResp `json:"result"`
	}{}

	if checkpointHeight != nil {
		if (startHeight + (count - 1)) > checkpointHeight[0] {
			return nil, ErrCheckpointHeight
		}
		err := s.request("blockchain.block.headers", []interface{}{startHeight, count, checkpointHeight[0]}, resp)
		if err != nil {
			return nil, err
		}

		return resp.Result, err
	}

	err := s.request("blockchain.block.headers", []interface{}{startHeight, count}, resp)
	if err != nil {
		return nil, err
	}

	return resp.Result, err
}
