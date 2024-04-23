package packages

func FanIn(participants []string, ch chan<- string) {
	for _, participant := range participants {
		ch <- participant
	}
}
