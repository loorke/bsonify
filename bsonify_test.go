package bsonify

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestToBSON(t *testing.T) {
	{
		opts := a{Ab: true}
		d := SetUpdateD(opts)
		for i := range d {
			d[i].Key += "OLOLo"
		}

		spew.Dump(d)
	}

	var (
		n    int = 1.0
		opts     = a{Ab: true, Ae: c{Ca: &n}}
	)

	require.Equal(t,
		bson.D{
			{"custom_name", true},
			{"Ac", ""},
			{"ae", bson.D{
				{"ca", 1},
				{"cb", (*uint)(nil)},
			}},
		}, Dump(opts))

	require.Equal(t,
		bson.D{
			{"custom_name", true},
			{"Ac", ""},
			{"ae.ca", 1},
			{"ae.cb", (*uint)(nil)},
		}, SetUpdateD(opts))

	require.Equal(t,
		bson.M{
			"custom_name": true,
			"Ac":          "",
			"ae.ca":       1,
			"ae.cb":       (*uint)(nil),
		}, SetUpdateM(opts))

}

type a struct {
	Aa string `bson:"aa,omitempty"`
	Ab bool   `bson:"custom_name,omitempty"`
	Ac string

	Ad b `bson:"ad,omitempty"`
	Ae c `bson:"ae"`
}

type b struct {
}

type c struct {
	Ca *int  `bson:"ca,omitempty"`
	Cb *uint `bson:"cb"`
}
