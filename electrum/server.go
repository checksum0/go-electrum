package electrum

// Ping ...
func (s *Server) Ping() error {
	err := s.request("server.ping", []interface{}{}, nil)

	return err
}

// ServerAddPeer ...
func (s *Server) ServerAddPeer() error {
	return ErrNotImplemented
}

// ServerBanner ...
func (s *Server) ServerBanner() (string, error) {
	resp := &basicResp{}
	err := s.request("server.banner", []interface{}{}, resp)

	return resp.Result, err
}

// ServerDonation ...
func (s *Server) ServerDonation() (string, error) {
	resp := &basicResp{}
	err := s.request("server.donation_address", []interface{}{}, resp)

	return resp.Result, err
}

type host struct {
	TCPPort uint16 `json:"tcp_port,omitempty"`
	SSLPort uint16 `json:"ssl_port,omitempty"`
}

type featuresResp struct {
	GenesisHash   string          `json:"genesis_hash"`
	Hosts         map[string]host `json:"hosts"`
	ProtocolMax   string          `json:"protocol_max"`
	ProtocolMin   string          `json:"protocol_min"`
	Pruning       bool            `json:"pruning,omitempty"`
	ServerVersion string          `json:"server_version"`
	HashFunction  string          `json:"hash_function"`
}

// ServerFeatures ...
func (s *Server) ServerFeatures() (interface{}, error) {
	resp := &struct {
		Result *featuresResp `json:"result"`
	}{}
	err := s.request("server.features", []interface{}{}, resp)

	return &resp.Result, err
}

// ServerPeers ...
func (s *Server) ServerPeers() (interface{}, error) {
	resp := &struct {
		Result [][]interface{} `json:"result"`
	}{}
	err := s.request("server.peers.subscribe", []interface{}{}, resp)

	return resp.Result, err
}

// ServerVersion ...
func (s *Server) ServerVersion() (serverVer, protocolVer string, err error) {
	resp := &struct {
		Result []string `json:"result"`
	}{}
	err = s.request("server.version", []interface{}{ClientVersion, ProtocolVersion}, resp)
	if err != nil {
		serverVer = ""
		protocolVer = ""
	} else {
		serverVer = resp.Result[0]
		protocolVer = resp.Result[1]
	}

	return
}
