package pkg

import "strings"

func SanitizeKeyword(keyword string) string {
	keyword = strings.ReplaceAll(keyword, "'", "''")  // handle single quote
	keyword = strings.ReplaceAll(keyword, "&", "\\&") // escape special FTS characters
	keyword = strings.ReplaceAll(keyword, "|", "\\|")
	keyword = strings.ReplaceAll(keyword, "!", "\\!")
	keyword = strings.ReplaceAll(keyword, "(", "\\(")
	keyword = strings.ReplaceAll(keyword, ")", "\\)")
	keyword = strings.ReplaceAll(keyword, ":", "\\:")
	keyword = strings.ReplaceAll(keyword, "*", "\\*")
	keyword = strings.ReplaceAll(keyword, "<", "\\<")
	keyword = strings.ReplaceAll(keyword, ">", "\\>")
	return keyword
}

func FormatKeywords(keyword string) string {
	keywords := strings.Split(keyword, " ")
	for i, keyword := range keywords {
		keyword = SanitizeKeyword(keyword)
		keywords[i] = keyword + ":*"
	}
	return strings.Join(keywords, " | ")
}
