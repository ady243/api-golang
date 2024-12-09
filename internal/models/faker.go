package models

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func randomStringElement(elements []string) string {
	return elements[rand.Intn(len(elements))]
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randomString(length int) string {
	const charset = "abcdefghijklmnop" + "qrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateFakeUser() Users {
	now := time.Now()
	birthDate := time.Date(rand.Intn(30)+1970, time.Month(rand.Intn(12)+1), rand.Intn(28)+1, 0, 0, 0, 0, time.UTC)

	return Users{
		ID:            ulid.Make().String(),
		Username:      "user" + randomString(5),
		Email:         randomString(5) + "@temUp.com",
		PasswordHash:  randomString(12),
		UpdatedAt:     now,
		DeletedAt:     nil,
		BirthDate:     &birthDate,
		ProfilePhoto:  "https://example.com/photo" + ulid.Make().String(),
		FavoriteSport: "sport" + ulid.Make().String(),
		SkillLevel:    "level" + ulid.Make().String(),
		Bio:           "This is a bio for user " + ulid.Make().String(),
		Pac:           randomInt(1, 100),
		Sho:           randomInt(1, 100),
		Pas:           randomInt(1, 100),
		Dri:           randomInt(1, 100),
		Def:           randomInt(1, 100),
		Phy:           randomInt(1, 100),
		MatchesPlayed: randomInt(1, 100),
		MatchesWon:    randomInt(1, 100),
		GoalsScored:   randomInt(1, 100),
		BehaviorScore: randomInt(1, 100),
		RefreshToken:  ulid.Make().String(),
	}
}

func GenerateFakeUsers(count int) []Users {
	if count >= 10 {
		count = 10
	}
	users := make([]Users, count)
	i := 0
	for i < count {
		users[i] = GenerateFakeUser()
		i += 1
	}
	return users
}
