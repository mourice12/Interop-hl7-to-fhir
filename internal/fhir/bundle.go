package fhir

import "fmt"

// NewBundle creates a new transaction bundle
func NewBundle() *Bundle {
	return &Bundle{
		ResourceType: "Bundle",
		Type:         "transaction",
		Entry:        []BundleEntry{},
	}
}

// AddEntry adds a resource to the bundle
func (b *Bundle) AddEntry(resourceType, id string, resource interface{}) {
	entry := BundleEntry{
		FullURL:  fmt.Sprintf("urn:uuid:%s-%s", resourceType, id),
		Resource: resource,
	}
	b.Entry = append(b.Entry, entry)
}
