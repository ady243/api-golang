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
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
		ID:            ulid.Make(),
		Username:      "user" + ulid.Make().String(),
		Email:         "user" + ulid.Make().String() + "@example.com",
		PasswordHash:  randomString(12),
		UpdatedAt:     now,
		DeletedAt:     nil,
		BirthDate:     &birthDate,
		Role:          Role(randomStringElement([]string{"player", "referee", "administrator"})),
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
	users := make([]Users, count)
	for i := 0; i < count; i++ {
		users[i] = GenerateFakeUser()
	}
	return users
}
