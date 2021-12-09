package server

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/goto-basic-object.jsonnet
	basicJsonnetContent string
)

func getVM() (vm *jsonnet.VM) {
	vm = jsonnet.MakeVM()
	vm.Importer(&jsonnet.MemoryImporter{
		Data: map[string]jsonnet.Contents{
			"goto-basic-object.jsonnet": jsonnet.MakeContents(basicJsonnetContent),
		},
	})
	return
}

func TestDefinition(t *testing.T) {
	testCases := []struct {
		name     string
		params   protocol.DefinitionParams
		expected *protocol.DefinitionLink
	}{
		{
			name: "test goto definition for var myvar",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_goto_definition.jsonnet",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 19,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/test_goto_definition.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      0,
						Character: 15,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      0,
						Character: 11,
					},
				},
			},
		},
		{
			name: "test goto definition on function helper",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_goto_definition.jsonnet",
					},
					Position: protocol.Position{
						Line:      7,
						Character: 8,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/test_goto_definition.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 6,
					},
					End: protocol.Position{
						Line:      1,
						Character: 23,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 6,
					},
					End: protocol.Position{
						Line:      1,
						Character: 12,
					},
				},
			},
		},
		{
			name: "test goto inner definition",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_goto_definition_multi_locals.jsonnet",
					},
					Position: protocol.Position{
						Line:      6,
						Character: 11,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/test_goto_definition_multi_locals.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      4,
						Character: 10,
					},
					End: protocol.Position{
						Line:      4,
						Character: 28,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      4,
						Character: 10,
					},
					End: protocol.Position{
						Line:      4,
						Character: 18,
					},
				},
			},
		},
		{
			name: "test goto super index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_combined_object.jsonnet",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 13,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/test_combined_object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 5,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      1,
						Character: 5,
					},
				},
			},
		},
		{
			name: "test goto super nested",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_combined_object.jsonnet",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 15,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/test_combined_object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 22,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 9,
					},
				},
			},
		},
		{
			name: "test goto self object field function",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_basic_lib.libsonnet",
					},
					Position: protocol.Position{
						Line:      4,
						Character: 19,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/test_basic_lib.libsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 20,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      1,
						Character: 9,
					},
				},
			},
		},
		{
			name: "test goto super object field local defined obj 'foo'",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/oo-contrived.jsonnet",
					},
					Position: protocol.Position{
						Line:      12,
						Character: 17,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/oo-contrived.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 2,
					},
					End: protocol.Position{
						Line:      1,
						Character: 8,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 2,
					},
					End: protocol.Position{
						Line:      1,
						Character: 5,
					},
				},
			},
		},
		{
			name: "test goto super object field local defined obj 'g'",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/oo-contrived.jsonnet",
					},
					Position: protocol.Position{
						Line:      13,
						Character: 17,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/oo-contrived.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 2,
					},
					End: protocol.Position{
						Line:      2,
						Character: 19,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 2,
					},
					End: protocol.Position{
						Line:      2,
						Character: 3,
					},
				},
			},
		},
		{
			name: "test goto local var from other local var",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/oo-contrived.jsonnet",
					},
					Position: protocol.Position{
						Line:      6,
						Character: 9,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/oo-contrived.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      3,
						Character: 1,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      0,
						Character: 10,
					},
				},
			},
		},
		{
			name: "test goto local obj field from 'self.attr' from other obj",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-indexes.jsonnet",
					},
					Position: protocol.Position{
						Line:      9,
						Character: 17,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/goto-indexes.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 23,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 11,
					},
				},
			},
		},
		{
			name: "test goto local object 'obj' via obj index 'obj.foo'",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-indexes.jsonnet",
					},
					Position: protocol.Position{
						Line:      8,
						Character: 15,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "./testdata/goto-indexes.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 5,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      1,
						Character: 7,
					},
				},
			},
		},
		{
			name: "test goto imported file",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 22,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "goto-basic-object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 0,
					},
					End: protocol.Position{
						Line:      0,
						Character: 0,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 0,
					},
					End: protocol.Position{
						Line:      0,
						Character: 0,
					},
				},
			},
		},
		{
			name: "test goto imported file at lhs index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      3,
						Character: 18,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "goto-basic-object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 7,
					},
				},
			},
		},
		{
			name: "test goto imported file at rhs index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      4,
						Character: 18,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "goto-basic-object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto import index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-import-attribute.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 48,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "goto-basic-object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto attribute of nested import",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-nested-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      2,
						Character: 15,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: "goto-basic-object.jsonnet",
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 7,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := string(tc.params.TextDocument.URI)
			var content, err = os.ReadFile(filename)
			require.NoError(t, err)
			ast, err := jsonnet.SnippetToAST(filename, string(content))
			require.NoError(t, err)
			got, err := Definition(ast, &tc.params, getVM())
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDefinitionFail(t *testing.T) {
	testCases := []struct {
		name     string
		params   protocol.DefinitionParams
		expected error
	}{
		{
			name: "goto local keyword fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-basic-object.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 3,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("cannot find definition"),
		},
		{
			name: "goto function argument from inside function fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-functions.libsonnet",
					},
					Position: protocol.Position{
						Line:      7,
						Character: 13,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("unable to find matching bind for arg1"),
		},
		{
			name: "goto index of std fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-std.jsonnet",
					},
					Position: protocol.Position{
						Line:      1,
						Character: 20,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("cannot get definition of std lib"),
		},
		{
			name: "goto comment fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-comment.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 1,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("cannot find definition"),
		},
		{
			name: "goto local func param fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-local-function.libsonnet",
					},
					Position: protocol.Position{
						Line:      2,
						Character: 25,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("unable to find matching bind for k"),
		},
		{
			name: "goto range index fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-local-function.libsonnet",
					},
					Position: protocol.Position{
						Line:      15,
						Character: 57,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("unexpected node type when finding bind for 'ports'"),
		},
		{
			name: "goto super fails as no LHS object exists",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-local-function.libsonnet",
					},
					Position: protocol.Position{
						Line:      33,
						Character: 23,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: fmt.Errorf("could not find a lhs object"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := string(tc.params.TextDocument.URI)
			var content, err = os.ReadFile(filename)
			require.NoError(t, err)
			ast, err := jsonnet.SnippetToAST(filename, string(content))
			require.NoError(t, err)
			got, err := Definition(ast, &tc.params, getVM())
			require.Error(t, err)
			assert.Equal(t, tc.expected.Error(), err.Error())
			assert.Nil(t, got)
		})
	}
}
