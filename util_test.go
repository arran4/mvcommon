package mvcommon

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseNumberRanges(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		max     int
		want    []int
		wantErr error
	}{
		{
			name:  "Valid single number",
			input: "3",
			max:   5,
			want:  []int{2},
		},
		{
			name:  "Valid single range",
			input: "1-3",
			max:   5,
			want:  []int{0, 1, 2},
		},
		{
			name:  "Valid mixed range and numbers",
			input: "1-2,4,5",
			max:   5,
			want:  []int{0, 1, 3, 4},
		},
		{
			name:  "Valid multiple ranges",
			input: "1-2,3-5",
			max:   5,
			want:  []int{0, 1, 2, 3, 4},
		},
		{
			name:  "Valid single range with max boundary",
			input: "4-5",
			max:   5,
			want:  []int{3, 4},
		},
		{
			name:    "Invalid range start > end",
			input:   "3-1",
			max:     5,
			wantErr: errors.New("invalid range: 3-1"),
		},
		{
			name:    "Invalid number exceeding max",
			input:   "6",
			max:     5,
			wantErr: errors.New("invalid number: 6"),
		},
		{
			name:    "Invalid range exceeding max",
			input:   "4-6",
			max:     5,
			wantErr: errors.New("invalid range: 4-6"),
		},
		{
			name:    "Invalid non-numeric input",
			input:   "a",
			max:     5,
			wantErr: errors.New("invalid number: a"),
		},
		{
			name:    "Invalid range with non-numeric input",
			input:   "1-a",
			max:     5,
			wantErr: errors.New("invalid range: 1-a"),
		},
		{
			name:    "Empty input",
			input:   "",
			max:     5,
			wantErr: errors.New("invalid number: "),
		},
		{
			name:  "Valid input with spaces",
			input: " 1-2 , 4 , 5 ",
			max:   5,
			want:  []int{0, 1, 3, 4},
		},
		{
			name:  "Valid input with overlapping ranges",
			input: "1-3,2-4",
			max:   5,
			want:  []int{0, 1, 2, 1, 2, 3}, // Overlaps preserved
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNumberRanges(tt.input, tt.max)

			// If we expect an error, ensure it matches
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("ParseNumberRanges() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			// If no error is expected, ensure the result matches
			if err != nil {
				t.Errorf("ParseNumberRanges() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNumberRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
