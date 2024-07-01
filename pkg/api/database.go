package api

type Database struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Href        string            `json:"href"`
	Annotation  string            `json:"annotation"`
	Permissions map[string]string `json:"permissions,omitempty"`
	Extensions  []string          `json:"extensions,omitempty"`
	Locale      string            `json:"locale,omitempty"`
}
