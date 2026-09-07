package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func knownOr[T attr.Value](planned T, server T) T {
	if planned.IsUnknown() {
		return server
	}
	return planned
}

func boolPointerIfKnown(v types.Bool) *bool {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueBoolPointer()
}

func stringPointerIfKnown(v types.String) *string {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueStringPointer()
}
