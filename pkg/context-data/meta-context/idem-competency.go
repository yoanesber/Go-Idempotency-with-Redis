package metacontext

import (
	"context"
)

// This struct defines the IdemCompetencyMeta struct
//
//	It can be used to store metadata about the idem competency information
type IdemCompetencyMeta struct {
	Key             string
	BodyHash        string
	ResponsePayload string
	StatusCode      int
}

// This struct defines the IdemCompetencyMetaKeyType struct
//
//	It is used as a key for storing and retrieving IdemCompetencyMeta from the context
type IdemCompetencyMetaKeyType struct{}

// Define a key for storing IdemCompetencyMeta in the context
var idemCompetencyMetaKey = IdemCompetencyMetaKeyType{}

// InjectIdemCompetencyMeta injects the IdemCompetencyMeta into the context.
// This function is used to add metadata to the context for later retrieval
func InjectIdemCompetencyMeta(ctx context.Context, meta IdemCompetencyMeta) context.Context {
	return context.WithValue(ctx, idemCompetencyMetaKey, meta)
}

// ExtractIdemCompetencyMeta retrieves the IdemCompetencyMeta from the context.
// This function is used to access the metadata stored in the context
func ExtractIdemCompetencyMeta(ctx context.Context) (IdemCompetencyMeta, bool) {
	meta, ok := ctx.Value(idemCompetencyMetaKey).(IdemCompetencyMeta)
	return meta, ok
}
