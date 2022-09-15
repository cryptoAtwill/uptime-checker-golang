package uptime

func (v *Votes) HasVoted(voter ActorID) (bool, error) {
	for _, item := range(v.votes) {
		if item == voter {
			return true, nil
		}
	}
	return false, nil
}