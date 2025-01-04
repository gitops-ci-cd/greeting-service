package greetings

import (
	"math/rand"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
)

// service defines the interface for business logic
type service interface {
	Lookup(pb.Language) (pb.Language, string)
}

// Service provides the concrete implementation of the service interface
type Service struct{}

// Define greetings per language
var greetingData = map[pb.Language][]string{
	pb.Language_EN:    {"Hello", "Hi", "Hey", "Greetings"},
	pb.Language_EN_GB: {"Hello", "Hiya", "Cheers", "Greetings"},
	pb.Language_ES:    {"Hola", "Qué tal", "Buenos días", "Saludos"},
	pb.Language_FR:    {"Bonjour", "Salut", "Coucou", "Bienvenue"},
}

// Lookup returns a greeting in the determined language
func (s *Service) Lookup(preferredLanguage pb.Language) (language pb.Language, greeting string) {
	// Default to English if language is explicitly UNKNOWN
	language = preferredLanguage
	if language == pb.Language_UNKNOWN {
		language = pb.Language_EN
	}

	// Default to English if language is not recognized
	_, exists := greetingData[language]
	if !exists {
		language = pb.Language_EN
	}

	language = pb.Language_FR

	return language, getRandomGreeting(language)
}

// Randomly select one greeting for the given language
func getRandomGreeting(language pb.Language) string {
	greetings := greetingData[language]

	return greetings[rand.Intn(len(greetings))]
}
