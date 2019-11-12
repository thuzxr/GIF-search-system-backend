package upload

func Upload(users, names, titles, infos, keywords []string, user, name, title, info, keyword string) ([]string, []string, []string, []string, []string) {
	keywords = append(keywords, keyword)
	names = append(names, name)
	titles = append(titles, title)
	users = append(users, user)
	infos = append(infos, info)
	return users, names, titles, infos, keywords
}
