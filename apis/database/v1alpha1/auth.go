package v1alpha1

// Auth defines the necessary conditions for connecting to the database
type Auth struct {
	//+optional

	// Username is to access the database
	Username ValueOrValueFrom `json:"username,omitempty"`

	// Password is the database access password
	Password ValueOrValueFrom `json:"password"`
}

type ValueOrValueFrom struct {
	Value     string     `json:"value,omitempty" yaml:"value,omitempty"`
	ValueFrom *ValueFrom `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

type ValueFrom struct {
	SecretKerRef *SecretKeyRef `json:"secretKeyRef,omitempty"`
	//TODO(yuanyp8): add some other secret platform, such as ssm.
}

type SecretKeyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// IsEmpty returns true if there is not a value or value and vlauefrom
func (v *ValueOrValueFrom) IsEmpty() bool {
	if v.Value != "" || v.ValueFrom != nil {
		return false
	}
	return true
}
