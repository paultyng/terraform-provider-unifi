// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type Map struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Lat        string  `json:"lat,omitempty"` // ^([-]?[\d]+[.]?[\d]*([eE][-+]?[\d]+)?)$
	Lng        string  `json:"lng,omitempty"` // ^([-]?[\d]+[.]?[\d]*([eE][-+]?[\d]+)?)$
	MapTypeID  string  `json:"mapTypeId"`     // satellite|roadmap|hybrid|terrain
	Name       string  `json:"name,omitempty"`
	OffsetLeft float64 `json:"offset_left,omitempty"`
	OffsetTop  float64 `json:"offset_top,omitempty"`
	Opacity    float64 `json:"opacity,omitempty"` // ^(0(\.[\d]{1,2})?|1)$|^$
	Selected   bool    `json:"selected"`
	Tilt       int     `json:"tilt,omitempty"`
	Type       string  `json:"type,omitempty"` // designerMap|imageMap|googleMap
	Unit       string  `json:"unit,omitempty"` // m|f
	Upp        float64 `json:"upp,omitempty"`
	Zoom       int     `json:"zoom,omitempty"`
}
