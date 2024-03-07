package customtypes_test

import (
	"context"
	"testing"

	"github.com/ansible/terraform-provider-aap/internal/provider/customtypes"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestCustomStringStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentData   customtypes.CustomStringValue
		givenData     basetypes.StringValuable
		expectedMatch bool
		expectedDiags diag.Diagnostics
	}{
		"not equal - mismatched field values": {
			currentData:   customtypes.NewCustomStringValue(`"{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}]}"`),
			givenData:     customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":1}]}`),
			expectedMatch: false,
		},
		"not equal - mismatched field names": {
			currentData:   customtypes.NewCustomStringValue(`{"exampleVariables":[{"Name":"bar","namespace":"bar-namespace","type":0}]}`),
			givenData:     customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}]}`),
			expectedMatch: false,
		},
		"not equal - additional field": {
			currentData:   customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}],"new-field": null}`),
			givenData:     customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}]}`),
			expectedMatch: false,
		},
		"not equal - array item order difference": {
			currentData:   customtypes.NewCustomStringValue(`{"exampleVariables":[{"namespace":"bar-namespace","name":"bar","type":0}]}`),
			givenData:     customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}]}`),
			expectedMatch: false,
		},
		"semantically equal - object byte-for-byte match": {
			currentData:   customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}]}`),
			givenData:     customtypes.NewCustomStringValue(`{"exampleVariables":[{"name":"bar","namespace":"bar-namespace","type":0}]}`),
			expectedMatch: true,
		},
		"semantically equal - object whitespace difference": {
			currentData: customtypes.NewCustomStringValue(`{
				"hello": "world",
				"nums": [1, 2, 3],
				"nested": {
					"test-bool": true
				}
			}`),
			givenData:     customtypes.NewCustomStringValue(`{"hello":"world","nums":[1,2,3],"nested":{"test-bool":true}}`),
			expectedMatch: false,
		},
		"semantically equal - yaml no difference": {
			currentData: customtypes.NewCustomStringValue(`os: Linux
			automation: ansible-devel`),
			givenData: customtypes.NewCustomStringValue(`os: Linux
			automation: ansible-devel`),
			expectedMatch: true,
		},
		"semantically equal - yaml no difference with newline": {
			currentData: customtypes.NewCustomStringValue(`os: Linux
			automation: ansible-devel`),
			givenData: customtypes.NewCustomStringValue(`os: Linux
			automation: ansible-devel

			`),
			expectedMatch: true,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentData.StringSemanticEquals(context.Background(), testCase.givenData)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
