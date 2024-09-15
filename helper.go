package main

// stringOrNA handle potentially nil strings
func stringOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
