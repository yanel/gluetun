package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func IpvanishCountryChoices(servers []models.IpvanishServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func IpvanishCityChoices(servers []models.IpvanishServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func IpvanishHostnameChoices(servers []models.IpvanishServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
