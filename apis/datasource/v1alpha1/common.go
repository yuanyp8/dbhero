package v1alpha1

type DataSourceType string

const (
	MYSQL   = DataSourceType("MySQL")
	POSTGRE = DataSourceType("Postgre")
	UNKNOWN = DataSourceType("UnKnown")
)

func (d DataSourceType) ToString() string {
	return string(d)
}

// AccessAddress defines 3 available connection address.
// 1) Host:Port is used to connect the backend database from this kubernetes cluster.
// 2) PrivateAddr is the inner network address, such as a `Service Name`.
// 3) PublicAddr is the outside network address, such as a Public IP addr.
type AccessAddress struct {
	// Host is the backend database access address for dbaas operator, e.g. 127.0.0.1
	Host string `json:"host"`

	//+kubebuilder:default=3306
	//+optional

	// Port is the backend database access port for dbaas operator, e.g. 3306
	Port int `json:"port,omitempty"`

	// PrivateAddr is the inner network access address
	PrivateAddr string `json:"privateAddr,omitempty"`

	// PublicAddr is the outside network access address
	PublicAddr string `json:"publicAddr,omitempty"`
}

// ConnectPoolConfig is used to maintain information about connection pool status.
type ConnectPoolConfig struct {
	//+kubebuilder:default=200
	//+optional
	MaxOpenConn int `json:"max_open_conn,omitempty"`

	//+kubebuilder:default=50
	//+optional
	MaxIdleConn int `json:"max_idle_conn,omitempty"`

	//+kubebuilder:default=1800
	//+optional
	MaxLifeTime int `json:"max_life_time,omitempty"`

	//+kubebuilder:default=600
	//+optional
	MaxIdleTime int `json:"max_idle_time,omitempty"`
}
