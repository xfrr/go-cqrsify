package apix

import (
	"bytes"
	"encoding/json"
	"errors"
	"maps"
)

// ApplyMergePatch implements RFC 7386 (JSON Merge Patch).
// It applies 'patch' to 'original', returning the merged JSON.
//
// If 'patch' is not a JSON object, it replaces 'original' per RFC.
//
// If 'original' is not valid JSON, it is treated as an empty object ({}).
//
// This function does not modify the input slices.
func ApplyMergePatch(original, patch []byte) ([]byte, error) {
	var p any
	if err := json.Unmarshal(patch, &p); err != nil {
		return nil, err
	}

	// If patch is not an object, per RFC return patch as the result
	if _, ok := p.(map[string]any); !ok {
		// return patch normalized JSON (preserve formatting via re-marshal)
		return json.Marshal(p)
	}

	// Ensure target is object (or empty object if not)
	var o any
	if len(bytes.TrimSpace(original)) == 0 {
		o = map[string]any{}
	} else if err := json.Unmarshal(original, &o); err != nil {
		// If original invalid → treat as empty object per practical implementations
		o = map[string]any{}
	}

	tgt, ok := o.(map[string]any)
	if !ok {
		tgt = map[string]any{}
	}

	pmap, _ := p.(map[string]any)
	res, err := mergeObjects(tgt, pmap)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func mergeObjects(target, patch map[string]any) (map[string]any, error) {
	// Copy target to avoid mutating caller's map
	out := make(map[string]any, len(target))
	maps.Copy(out, target)

	for k, pv := range patch {
		if pv == nil {
			// remove key
			delete(out, k)
			continue
		}
		// recurse if both are objects
		if pObj, ok := pv.(map[string]any); ok {
			if tObj, objIsMap := out[k].(map[string]any); objIsMap {
				merged, err := mergeObjects(tObj, pObj)
				if err != nil {
					return nil, err
				}
				out[k] = merged
			} else {
				// replace entirely
				out[k] = pObj
			}
			continue
		}
		// arrays or primitives -> replace
		out[k] = pv
	}
	return out, nil
}

// ApplyMergePatchTo is a typed convenience: it applies 'patch' to 'orig' (T),
// returning the updated T.
//
//	var user User
//	b, _ := json.Marshal(user)
//	out, _ := ApplyMergePatchTo[User](b, patchBytes)
//
// If patch is not compatible with T, you'll get a JSON unmarshal error.
func ApplyMergePatchTo[T any](original []byte, patch []byte) (T, error) {
	var zero T
	merged, err := ApplyMergePatch(original, patch)
	if err != nil {
		return zero, err
	}
	var out T
	if err = json.Unmarshal(merged, &out); err != nil {
		return zero, err
	}
	return out, nil
}

// MustApplyMergePatch is a convenience panic-on-error variant (use carefully).
func MustApplyMergePatch(original, patch []byte) []byte {
	out, err := ApplyMergePatch(original, patch)
	if err != nil {
		panic(err)
	}
	return out
}

// ValidateMergePatch ensures the payload is valid JSON. Optionally enforce object-only.
func ValidateMergePatch(b []byte, requireObject bool) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if requireObject {
		if _, ok := v.(map[string]any); !ok {
			return errors.New("merge-patch must be a JSON object")
		}
	}
	return nil
}
