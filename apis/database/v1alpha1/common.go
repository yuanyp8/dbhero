package v1alpha1

type CharacterSet string
type Collation string

const (
	UTF8MB4           = CharacterSet("utf8mb4")
	LATIN1            = CharacterSet("latin1")
	UTF8_GENERAL_CI   = Collation("utf8mb4_general_ci")
	LATIN1_GENERAL_CS = Collation("latin1_general_cs")
)

// AccessAddress defines 3 available connection address.
// 1) PrivateAddr is the inner network address, such as a `Service Name`.
// 2) PublicAddr is the outside network address, such as a Public IP addr.
type AccessAddress struct {
	// PrivateAddr is the inner network access address
	PrivateAddr string `json:"privateAddr,omitempty"`

	// PublicAddr is the outside network access address
	PublicAddr string `json:"publicAddr,omitempty"`
}
