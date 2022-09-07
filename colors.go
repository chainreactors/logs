package logs

func Red(s string) string {
	return "\033[1;31m" + s + "\033[0m"
}

func Green(s string) string {
	return "\033[1;32m" + s + "\033[0m"
}

func Yellow(s string) string {
	return "\033[4;33m" + s + "\033[0m"
}

func Blue(s string) string {
	return "\033[1;34m" + s + "\033[0m"
}

func Purple(s string) string {
	return "\033[1;35m" + s + "\033[0m"
}

func Cyan(s string) string {
	return "\033[1;36m" + s + "\033[0m"
}

func White(s string) string {
	return "\033[1;37m" + s + "\033[0m"
}
