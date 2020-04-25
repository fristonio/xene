package response

// RegistryItemsFromPrefix is the response of list prefix query on registry items.
type RegistryItemsFromPrefix struct {
	Count int `json:"count" example:"2"`

	// Items contains the Serialized kvstore items
	Items []string `json:"items"`
}

// RegistryItem is the reponse of registry item get on the apiserver.
type RegistryItem struct {
	// Items contains the Serialized kvstore item
	Item string `json:"item"`
}
