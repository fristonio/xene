package response

import (
	types "github.com/fristonio/xene/pkg/types/v1alpha1"
)

// RegistryItemsFromPrefix is the response of list prefix query on registry items.
type RegistryItemsFromPrefix struct {
	Count int `json:"count" example:"2"`

	Items []types.KVPairStruct `json:"items"`
}

// RegistryItem is the reponse of registry item get on the apiserver.
type RegistryItem struct {
	Item types.KVPairStruct `json:"item" example:"Workflow Document"`
}
