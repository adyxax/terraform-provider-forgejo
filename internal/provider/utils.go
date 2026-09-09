package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// boolPointerIfKnown returns a pointer to the bool value, or nil when the value
// is null or unknown so that it is omitted from API request payloads.
//
// It differs from the SDK's types.Bool.ValueBoolPointer() in the unknown case:
// that method returns a pointer to false, which would override the server's
// default.
func boolPointerIfKnown(v types.Bool) *bool {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueBoolPointer()
}

// knownOr returns the planned value when it is known, and the server value
// otherwise. It is meant to be used after an update API call: values that were
// planned are kept as is (avoiding "inconsistent result after apply" errors on
// incomplete API responses), while values that were unknown at plan time are
// filled in from the server response.
func knownOr[T attr.Value](planned T, server T) T {
	if planned.IsUnknown() {
		return server
	}
	return planned
}

// stringPointerIfKnown returns a pointer to the string value, or nil when the
// value is null or unknown so that it is omitted from API request payloads.
//
// It differs from the SDK's types.String.ValueStringPointer() in the unknown
// case: that method returns a pointer to "", which would override the server's
// default.
func stringPointerIfKnown(v types.String) *string {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueStringPointer()
}
