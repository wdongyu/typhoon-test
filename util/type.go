package util

type ManagedServiceStatus struct {
	LastUpdateBegin int64 `json:"lastUpdateBegin,omitempty"`

	LastUpdateEnd int64 `json:"lastUpdateEnd,omitempty"`

	ElapseTime int64 `json:"elapseTime,omitempty"`
}

type ManagedService struct {
	Status ManagedServiceStatus `json:"status,omitempty"`
}

type ManagedServiceList struct {
	ManagedServices []ManagedService `json:"managedServices"`
}
