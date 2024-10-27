package helper

// StringOrNA handle potentially nil strings
func StringOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
