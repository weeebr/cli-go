package figma

import "fmt"

// FigmaNode represents a node in the Figma file tree
type FigmaNode struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Children []FigmaNode `json:"children,omitempty"`

	// For full metadata extraction
	ComponentPropertyDefinitions map[string]interface{} `json:"componentPropertyDefinitions,omitempty"`
	ComponentProperties          map[string]interface{} `json:"componentProperties,omitempty"`
	Interactions                 []interface{}          `json:"interactions,omitempty"`
	Description                  string                 `json:"description,omitempty"`
}

// FigmaFile represents the API response from /v1/files/{key}
type FigmaFile struct {
	Document FigmaNode `json:"document"`
	Name     string    `json:"name"`
}

// FullMetadata represents detailed component information
type FullMetadata struct {
	Component       Component              `json:"component"`
	Properties      map[string]interface{} `json:"properties,omitempty"`
	Interactions    []interface{}          `json:"interactions,omitempty"`
	Description     string                 `json:"description,omitempty"`
	HasInteractions bool                   `json:"has_interactions"`
	PropertyCount   int                    `json:"property_count"`
}

// ComponentFromCacheData safely converts cache data to Component
func ComponentFromCacheData(data map[string]interface{}) (Component, error) {
	var comp Component
	var ok bool

	// Safe type assertions with error handling
	if comp.NodeID, ok = data["node_id"].(string); !ok {
		return comp, fmt.Errorf("invalid node_id in cache data")
	}
	if comp.Name, ok = data["name"].(string); !ok {
		return comp, fmt.Errorf("invalid name in cache data")
	}
	if comp.Type, ok = data["type"].(string); !ok {
		return comp, fmt.Errorf("invalid type in cache data")
	}

	// URL is optional, so don't fail if it's missing
	if comp.URL, ok = data["url"].(string); !ok {
		comp.URL = "" // Set empty string if URL is not present
	}

	return comp, nil
}
