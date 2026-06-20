package main

import "sigs.k8s.io/controller-tools/pkg/markers"

var ProviderMarkers = []*markers.Definition{
	markers.Must(markers.MakeDefinition("provider", markers.DescribesType, nil)),
}

func Register(reg *markers.Registry) error {
	return markers.RegisterAll(reg, ProviderMarkers...)
}
