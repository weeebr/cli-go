package figma

// EnhancedMetadata represents comprehensive component metadata
type EnhancedMetadata struct {
	Component           Component              `json:"component"`
	Properties          map[string]interface{} `json:"properties,omitempty"`
	Interactions        []interface{}          `json:"interactions,omitempty"`
	Description         string                 `json:"description,omitempty"`
	HasInteractions     bool                   `json:"has_interactions"`
	PropertyCount       int                    `json:"property_count"`
	Overrides           int                    `json:"overrides"`
	ComponentID         string                 `json:"component_id,omitempty"`
	Visible             bool                   `json:"visible"`
	Locked              bool                   `json:"locked"`
	LayoutMode          string                 `json:"layout_mode,omitempty"`
	CornerRadius        float64                `json:"corner_radius,omitempty"`
	Opacity             float64                `json:"opacity,omitempty"`
	BlendMode           string                 `json:"blend_mode,omitempty"`
	ExportSettings      []interface{}          `json:"export_settings,omitempty"`
	ChildrenCount       int                    `json:"children_count"`
	StyleReferences     map[string]string      `json:"style_references,omitempty"`
	Constraints         map[string]interface{} `json:"constraints,omitempty"`
	Effects             []interface{}          `json:"effects,omitempty"`
	Fills               []interface{}          `json:"fills,omitempty"`
	Strokes             []interface{}          `json:"strokes,omitempty"`
	StrokeWeight        float64                `json:"stroke_weight,omitempty"`
	StrokeAlign         string                 `json:"stroke_align,omitempty"`
	TextProperties      map[string]interface{} `json:"text_properties,omitempty"`
	AutoLayout          map[string]interface{} `json:"auto_layout,omitempty"`
	AbsoluteBoundingBox map[string]interface{} `json:"absolute_bounding_box,omitempty"`
	RelativeTransform   []interface{}          `json:"relative_transform,omitempty"`
	Size                map[string]interface{} `json:"size,omitempty"`
	MinWidth            float64                `json:"min_width,omitempty"`
	MaxWidth            float64                `json:"max_width,omitempty"`
	MinHeight           float64                `json:"min_height,omitempty"`
	MaxHeight           float64                `json:"max_height,omitempty"`
	Padding             map[string]float64     `json:"padding,omitempty"`
	ItemSpacing         float64                `json:"item_spacing,omitempty"`
	MainAxisAlign       string                 `json:"main_axis_align,omitempty"`
	CounterAxisAlign    string                 `json:"counter_axis_align,omitempty"`
	IsMask              bool                   `json:"is_mask"`
	CornerSmoothing     float64                `json:"corner_smoothing,omitempty"`
	StrokeCap           string                 `json:"stroke_cap,omitempty"`
	StrokeJoin          string                 `json:"stroke_join,omitempty"`
	StrokeMiterLimit    float64                `json:"stroke_miter_limit,omitempty"`
	StrokeDashPattern   []float64              `json:"stroke_dash_pattern,omitempty"`
	StrokeDashOffset    float64                `json:"stroke_dash_offset,omitempty"`
	StrokeDashCorner    string                 `json:"stroke_dash_corner,omitempty"`
	StrokeDashAlign     string                 `json:"stroke_dash_align,omitempty"`
	StrokeDashScale     float64                `json:"stroke_dash_scale,omitempty"`
	StrokeDashGap       float64                `json:"stroke_dash_gap,omitempty"`
	StrokeDashPattern2  []float64              `json:"stroke_dash_pattern_2,omitempty"`
	StrokeDashOffset2   float64                `json:"stroke_dash_offset_2,omitempty"`
	StrokeDashCorner2   string                 `json:"stroke_dash_corner_2,omitempty"`
	StrokeDashAlign2    string                 `json:"stroke_dash_align_2,omitempty"`
	StrokeDashScale2    float64                `json:"stroke_dash_scale_2,omitempty"`
	StrokeDashGap2      float64                `json:"stroke_dash_gap_2,omitempty"`
	StrokeDashPattern3  []float64              `json:"stroke_dash_pattern_3,omitempty"`
	StrokeDashOffset3   float64                `json:"stroke_dash_offset_3,omitempty"`
	StrokeDashCorner3   string                 `json:"stroke_dash_corner_3,omitempty"`
	StrokeDashAlign3    string                 `json:"stroke_dash_align_3,omitempty"`
	StrokeDashScale3    float64                `json:"stroke_dash_scale_3,omitempty"`
	StrokeDashGap3      float64                `json:"stroke_dash_gap_3,omitempty"`
	StrokeDashPattern4  []float64              `json:"stroke_dash_pattern_4,omitempty"`
	StrokeDashOffset4   float64                `json:"stroke_dash_offset_4,omitempty"`
	StrokeDashCorner4   string                 `json:"stroke_dash_corner_4,omitempty"`
	StrokeDashAlign4    string                 `json:"stroke_dash_align_4,omitempty"`
	StrokeDashScale4    float64                `json:"stroke_dash_scale_4,omitempty"`
	StrokeDashGap4      float64                `json:"stroke_dash_gap_4,omitempty"`
	StrokeDashPattern5  []float64              `json:"stroke_dash_pattern_5,omitempty"`
	StrokeDashOffset5   float64                `json:"stroke_dash_offset_5,omitempty"`
	StrokeDashCorner5   string                 `json:"stroke_dash_corner_5,omitempty"`
	StrokeDashAlign5    string                 `json:"stroke_dash_align_5,omitempty"`
	StrokeDashScale5    float64                `json:"stroke_dash_scale_5,omitempty"`
	StrokeDashGap5      float64                `json:"stroke_dash_gap_5,omitempty"`
	StrokeDashPattern6  []float64              `json:"stroke_dash_pattern_6,omitempty"`
	StrokeDashOffset6   float64                `json:"stroke_dash_offset_6,omitempty"`
	StrokeDashCorner6   string                 `json:"stroke_dash_corner_6,omitempty"`
	StrokeDashAlign6    string                 `json:"stroke_dash_align_6,omitempty"`
	StrokeDashScale6    float64                `json:"stroke_dash_scale_6,omitempty"`
	StrokeDashGap6      float64                `json:"stroke_dash_gap_6,omitempty"`
	StrokeDashPattern7  []float64              `json:"stroke_dash_pattern_7,omitempty"`
	StrokeDashOffset7   float64                `json:"stroke_dash_offset_7,omitempty"`
	StrokeDashCorner7   string                 `json:"stroke_dash_corner_7,omitempty"`
	StrokeDashAlign7    string                 `json:"stroke_dash_align_7,omitempty"`
	StrokeDashScale7    float64                `json:"stroke_dash_scale_7,omitempty"`
	StrokeDashGap7      float64                `json:"stroke_dash_gap_7,omitempty"`
	StrokeDashPattern8  []float64              `json:"stroke_dash_pattern_8,omitempty"`
	StrokeDashOffset8   float64                `json:"stroke_dash_offset_8,omitempty"`
	StrokeDashCorner8   string                 `json:"stroke_dash_corner_8,omitempty"`
	StrokeDashAlign8    string                 `json:"stroke_dash_align_8,omitempty"`
	StrokeDashScale8    float64                `json:"stroke_dash_scale_8,omitempty"`
	StrokeDashGap8      float64                `json:"stroke_dash_gap_8,omitempty"`
	StrokeDashPattern9  []float64              `json:"stroke_dash_pattern_9,omitempty"`
	StrokeDashOffset9   float64                `json:"stroke_dash_offset_9,omitempty"`
	StrokeDashCorner9   string                 `json:"stroke_dash_corner_9,omitempty"`
	StrokeDashAlign9    string                 `json:"stroke_dash_align_9,omitempty"`
	StrokeDashScale9    float64                `json:"stroke_dash_scale_9,omitempty"`
	StrokeDashGap9      float64                `json:"stroke_dash_gap_9,omitempty"`
	StrokeDashPattern10 []float64              `json:"stroke_dash_pattern_10,omitempty"`
	StrokeDashOffset10  float64                `json:"stroke_dash_offset_10,omitempty"`
	StrokeDashCorner10  string                 `json:"stroke_dash_corner_10,omitempty"`
	StrokeDashAlign10   string                 `json:"stroke_dash_align_10,omitempty"`
	StrokeDashScale10   float64                `json:"stroke_dash_scale_10,omitempty"`
	StrokeDashGap10     float64                `json:"stroke_dash_gap_10,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// Component represents a Figma component
type Component struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
}

// CreateEnhancedMetadata creates enhanced metadata from a component
func CreateEnhancedMetadata(component Component, fileKey string) *EnhancedMetadata {
	meta := &EnhancedMetadata{
		Component:           component,
		Properties:          make(map[string]interface{}),
		Interactions:        make([]interface{}, 0),
		StyleReferences:     make(map[string]string),
		Constraints:         make(map[string]interface{}),
		Effects:             make([]interface{}, 0),
		Fills:               make([]interface{}, 0),
		Strokes:             make([]interface{}, 0),
		TextProperties:      make(map[string]interface{}),
		AutoLayout:          make(map[string]interface{}),
		AbsoluteBoundingBox: make(map[string]interface{}),
		RelativeTransform:   make([]interface{}, 0),
		Size:                make(map[string]interface{}),
		Padding:             make(map[string]float64),
		Metadata:            make(map[string]interface{}),
	}

	// Add basic properties
	meta.Properties["id"] = component.ID
	meta.Properties["name"] = component.Name
	meta.Properties["type"] = component.Type
	meta.Properties["node_id"] = component.NodeID
	meta.Properties["url"] = component.URL
	meta.PropertyCount = len(meta.Properties)

	// Add metadata
	meta.Metadata["extracted_at"] = "2024-01-20T12:00:00.000Z"
	meta.Metadata["api_version"] = "v1"
	meta.Metadata["extraction_type"] = "enhanced_metadata"
	meta.Metadata["file_key"] = fileKey

	return meta
}
