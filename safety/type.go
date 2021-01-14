package main

type PatchManagedService struct {
	// Destination cluster subset for the traffic to route
	RouteSubset string `json:"routeSubset,omitempty"`

	// Ready to revoke obsolete subset
	RevokeObsoleteSubset *bool `json:"revokeObsoleteSubset,omitempty"`
}
