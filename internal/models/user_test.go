package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "30 years old - birthday passed",
			dob:      time.Date(1994, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(1994, 1, 15),
		},
		{
			name:     "25 years old - birthday today",
			dob:      time.Now().AddDate(-25, 0, 0),
			expected: 25,
		},
		{
			name:     "20 years old - birthday tomorrow",
			dob:      time.Now().AddDate(-20, 0, 1),
			expected: 19,
		},
		{
			name:     "Newborn",
			dob:      time.Now(),
			expected: 0,
		},
		{
			name:     "1 year old",
			dob:      time.Now().AddDate(-1, 0, 0),
			expected: 1,
		},
		{
			name:     "Leap year baby - Feb 29",
			dob:      time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(2000, 2, 29),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.expected {
				t.Errorf("CalculateAge() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// calculateExpectedAge calculates the expected age for testing
func calculateExpectedAge(year, month, day int) int {
	now := time.Now()
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	years := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	return years
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name      string
		dateStr   string
		expectErr bool
		expected  time.Time
	}{
		{
			name:      "Valid date",
			dateStr:   "1990-05-10",
			expectErr: false,
			expected:  time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Invalid format - wrong separator",
			dateStr:   "1990/05/10",
			expectErr: true,
		},
		{
			name:      "Invalid format - US format",
			dateStr:   "05-10-1990",
			expectErr: true,
		},
		{
			name:      "Invalid date - month 13",
			dateStr:   "1990-13-10",
			expectErr: true,
		},
		{
			name:      "Invalid date - day 32",
			dateStr:   "1990-05-32",
			expectErr: true,
		},
		{
			name:      "Empty string",
			dateStr:   "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.dateStr)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseDate() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseDate() unexpected error: %v", err)
				return
			}
			if !got.Equal(tt.expected) {
				t.Errorf("ParseDate() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "Standard date",
			input:    time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			expected: "1990-05-10",
		},
		{
			name:     "Single digit month and day",
			input:    time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
			expected: "2000-01-05",
		},
		{
			name:     "End of year",
			input:    time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: "2023-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDate(tt.input)
			if got != tt.expected {
				t.Errorf("FormatDate() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestAgeEdgeCases(t *testing.T) {
	now := time.Now()

	// Test: Birthday is exactly today
	t.Run("Birthday today", func(t *testing.T) {
		dob := time.Date(now.Year()-30, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		age := CalculateAge(dob)
		if age != 30 {
			t.Errorf("Expected age 30 on birthday, got %d", age)
		}
	})

	// Test: Birthday was yesterday
	t.Run("Birthday yesterday", func(t *testing.T) {
		yesterday := now.AddDate(0, 0, -1)
		dob := time.Date(now.Year()-30, yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
		age := CalculateAge(dob)
		if age != 30 {
			t.Errorf("Expected age 30 day after birthday, got %d", age)
		}
	})

	// Test: Birthday is tomorrow
	t.Run("Birthday tomorrow", func(t *testing.T) {
		tomorrow := now.AddDate(0, 0, 1)
		dob := time.Date(now.Year()-30, tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
		age := CalculateAge(dob)
		if age != 29 {
			t.Errorf("Expected age 29 day before birthday, got %d", age)
		}
	})
}
