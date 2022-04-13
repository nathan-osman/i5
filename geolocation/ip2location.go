package geolocation

import (
	"github.com/ip2location/ip2location-go"
)

const ip2locationType = "ip2location"

type ip2locationProvider struct {
	db *ip2location.DB
}

func newIp2locationProvider(path string) (geolocationProvider, error) {
	d, err := ip2location.OpenDB(path)
	if err != nil {
		return nil, err
	}
	return &ip2locationProvider{
		db: d,
	}, nil
}

func (i *ip2locationProvider) Country(ip string) string {
	r, err := i.db.Get_country_short(ip)
	if err != nil {
		return ""
	}
	return r.Country_short
}

func (i *ip2locationProvider) Close() {
	i.db.Close()
}
