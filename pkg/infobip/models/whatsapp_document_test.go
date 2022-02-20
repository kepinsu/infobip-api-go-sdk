package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidDocumentMessage(t *testing.T) {
	tests := []struct {
		name     string
		instance DocumentMsg
	}{
		{name: "minimum input",
			instance: DocumentMsg{
				MsgCommon: MsgCommon{
					From: "16175551213",
					To:   "16175551212",
				},
				Content: DocumentContent{MediaURL: "https://www.mypath.com/my_doc.txt"},
			}},
		{
			name: "complete input",
			instance: DocumentMsg{
				MsgCommon: GenerateTestMsgCommon(),
				Content: DocumentContent{
					MediaURL: "https://www.mypath.com/my_doc.txt",
					Caption:  "hello world",
					Filename: "my_doc.txt",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.instance.Validate()
			require.Nil(t, err)
		})
	}
}

func TestDocumentMessageConstraints(t *testing.T) {
	msgCommon := GenerateTestMsgCommon()
	tests := []struct {
		name    string
		content DocumentContent
	}{
		{
			name:    "empty Content field",
			content: DocumentContent{},
		},
		{
			name:    "missing Content MediaURL",
			content: DocumentContent{Filename: "asd"},
		},
		{
			name:    "Content MediaURL too long",
			content: DocumentContent{MediaURL: fmt.Sprintf("https://www.g%sgle.com", strings.Repeat("o", 2040))},
		},
		{
			name:    "Content invalid MediaURL",
			content: DocumentContent{MediaURL: "asd"},
		},
		{
			name: "Content Caption too long",
			content: DocumentContent{
				MediaURL: "https://www.mypath.com/my_doc.txt",
				Caption:  strings.Repeat("a", 3001),
			},
		},
		{
			name: "Content Filename too long",
			content: DocumentContent{
				MediaURL: "https://www.mypath.com/my_doc.txt",
				Filename: strings.Repeat("a", 241),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := DocumentMsg{
				MsgCommon: msgCommon,
				Content:   tc.content,
			}
			err := msg.Validate()
			require.NotNil(t, err)
		})
	}
}
