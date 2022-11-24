package v1alpha1

type MysqlConnection struct {
	Access AccessAddress `json:"access"`

	Auth Auth `json:"auth"`

	//+kubebuilder:default="5.7"
	//+optional

	// Version is the mysql protocol version
	Version string `json:"version,omitempty"`

	//+kubebuilder:default={max_open_conn:200,max_idle_conn:50,max_life_time:1800,max_idle_time:600}
	//+optional
	PoolConfig ConnectPoolConfig `json:"pool_config,omitempty"`
}
