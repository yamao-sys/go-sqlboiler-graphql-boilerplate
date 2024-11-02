package factories

import (
	models "app/models/generated"
	"log"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"golang.org/x/crypto/bcrypt"
)

var UserFactory = factory.NewFactory(
	&models.User{
		Name:  randomdata.FullName(randomdata.RandomGender),
		Email: randomdata.Email(),
	},
).Attr("Password", func(args factory.Args) (interface{}, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to generate hash %v", err)
	}
	return string(hash), nil
})
