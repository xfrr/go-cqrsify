package messaginghttp_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/xfrr/go-cqrsify/messaging"
	messaginghttp "github.com/xfrr/go-cqrsify/messaging/http"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

type CreateBaseCommandFromSingleDocumentTestSuite struct {
	suite.Suite
}

func TestCreateBaseCommandFromSingleDocumentTestSuite(t *testing.T) {
	suite.Run(t, new(CreateBaseCommandFromSingleDocumentTestSuite))
}

func (st *CreateBaseCommandFromSingleDocumentTestSuite) Test_CreateBaseCommandFromSingleDocument() {
	type args struct {
		cmdType string
		sd      apix.SingleDocument[any]
	}
	tests := []struct {
		name string
		args args
		want messaging.BaseCommand
	}{
		{
			name: "BaseCommand with all valid fields",
			args: args{
				cmdType: "createUser",
				sd: apix.SingleDocument[any]{
					Data: apix.Resource[any]{
						Type: "createUser",
						ID:   "123e4567-e89b-12d3-a456-426614174000",
						Attributes: map[string]any{
							"name": "John Doe",
						},
						Meta: map[string]any{
							"schema":    "testSchema",
							"source":    "testSource",
							"timestamp": "2024-01-01T12:00:00Z",
							"extra":     "metadata",
						},
					},
					Meta: map[string]any{},
				},
			},
			want: messaging.NewBaseCommand("createUser",
				messaging.WithID("123e4567-e89b-12d3-a456-426614174000"),
				messaging.WithSchema("testSchema"),
				messaging.WithSource("testSource"),
				messaging.WithTimestamp(func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
					return t
				}()),
				messaging.WithMetadata(map[string]string{
					"extra": "metadata",
				}),
			),
		},
		{
			name: "BaseCommand with duplicated meta fields in resource and document",
			args: args{
				cmdType: "updateUser",
				sd: apix.SingleDocument[any]{
					Data: apix.Resource[any]{
						Type: "updateUser",
						ID:   "223e4567-e89b-12d3-a456-426614174001",
						Attributes: map[string]any{
							"name": "Jane Doe",
						},
						Meta: map[string]any{
							"schema":    "resourceSchema",
							"source":    "resourceSource",
							"timestamp": "2024-02-01T12:00:00Z",
							"role":      "admin",
						},
					},
					Meta: map[string]any{
						"role": "superadmin", // This should take precedence
					},
				},
			},
			want: messaging.NewBaseCommand("updateUser",
				messaging.WithID("223e4567-e89b-12d3-a456-426614174001"),
				messaging.WithSchema("resourceSchema"),
				messaging.WithSource("resourceSource"),
				messaging.WithTimestamp(func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-01T12:00:00Z")
					return t
				}()),
				messaging.WithMetadata(map[string]string{
					"role": "superadmin",
				}),
			),
		},
		{
			name: "BaseCommand with not string meta values",
			args: args{
				cmdType: "deleteUser",
				sd: apix.SingleDocument[any]{
					Data: apix.Resource[any]{
						Type: "deleteUser",
						ID:   "323e4567-e89b-12d3-a456-426614174002",
						Attributes: map[string]any{
							"confirm": true,
						},
						Meta: map[string]any{
							"int":     42,
							"float":   3.14,
							"bool":    false,
							"int8":    int8(8),
							"int16":   int16(16),
							"int32":   int32(32),
							"int64":   int64(64),
							"uint":    uint(1),
							"uint8":   uint8(8),
							"uint16":  uint16(16),
							"uint32":  uint32(32),
							"uint64":  uint64(64),
							"ignored": strconv.ErrSyntax, // Will be ignored
						},
					},
					Meta: map[string]any{},
				},
			},
			want: messaging.NewBaseCommand("deleteUser",
				messaging.WithID("323e4567-e89b-12d3-a456-426614174002"),
				messaging.WithTimestamp(time.Now()), // Approximate, will be checked in test
				messaging.WithMetadata(map[string]string{
					"int":    "42",
					"float":  "3.140000",
					"bool":   "false",
					"int8":   "8",
					"int16":  "16",
					"int32":  "32",
					"int64":  "64",
					"uint":   "1",
					"uint8":  "8",
					"uint16": "16",
					"uint32": "32",
					"uint64": "64",
				}),
			),
		},
		{
			name: "BaseCommand with empty meta, no ID, no schema, no source, no timestamp",
			args: args{
				cmdType: "simpleCommand",
				sd: apix.SingleDocument[any]{
					Data: apix.Resource[any]{
						Type:       "simpleCommand",
						Attributes: map[string]any{},
						Meta:       map[string]any{},
					},
					Meta: map[string]any{},
				},
			},
			want: messaging.NewBaseCommand("simpleCommand"),
		},
	}
	for _, tt := range tests {
		st.Run(tt.name, func() {
			got := messaginghttp.CreateBaseCommandFromSingleDocument(tt.args.cmdType, tt.args.sd)
			st.Require().Equal(tt.want.CommandID(), got.CommandID())
			st.Require().Equal(tt.want.MessageType(), got.MessageType())
			st.Require().Equal(tt.want.MessageSchemaURI(), got.MessageSchemaURI())
			st.Require().Equal(tt.want.MessageSource(), got.MessageSource())
			st.Require().Equal(tt.want.MessageMetadata(), got.MessageMetadata())
			st.Require().WithinDuration(tt.want.MessageTimestamp(), got.MessageTimestamp(), time.Millisecond)
		})
	}
}
