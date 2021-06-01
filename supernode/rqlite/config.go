package rqlite

const (
	defaultHTTPAddress            = "0.0.0.0:4001"
	defaultRaftAddress            = "0.0.0.0:4002"
	defaultJoinAttempts           = 5
	defaultJoinInterval           = "5s"
	defaultRaftNotVoter           = false
	defaultRaftHeartbeatTimeout   = "1s"
	defaultRaftElectionTimeout    = "1s"
	defaultRaftApplyTimeout       = "10s"
	defaultRaftOpenTimeout        = "120s"
	defaultRaftWaitForLeader      = true
	defaultRaftSnapThreshold      = 8192
	defaultRaftSnapInterval       = "30s"
	defaultRaftLeaderLeaseTimeout = "0s"
	defaultRaftShutdownOnRemove   = false
	defaultRaftLogLevel           = "error"
	defaultCompressionSize        = 150
	defaultCompressionBatch       = 5
	defaultOnDisk                 = false
	defaultDataDir                = "data"
	defaultNodeX509Cert           = "cert.pem"
	defaultNodeX509Key            = "key.pem"
)

// Config contains settings of the rqlite server
type Config struct {
	// http server bind address. For HTTPS, set X.509 cert and key
	HTTPAddress string `mapstructure:"http_address" json:"http_address,omitempty"`
	// advertised HTTP address. If not set, same as HTTP server
	HTTPAdvertiseAddress string `mapstructure:"http_advertise_address" json:"http_advertise_address"`
	// set source IP address during Join request
	JoinSourceIP string `mapstructure:"join_source_ip" json:"join_source_ip,omitempty"`
	// unique name for node. If not set, set to hostname
	NodeID string `mapstructure:"node_id" json:"node_id,omitempty"`
	// raft communication bind address
	RaftAddress string `mapstructure:"raft_address" json:"raft_address,omitempty"`
	// advertised Raft communication address. If not set, same as Raft bind
	RaftAdvertiseAddress string `mapstructure:"raft_advertise_address" json:"raft_advertise_address,omitempty"`
	// comma-delimited list of nodes, through which a cluster can be joined (proto://host:port)"
	JoinAddress string `mapstructure:"join_address" json:"join_address,omitempty"`
	// number of join attempts to make
	JoinAttempts int `mapstructure:"join_attempts" json:"join_attempts,omitempty"`
	// period between join attempts
	JoinInterval string `mapstructure:"join_interval" json:"join_interval,omitempty"`
	// set Discovery Service URL
	DiscoveryURL string `mapstructure:"discovery_url" json:"discovery_url,omitempty"`
	// set Discovery ID. If not set, Discovery Service not used
	DiscoveryID string `mapstructure:"discovery_id" json:"discovery_id,omitempty"`
	// configure as non-voting node
	RaftNotVoter bool `mapstructure:"raft_not_voter" json:"raft_not_voter,omitempty"`
	// raft heartbeat timeout
	RaftHeartbeatTimeout string `mapstructure:"raft_heartbeat_timeout" json:"raft_heartbeat_timeout,omitempty"`
	// raft election timeout
	RaftElectionTimeout string `mapstructure:"raft_election_timeout" json:"raft_election_timeout,omitempty"`
	// raft apply timeout
	RaftApplyTimeout string `mapstructure:"raft_apply_timeout" json:"raft_apply_timeout,omitempty"`
	// time for initial Raft logs to be applied. Use 0s duration to skip wait
	RaftOpenTimeout string `mapstructure:"raft_open_timeout" json:"raft_open_timeout,omitempty"`
	// node waits for a leader before answering requests
	RaftWaitForLeader bool `mapstructure:"raft_wait_for_leader" json:"raft_wait_for_leader,omitempty"`
	// number of outstanding log entries that trigger snapshot
	RaftSnapThreshold uint64 `mapstructure:"raft_snap_threshold" json:"raft_snap_threshold,omitempty"`
	// snapshot threshold check interval
	RaftSnapInterval string `mapstructure:"raft_snap_interval" json:"raft_snap_interval,omitempty"`
	// raft leader lease timeout. Use 0s for Raft default
	RaftLeaderLeaseTimeout string `mapstructure:"raft_leader_lease_timeout" json:"raft_leader_lease_timeout,omitempty"`
	// shutdown Raft if node removed
	RaftShutdownOnRemove bool `mapstructure:"raft_shutdown_on_remove" json:"raft_shutdown_on_remove,omitempty"`
	// minimum log level for Raft module
	RaftLogLevel string `mapstructure:"raft_log_level" json:"raft_log_level,omitempty"`
	// request query size for compression attempt
	CompressionSize int `mapstructure:"compression_size" json:"compression_size,omitempty"`
	// request batch threshold for compression attempt
	CompressionBatch int `mapstructure:"compression_batch" json:"compression_batch,omitempty"`
	// use an on-disk SQLite database
	OnDisk bool `mapstructure:"on_disk" json:"on_disk,omitempty"`
	// SQLite DSN parameters. E.g. "cache=shared&mode=memory"
	DNS string `mapstructure:"dns" json:"dns,omitempty"`
	// data directory for rqlite
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`
	// Support deprecated TLS versions 1.0 and 1.1
	TLS1011 bool `mapstructure:"tls1011" json:"tls1011,omitempty"`
	// Path to root X.509 certificate for HTTP endpoint
	X509CACert string `mapstructure:"http_ca_cert" json:"http_ca_cert,omitempty"`
	// Path to X.509 certificate for HTTP endpoint
	X509Cert string `mapstructure:"http_cert" json:"http_cert,omitempty"`
	// Path to X.509 private key for HTTP endpoint
	X509Key string `mapstructure:"http_key" json:"http_key,omitempty"`
	// Skip verification of remote HTTPS cert when joining cluster
	NoVerify bool `mapstructure:"http_no_verify" json:"http_no_verify,omitempty"`
	// Enable node-to-node encryption
	NodeEncrypt bool `mapstructure:"node_encrypt" json:"node_encrypt,omitempty"`
	// Path to root X.509 certificate for node-to-node encryption
	NodeX509CACert string `mapstructure:"node_ca_cert" json:"node_ca_cert,omitempty"`
	// Path to X.509 certificate for node-to-node encryption
	NodeX509Cert string `mapstructure:"node_cert" json:"node_cert,omitempty"`
	// Path to X.509 private key for node-to-node encryption
	NodeX509Key string `mapstructure:"node_key" json:"node_key,omitempty"`
	// Skip verification of a remote node cert
	NodeNoVerify bool `mapstructure:"node_no_verify" json:"node_no_verify,omitempty"`
	// Path to authentication and authorization file. If not set, not enabled
	AuthFile string `mapstructure:"auth_file" json:"auth_file,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		HTTPAddress:            defaultHTTPAddress,
		RaftAddress:            defaultRaftAddress,
		JoinAttempts:           defaultJoinAttempts,
		JoinInterval:           defaultJoinInterval,
		RaftNotVoter:           defaultRaftNotVoter,
		RaftHeartbeatTimeout:   defaultRaftHeartbeatTimeout,
		RaftElectionTimeout:    defaultRaftElectionTimeout,
		RaftApplyTimeout:       defaultRaftApplyTimeout,
		RaftOpenTimeout:        defaultRaftOpenTimeout,
		RaftWaitForLeader:      defaultRaftWaitForLeader,
		RaftSnapThreshold:      defaultRaftSnapThreshold,
		RaftSnapInterval:       defaultRaftSnapInterval,
		RaftLeaderLeaseTimeout: defaultRaftLeaderLeaseTimeout,
		RaftShutdownOnRemove:   defaultRaftShutdownOnRemove,
		RaftLogLevel:           defaultRaftLogLevel,
		CompressionSize:        defaultCompressionSize,
		CompressionBatch:       defaultCompressionBatch,
		OnDisk:                 defaultOnDisk,
		DataDir:                defaultDataDir,
		NodeX509Cert:           defaultNodeX509Cert,
		NodeX509Key:            defaultNodeX509Key,
	}
}
