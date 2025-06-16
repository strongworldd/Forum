package models

import (
	"sync"
)

type Vote struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Type     int    `json:"type"` // 1 pour upvote, -1 pour downvote
}

var (
	votes     = make(map[string]map[string]int) // map[username]map[title]voteType
	votesLock sync.RWMutex
)

func AddVote(username, title string, voteType int) error {
	votesLock.Lock()
	defer votesLock.Unlock()

	if _, exists := votes[username]; !exists {
		votes[username] = make(map[string]int)
	}

	// Vérifier si l'utilisateur a déjà voté
	if currentVote, exists := votes[username][title]; exists {
		// Si le vote est le même, on le retire
		if currentVote == voteType {
			delete(votes[username], title)
			return nil
		}
		// Si le vote est différent, on le change
		votes[username][title] = voteType
		return nil
	}

	// Nouveau vote
	votes[username][title] = voteType
	return nil
}

func GetUserVotes(username string) map[string]int {
	votesLock.RLock()
	defer votesLock.RUnlock()

	if userVotes, exists := votes[username]; exists {
		return userVotes
	}
	return make(map[string]int)
}

func GetVoteCount(title string) int {
	votesLock.RLock()
	defer votesLock.RUnlock()

	count := 0
	for _, userVotes := range votes {
		if voteType, exists := userVotes[title]; exists {
			count += voteType
		}
	}
	return count
}
