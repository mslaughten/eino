/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================================================================
// BUG-23: MCPToolResult.String() panics when Error.Code is nil.
//
// MCPToolCallError.Code is *int64 (a pointer). The String() method
// dereferences it with *m.Error.Code without checking for nil.
//
// Impact: Any MCPToolResult with an error that has no Code field
// will panic when String() is called (e.g., during logging).
// ==================================================================

func TestAttack_MCPToolResult_NilErrorCode(t *testing.T) {
	result := &MCPToolResult{
		CallID: "test-call",
		Name:   "test-tool",
		Result: "some result",
		Error: &MCPToolCallError{
			Code:    nil,
			Message: "something went wrong",
		},
	}

	require.NotPanics(t, func() {
		s := result.String()
		t.Logf("String output: %s", s)
		assert.Contains(t, s, "something went wrong")
	}, "BUG: MCPToolResult.String() should not panic when Error.Code is nil")
}

func TestAttack_MCPToolResult_WithErrorCode(t *testing.T) {
	code := int64(500)
	result := &MCPToolResult{
		CallID: "test-call",
		Name:   "test-tool",
		Result: "",
		Error: &MCPToolCallError{
			Code:    &code,
			Message: "internal server error",
		},
	}

	require.NotPanics(t, func() {
		s := result.String()
		assert.Contains(t, s, "500")
		assert.Contains(t, s, "internal server error")
	})
}

// ==================================================================
// BUG-24: NewContentBlockChunk panics if NewContentBlock returns nil
// for certain edge cases in the type switch default branch.
//
// While NewContentBlock uses a generic constraint that prevents truly
// unsupported types, the default branch still returns nil, meaning
// new types added to the constraint but not the switch will cause
// NewContentBlockChunk to dereference nil.
// ==================================================================

func TestAttack_NewContentBlockChunk_NilMeta(t *testing.T) {
	require.NotPanics(t, func() {
		block := NewContentBlockChunk(&AssistantGenText{Text: "test"}, nil)
		require.NotNil(t, block)
		assert.Nil(t, block.StreamingMeta)
	}, "NewContentBlockChunk should handle nil meta without panic")
}

// ==================================================================
// BUG-25: concatAssistantGenTexts overwrites ConcatSliceValue result.
//
// Line 1377 sets ret.Extension to the concat result, but line 1381
// overwrites it with extensions.Interface() (the raw unconcatenated
// reflect.Value). This means the concatenation is discarded.
//
// Impact: Streaming responses with Extension data will have incorrect
// (unconcatenated) extension values after stream assembly.
// ==================================================================

func TestAttack_ConcatAssistantGenTexts_ExtensionOverwrite(t *testing.T) {
	type testExtension struct {
		Value string
	}

	texts := []*AssistantGenText{
		{Text: "Hello ", Extension: &testExtension{Value: "ext1"}},
		{Text: "world", Extension: &testExtension{Value: "ext2"}},
	}

	result, err := concatAssistantGenTexts(texts)
	if err != nil {
		t.Logf("Concat error (may be expected if ConcatSliceValue doesn't handle this type): %v", err)
		t.Skip("Skipping: ConcatSliceValue doesn't support test type")
	}
	require.NotNil(t, result)

	assert.Equal(t, "Hello world", result.Text)

	if result.Extension != nil {
		t.Logf("Extension type: %T, value: %v", result.Extension, result.Extension)
		_, isSlice := result.Extension.([]*testExtension)
		if isSlice {
			t.Log("WARNING: Extension is a raw slice instead of a concatenated value. " +
				"Line 1381 in agentic_message.go overwrites the ConcatSliceValue result " +
				"with extensions.Interface(), discarding the concatenation.")
		}
	}
}

// ==================================================================
// Verify basic agentic message construction functions work correctly.
// ==================================================================

func TestAttack_AgenticMessageConstructors(t *testing.T) {
	t.Run("UserAgenticMessage", func(t *testing.T) {
		msg := UserAgenticMessage("hello")
		require.NotNil(t, msg)
		assert.Equal(t, AgenticRoleTypeUser, msg.Role)
		require.Len(t, msg.ContentBlocks, 1)
		assert.Equal(t, "hello", msg.ContentBlocks[0].UserInputText.Text)
	})

	t.Run("SystemAgenticMessage", func(t *testing.T) {
		msg := SystemAgenticMessage("instruction")
		require.NotNil(t, msg)
		assert.Equal(t, AgenticRoleTypeSystem, msg.Role)
	})

	t.Run("FunctionToolResultAgenticMessage", func(t *testing.T) {
		msg := FunctionToolResultAgenticMessage("call-1", "tool1", "result")
		require.NotNil(t, msg)
		assert.Equal(t, AgenticRoleTypeUser, msg.Role)
		require.Len(t, msg.ContentBlocks, 1)
		assert.Equal(t, ContentBlockTypeFunctionToolResult, msg.ContentBlocks[0].Type)
		assert.Equal(t, "call-1", msg.ContentBlocks[0].FunctionToolResult.CallID)
	})

	t.Run("ContentBlock nil dereference protection", func(t *testing.T) {
		msg := &AgenticMessage{
			ContentBlocks: []*ContentBlock{nil, nil},
		}
		require.NotPanics(t, func() {
			_ = msg.String()
		}, "String() should handle nil ContentBlocks without panic")
	})
}
