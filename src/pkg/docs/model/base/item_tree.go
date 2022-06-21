package base

type ItemTree struct {
	ItemStructure
	Items []Item `json:"item"`
}
