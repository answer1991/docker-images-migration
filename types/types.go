package types

type Source struct {
	Version string `json:"version"`

	//FromRegistryDomain   string `json:"fromRegistryDomain"`
	TargetRegistryDomain string `json:"targetRegistryDomain"`

	Images []string `json:"images"`
}

type AuthInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
