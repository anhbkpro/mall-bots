package domain

import "testing"

func TestBasketStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   BasketStatus
		expected string
	}{
		{
			name:     "Open status",
			status:   BasketStatusOpen,
			expected: "open",
		},
		{
			name:     "Canceled status",
			status:   BasketStatusCanceled,
			expected: "canceled",
		},
		{
			name:     "Checked out status",
			status:   BasketStatusCheckedOut,
			expected: "checked_out",
		},
		{
			name:     "Unknown status",
			status:   BasketStatusUnknown,
			expected: "",
		},
		{
			name:     "Invalid status",
			status:   "invalid",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("BasketStatus.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestToBasketStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BasketStatus
	}{
		{
			name:     "Open status",
			input:    "open",
			expected: BasketStatusOpen,
		},
		{
			name:     "Canceled status",
			input:    "canceled",
			expected: BasketStatusCanceled,
		},
		{
			name:     "Checked out status",
			input:    "checked_out",
			expected: BasketStatusCheckedOut,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: BasketStatusUnknown,
		},
		{
			name:     "Invalid status",
			input:    "invalid",
			expected: BasketStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBasketStatus(tt.input); got != tt.expected {
				t.Errorf("ToBasketStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}
