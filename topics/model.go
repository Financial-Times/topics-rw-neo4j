package topics

type Topic struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier,omitempty"`
	Type          string `json:"type,omitempty"`
}

type TopicLink struct {
	ApiUrl string `json:"apiUrl"`
}
