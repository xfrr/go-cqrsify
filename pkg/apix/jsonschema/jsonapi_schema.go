package jsonschema

var jsonapiSchema = []byte(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/schemas/jsonapi-data-document.schema.json",
  "title": "JSON:API Data Document (Single or Many)",
  "$comment": "A JSON:API document as per https://jsonapi.org/format/1.1/#document-top-level",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "links": { "$ref": "#/$defs/Links" },
    "included": {
      "type": "array",
      "items": { "$ref": "#/$defs/ResourceIncluded" }
    },
    "meta": { "$ref": "#/$defs/Meta" },
    "data": {
      "oneOf": [
        { "$ref": "#/$defs/Resource" },
        {
          "type": "array",
          "items": { "$ref": "#/$defs/Resource" }
        }
      ]
    }
  },
  "required": ["data"],

  "$defs": {
    "Links": {
      "title": "Links",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "self": { "type": "string", "format": "uri" },
        "next": { "type": "string", "format": "uri" },
        "prev": { "type": "string", "format": "uri" }
      }
    },

    "Meta": {
      "title": "Meta",
      "type": "object",
      "additionalProperties": true
    },

    "ResourceIdentifier": {
      "title": "Resource Identifier",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "type": { "type": "string", "minLength": 1 },
        "id":   { "type": "string", "minLength": 1 }
      },
      "required": ["type", "id"]
    },

    "Relationship": {
      "title": "Relationship",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "links": { "$ref": "#/$defs/Links" },
        "meta":  { "$ref": "#/$defs/Meta" },
        "data": {
          "oneOf": [
            { "$ref": "#/$defs/ResourceIdentifier" },
            {
              "type": "array",
              "items": { "$ref": "#/$defs/ResourceIdentifier" }
            },
            { "type": "null" }
          ]
        }
      }
    },

    "RelationshipsMap": {
      "title": "Relationships Map",
      "type": "object",
      "propertyNames": { "type": "string", "minLength": 1 },
      "patternProperties": {
        "^.+$": { "$ref": "#/$defs/Relationship" }
      },
      "additionalProperties": false
    },

    "Attributes": {
      "title": "Attributes (replace with your concrete schema if desired)",
      "type": "object",
      "additionalProperties": true
    },

    "Resource": {
      "title": "Resource",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "type": { "type": "string", "minLength": 1 },
        "id":   { "type": "string", "minLength": 1 },
        "attributes": { "$ref": "#/$defs/Attributes" },
        "relationships": { "$ref": "#/$defs/RelationshipsMap" },
        "links": { "$ref": "#/$defs/Links" },
        "meta":  { "$ref": "#/$defs/Meta" }
      },
      "required": ["type", "attributes"]
    },

    "ResourceIncluded": {
      "title": "Included Resource",
      "$comment": "Included resources follow the same structure as Resource. Attributes remain generic unless you specialize them.",
      "allOf": [{ "$ref": "#/$defs/Resource" }]
    }
  }
}`)
