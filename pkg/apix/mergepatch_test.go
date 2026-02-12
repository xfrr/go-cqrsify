package apix_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	apix "github.com/xfrr/go-cqrsify/pkg/apix"
)

func asMap(s *suite.Suite, b []byte) map[string]any {
	var m map[string]any
	s.Require().NoError(json.Unmarshal(b, &m))
	return m
}

func asAny(s *suite.Suite, b []byte) any {
	var v any
	s.Require().NoError(json.Unmarshal(b, &v))
	return v
}

type User struct {
	ID      string         `json:"id,omitempty"`
	Name    string         `json:"name,omitempty"`
	Age     int            `json:"age,omitempty"`
	Address map[string]any `json:"address,omitempty"`
	Tags    []string       `json:"tags,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type ApplyMergePatchSuite struct {
	suite.Suite
}

func (s *ApplyMergePatchSuite) Test_PatchIsNotObject_ReplacesOriginal() {
	orig := []byte(`{"a":1}`)
	tests := []struct {
		name  string
		patch []byte
		want  any
	}{
		{"Number", []byte(`5`), float64(5)},
		{"String", []byte(`"x"`), "x"},
		{"Bool", []byte(`true`), true},
		{"Null", []byte(`null`), nil},
		{"Array", []byte(`[1,2,3]`), []any{float64(1), float64(2), float64(3)}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			out, err := apix.ApplyMergePatch(orig, tt.patch)
			s.Require().NoError(err)
			got := asAny(&s.Suite, out)
			s.Equal(tt.want, got)
		})
	}
}

func (s *ApplyMergePatchSuite) Test_InvalidPatch_ReturnsError() {
	orig := []byte(`{"a":1}`)
	_, err := apix.ApplyMergePatch(orig, []byte(`{`))
	s.Error(err)
}

func (s *ApplyMergePatchSuite) Test_InvalidOrEmptyOriginal_TreatedAsEmptyObject() {
	tests := []struct {
		name  string
		orig  []byte
		patch []byte
		want  map[string]any
	}{
		{"EmptyBytes", []byte(``), []byte(`{"a":2}`), map[string]any{"a": float64(2)}},
		{"WhitespaceOnly", []byte(`   `), []byte(`{"a":2}`), map[string]any{"a": float64(2)}},
		{"InvalidJSON", []byte(`{`), []byte(`{"a":2}`), map[string]any{"a": float64(2)}},
		{"OriginalNotObject(Array)", []byte(`[1,2]`), []byte(`{"a":2}`), map[string]any{"a": float64(2)}},
		{"OriginalNotObject(Primitive)", []byte(`"x"`), []byte(`{"a":2}`), map[string]any{"a": float64(2)}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			out, err := apix.ApplyMergePatch(tt.orig, tt.patch)
			s.Require().NoError(err)
			s.Equal(tt.want, asMap(&s.Suite, out))
		})
	}
}

func (s *ApplyMergePatchSuite) Test_ObjectPatch_BasicMergeReplaceAndDelete() {
	orig := []byte(`{"a":1, "b":2, "c":{"n":1,"m":2}, "d":[1,2,3], "e":"keep"}`)
	patch := []byte(`{
		"a": 10,
		"b": null, 
		"c": {"m": 20},
		"d": [9,9],   
		"f": "new"        
	}`)
	out, err := apix.ApplyMergePatch(orig, patch)
	s.Require().NoError(err)

	got := asMap(&s.Suite, out)
	// a replaced
	s.InEpsilon(float64(10), got["a"], 0.0001)
	// b removed
	_, exists := got["b"]
	s.False(exists)
	// c merged
	c, ok := got["c"].(map[string]any)
	s.Require().True(ok)
	s.Require().NotNil(c)
	s.InEpsilon(float64(1), c["n"], 0.0001)
	s.InEpsilon(float64(20), c["m"], 0.0001)
	// d replaced entirely
	s.Equal([]any{float64(9), float64(9)}, got["d"])
	// f added
	s.Equal("new", got["f"])
	// e untouched
	s.Equal("keep", got["e"])
}

func (s *ApplyMergePatchSuite) Test_ObjectPatch_ReplacingNonMapWithMap() {
	orig := []byte(`{"x": 1}`)
	patch := []byte(`{"x": {"a":1}}`)
	out, err := apix.ApplyMergePatch(orig, patch)
	s.Require().NoError(err)
	got := asMap(&s.Suite, out)
	s.Equal(map[string]any{"a": float64(1)}, got["x"])
}

func (s *ApplyMergePatchSuite) Test_ObjectPatch_DeepMerge() {
	orig := []byte(`{"a":{"b":{"c":1,"keep":true}}}`)
	patch := []byte(`{"a":{"b":{"c":2,"d":3}}}`)
	out, err := apix.ApplyMergePatch(orig, patch)
	s.Require().NoError(err)
	a, ok := asMap(&s.Suite, out)["a"].(map[string]any)
	s.Require().True(ok)
	s.Require().NotNil(a)
	b, ok := a["b"].(map[string]any)
	s.Require().True(ok)
	s.Require().NotNil(b)
	s.InEpsilon(float64(2), b["c"], 0.0001)
	s.InEpsilon(float64(3), b["d"], 0.0001)
	s.Equal(true, b["keep"]) // preserved
}

func TestApplyMergePatchSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ApplyMergePatchSuite))
}

type ApplyMergePatchToSuite struct {
	suite.Suite
}

func (s *ApplyMergePatchToSuite) Test_SuccessfullyAppliesAndUnmarshals() {
	orig := User{ID: "u1", Name: "Alice", Age: 30, Address: map[string]any{"city": "Málaga"}, Tags: []string{"a"}}
	ob, err := json.Marshal(orig)
	s.Require().NoError(err)

	patch := []byte(`{
		"name": "Alicia",
		"age": 31,
		"address": { "country": "ES" },
		"tags": ["x","y"],
		"meta": {"tier":"gold"}
	}`)

	out, err := apix.ApplyMergePatchTo[User](ob, patch)
	s.Require().NoError(err)

	s.Equal("u1", out.ID)
	s.Equal("Alicia", out.Name)
	s.Equal(31, out.Age)
	s.Equal([]string{"x", "y"}, out.Tags)
	s.Equal("Málaga", out.Address["city"])
	s.Equal("ES", out.Address["country"])
	s.Equal("gold", out.Meta["tier"])
}

func (s *ApplyMergePatchToSuite) Test_UnmarshalError_WhenTypesMismatch() {
	// Patch sets "age" to string, but struct expects int -> JSON unmarshal should fail.
	orig := User{ID: "u1", Age: 30}
	ob, err := json.Marshal(orig)
	s.Require().NoError(err)

	_, err = apix.ApplyMergePatchTo[User](ob, []byte(`{"age":"oops"}`))
	s.Error(err)
}

func (s *ApplyMergePatchToSuite) Test_UsesEmptyObjectWhenOriginalInvalid() {
	ob := []byte(`{`) // invalid original, treated as {}
	out, err := apix.ApplyMergePatchTo[User](ob, []byte(`{"name":"Neo"}`))
	s.Require().NoError(err)
	s.Equal("Neo", out.Name)
}

func TestApplyMergePatchToSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ApplyMergePatchToSuite))
}

type MustApplyMergePatchSuite struct {
	suite.Suite
}

func (s *MustApplyMergePatchSuite) Test_NoPanic_OnValidInputs() {
	out := apix.MustApplyMergePatch([]byte(`{"a":1}`), []byte(`{"a":2}`))
	m := asMap(&s.Suite, out)
	s.InEpsilon(float64(2), m["a"], 0.0001)
}

func (s *MustApplyMergePatchSuite) Test_Panics_OnInvalidPatch() {
	s.Panics(func() {
		_ = apix.MustApplyMergePatch([]byte(`{"a":1}`), []byte(`{`))
	})
}

func TestMustApplyMergePatchSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MustApplyMergePatchSuite))
}

type ValidateMergePatchSuite struct {
	suite.Suite
}

func (s *ValidateMergePatchSuite) Test_ValidJSON_NoObjectRequirement() {
	s.NoError(apix.ValidateMergePatch([]byte(`123`), false))
	s.NoError(apix.ValidateMergePatch([]byte(`"ok"`), false))
	s.NoError(apix.ValidateMergePatch([]byte(`[1,2]`), false))
	s.NoError(apix.ValidateMergePatch([]byte(`{"a":1}`), false))
}

func (s *ValidateMergePatchSuite) Test_InvalidJSON_ReturnsError() {
	s.Error(apix.ValidateMergePatch([]byte(`{`), false))
}

func (s *ValidateMergePatchSuite) Test_RequireObject_RejectsNonObject() {
	s.Require().NoError(apix.ValidateMergePatch([]byte(`{"a":1}`), true))
	s.Require().Error(apix.ValidateMergePatch([]byte(`123`), true))
	s.Require().Error(apix.ValidateMergePatch([]byte(`"x"`), true))
	s.Require().Error(apix.ValidateMergePatch([]byte(`[1,2]`), true))
	s.Require().Error(apix.ValidateMergePatch([]byte(`null`), true))
}

func TestValidateMergePatchSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateMergePatchSuite))
}
