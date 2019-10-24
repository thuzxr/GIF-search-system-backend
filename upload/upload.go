package upload

func Upload(keyword, name, title string, keywords, names, titles []string) ([]string, []string, []string) {
	keywords = append(keywords, keyword)
	names = append(names, name)
	titles = append(titles, title)
	return keywords, names, titles
}
