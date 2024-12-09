package mvcommon

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseNumberRanges(input string, max int) ([]int, error) {
	var indices []int
	ranges := strings.Split(input, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", r)
			}
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || start < 1 || end > max || start > end {
				return nil, fmt.Errorf("invalid range: %s", r)
			}
			for i := start; i <= end; i++ {
				indices = append(indices, i-1)
			}
		} else {
			num, err := strconv.Atoi(r)
			if err != nil || num < 1 || num > max {
				return nil, fmt.Errorf("invalid number: %s", r)
			}
			indices = append(indices, num-1)
		}
	}
	return indices, nil
}
