package worker

var subscribers = make(map[string]bool)

func Register(url string) {
	subscribers[url] = true
}

func Unregister(url string) {
	delete(subscribers, url)
}
