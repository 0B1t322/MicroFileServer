package encode_test

import (
	"fmt"
	"testing"

	"github.com/MicroFileServer/pkg/urlvalue/encode"
)

func TestFunc_Unmarshall(t *testing.T) {
	type B struct {
		Number int `query:"number,int"`
	}
	type A struct {
		Start 		int  	`query:"start,int" json:",inline"`
		Count 		int    	`query:"count,int"`
		Name  		string 	`query:"name,string"`
		False	  	bool	`query:"False,bool"`
		Empty		bool	`query:"empty,bool"`
		True		bool	`query:"True,bool"`
		B
	}

	a := &A{}

	if err := encode.UrlQueryUnmarshall(
		a,
		map[string][]string{
			"start":  {"10"},
			"count":  {"12"},
			"name":   {"orbit"},
			"number": {"27"},
			"False":	{"false"},
			"True":	{"true"},
		},
	); err != nil {
		t.Log(err)
	}

	fmt.Printf("%+v", a)
}
