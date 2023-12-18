package solr

// FieldDef represents a Solr SchemaUpdates SolrFieldDef (does not cover all possible fields)
type FieldDef struct {
	// No omitempty to avoid dropping set values that happen to coincide with Go default initialization
	Name         string `json:"name"`
	Type         string `json:"type"`
	MultiValued  bool   `json:"multiValued"`
	Stored       bool   `json:"stored"`
	Indexed      bool   `json:"indexed"`
	Required     bool   `json:"required"`
	DocValues    bool   `json:"docValues"`
	Uninvertible bool   `json:"uninvertible"`
}

// CopyFieldDef represents a Solr SchemaUpdates copy Field
type CopyFieldDef struct {
	Source string `json:"source"`
	// Note that we can state multiple destinations
	Destination []string `json:"dest"`
	// MaxChars is deliberately left out when set to Go default to avoid sending '0' if the field is not set
	MaxChars uint `json:"maxChars,omitempty"`
}

// DynamicFieldDef represents a  Solr SchemaUpdates dynamic Field - alias for Field since the fields are the same
type DynamicFieldDef FieldDef

// SchemaUpdates represents the extra fields to be added to a blank schema
type SchemaUpdates struct {
	FieldDefs        []FieldDef
	CopyFieldDefs    []CopyFieldDef
	DynamicFieldDefs []DynamicFieldDef
}

type ClusterInfo struct {
	Collections map[string]interface{} `json:"collections"`
	Aliases     map[string]interface{} `json:"aliases"`
	Roles       map[string]interface{} `json:"roles"`
	LiveNodes   []string               `json:"live_nodes"`
}

type ClusterStatus struct {
	Cluster ClusterInfo `json:"cluster"`
}
