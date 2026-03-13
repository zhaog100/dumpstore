// Package schema defines ZFS dataset property metadata used by the API and frontend.
// Adding a new property: append one entry to allProperties; the UI, allowed-list,
// and OS filtering all update automatically.
package schema

import "slices"

// Option is a single selectable value for a Property.
type Option struct {
	Value   string `json:"value"`
	Label   string `json:"label"`
	Default bool   `json:"default,omitempty"`
}

// Property describes one ZFS dataset property as understood by the UI.
type Property struct {
	Name      string   `json:"name"`
	Label     string   `json:"label"`
	InputType string   `json:"input_type"` // "select" or "text"
	Editable  bool     `json:"editable"`   // shown in edit dialog
	Create    bool     `json:"create"`     // shown in create dialog
	AppliesTo []string `json:"applies_to"` // "filesystem", "volume", or both
	OS        []string `json:"-"`          // empty = all OS; otherwise filter by runtime.GOOS
	Options   []Option `json:"options,omitempty"`
}

var both = []string{"filesystem", "volume"}
var fsOnly = []string{"filesystem"}
var volOnly = []string{"volume"}

var inheritOpt = Option{Value: "", Label: "— inherit —"}

var allProperties = []Property{
	{
		Name: "compression", Label: "Compression",
		InputType: "select", Editable: true, Create: true, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "lz4", Label: "lz4", Default: true},
			{Value: "zstd", Label: "zstd"},
			{Value: "gzip", Label: "gzip"},
			{Value: "on", Label: "on"},
			{Value: "off", Label: "off"},
		},
	},
	{
		Name: "sync", Label: "Sync",
		InputType: "select", Editable: true, Create: true, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "standard", Label: "standard"},
			{Value: "always", Label: "always"},
			{Value: "disabled", Label: "disabled"},
		},
	},
	{
		Name: "atime", Label: "Access Time",
		InputType: "select", Editable: true, Create: true, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "on", Label: "on"},
			{Value: "off", Label: "off"},
		},
	},
	{
		Name: "exec", Label: "Exec",
		InputType: "select", Editable: true, Create: true, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "on", Label: "on"},
			{Value: "off", Label: "off"},
		},
	},
	{
		Name: "copies", Label: "Copies",
		InputType: "select", Editable: true, Create: true, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "1", Label: "1"},
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3"},
		},
	},
	{
		Name: "xattr", Label: "Extended Attributes",
		InputType: "select", Editable: true, Create: true, AppliesTo: fsOnly,
		Options: []Option{
			inheritOpt,
			{Value: "on", Label: "on"},
			{Value: "sa", Label: "sa (kernel)"},
			{Value: "off", Label: "off"},
		},
	},
	{
		Name: "dedup", Label: "Dedup",
		InputType: "select", Editable: true, Create: true, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "off", Label: "off"},
			{Value: "on", Label: "on"},
			{Value: "verify", Label: "verify"},
		},
	},
	{
		Name: "readonly", Label: "Read-only",
		InputType: "select", Editable: true, Create: false, AppliesTo: both,
		Options: []Option{
			inheritOpt,
			{Value: "on", Label: "on"},
			{Value: "off", Label: "off"},
		},
	},
	{
		Name: "recordsize", Label: "Record Size",
		InputType: "select", Editable: true, Create: true, AppliesTo: fsOnly,
		Options: []Option{
			inheritOpt,
			{Value: "512", Label: "512 B"},
			{Value: "4K", Label: "4K"},
			{Value: "8K", Label: "8K"},
			{Value: "16K", Label: "16K"},
			{Value: "32K", Label: "32K"},
			{Value: "64K", Label: "64K"},
			{Value: "128K", Label: "128K"},
			{Value: "256K", Label: "256K"},
			{Value: "512K", Label: "512K"},
			{Value: "1M", Label: "1M"},
		},
	},
	{
		Name: "volblocksize", Label: "Block Size",
		InputType: "select", Editable: false, Create: true, AppliesTo: volOnly,
		Options: []Option{
			inheritOpt,
			{Value: "512", Label: "512 B"},
			{Value: "4K", Label: "4K"},
			{Value: "8K", Label: "8K"},
			{Value: "16K", Label: "16K"},
			{Value: "32K", Label: "32K"},
			{Value: "64K", Label: "64K"},
			{Value: "128K", Label: "128K"},
		},
	},
	{
		Name: "quota", Label: "Quota",
		InputType: "text", Editable: true, Create: true, AppliesTo: fsOnly,
	},
	{
		Name: "mountpoint", Label: "Mountpoint",
		InputType: "text", Editable: true, Create: true, AppliesTo: fsOnly,
	},
}

// All returns every property definition regardless of OS.
func All() []Property {
	return allProperties
}

// ForOS returns properties applicable to the given GOOS value.
// Properties with an empty OS list apply to all platforms.
func ForOS(goos string) []Property {
	result := make([]Property, 0, len(allProperties))
	for _, p := range allProperties {
		if len(p.OS) == 0 || slices.Contains(p.OS, goos) {
			result = append(result, p)
		}
	}
	return result
}

// AllowedNames returns the names of all properties accepted by the PATCH handler.
// Text-input properties (quota, mountpoint) and select properties that are editable
// are included. volblocksize is create-only and excluded.
func AllowedNames() []string {
	names := make([]string, 0, len(allProperties))
	for _, p := range allProperties {
		if p.Editable {
			names = append(names, p.Name)
		}
	}
	// acltype, sharenfs, sharesmb are managed via separate API paths but must
	// remain accepted by the PATCH handler for NFS/SMB/ACL wiring.
	names = append(names, "acltype", "sharenfs", "sharesmb")
	return names
}
