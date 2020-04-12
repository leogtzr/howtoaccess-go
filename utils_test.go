package main

import (
	"io"
	"strings"
	"testing"
)

func TestAccess_String(t *testing.T) {
	type test struct {
		acc  Access
		want string
	}

	tests := []test{
		{
			acc: Access{
				ID:                1,
				ServerDestination: "alla",
				UserDestination:   "el",
				From:              "aca",
				Notes:             "jijij",
			},
			want: "1: alla@el from aca (jijij)",
		},

		{
			acc: Access{
				ID:                3,
				ServerDestination: "alla",
				UserDestination:   "el",
				From:              "aca",
				Notes:             "",
			},
			want: "3: alla@el from aca",
		},
	}

	for _, tt := range tests {
		if got := tt.acc.String(); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_extractAccessesFromFile(t *testing.T) {
	type test struct {
		file       io.Reader
		records    []Access
		shouldFail bool
	}

	tests := []test{
		{
			file: strings.NewReader(`Destination (Server Name),User (Destination),Access from,Notes
pr-galaxie-xl25,up000176,corolla01,
pr-galaxie-xl25,up000176,hola123,Password`),
			records: []Access{
				Access{
					ID:                1,
					ServerDestination: "pr-galaxie-xl25",
					UserDestination:   "up000176",
					From:              "corolla01",
					Notes:             "",
				},
				Access{
					ID:                2,
					ServerDestination: "pr-galaxie-xl25",
					UserDestination:   "up000176",
					From:              "hola123",
					Notes:             "Password",
				},
			},
			shouldFail: false,
		},
		{
			file: strings.NewReader(`Destination (Server Name),User (Destination),Access from,Notes
pr-galaxie-xl25,up000176,corolla01,
pr-galaxie-xl25,up000176,hola123,Password,Another,Hmmmm`),
			records:    []Access{},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		got, err := extractAccessesFromFile(tt.file)
		if hasError := err != nil; hasError != tt.shouldFail {
			t.Errorf("It should have failed to parse, flag is = [%t]", tt.shouldFail)
		}
		if !Equal(got, tt.records) {
			t.Errorf("got=[%s], want=[%s]", got, tt.records)
		}
	}
}

func TestEqual(t *testing.T) {
	type test struct {
		a      []Access
		b      []Access
		result bool
	}

	tests := []test{
		{
			a: []Access{
				Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
				Access{ID: 2, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
			},
			b: []Access{
				Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
				Access{ID: 2, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
			},
			result: true,
		},
		{
			a: []Access{
				Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
				Access{ID: 2, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
			},
			b: []Access{
				Access{ID: 3, ServerDestination: "s2", UserDestination: "u1", From: "f1", Notes: "n1"},
			},
			result: false,
		},
		{
			a: []Access{
				Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
				Access{ID: 2, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
			},
			b: []Access{
				Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
				Access{ID: 3, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
			},
			result: false,
		},
	}

	for _, tt := range tests {
		if got := Equal(tt.a, tt.b); got != tt.result {
			t.Errorf("[%s] and [%s] should be equal", tt.a, tt.b)
		}
	}
}

func Test_getNextIndex(t *testing.T) {
	type test struct {
		accesses []Access
		want     int
	}

	tests := []test{
		{
			accesses: []Access{
				Access{ID: 1},
				Access{ID: 2},
			},
			want: 3,
		},
		{
			accesses: []Access{},
			want:     -1,
		},
	}

	for _, tt := range tests {
		if got := getNextIndex(&tt.accesses); got != tt.want {
			t.Errorf("got=[%d], want=[%d]", got, tt.want)
		}
	}
}

func Test_removeElementByID(t *testing.T) {

	type test struct {
		accesses   []Access
		idToRemove int
		want       []Access
	}

	tests := []test{
		{
			accesses: []Access{
				Access{ID: 3},
				Access{ID: 5},
				Access{ID: 6},
			},
			idToRemove: 5,
			want: []Access{
				Access{ID: 3},
				Access{ID: 6},
			},
		},

		{
			accesses: []Access{
				Access{ID: 3},
				Access{ID: 5},
				Access{ID: 6},
			},
			idToRemove: 56,
			want: []Access{
				Access{ID: 3},
				Access{ID: 5},
				Access{ID: 6},
			},
		},
	}

	for _, tt := range tests {
		removeElementByID(tt.idToRemove, &tt.accesses)
		if !Equal(tt.accesses, tt.want) {
			t.Errorf("got=[%s], want=[%s]", tt.accesses, tt.want)
		}
	}

}

func TestAccess_ToCSV(t *testing.T) {
	type test struct {
		access Access
		want   string
	}

	tests := []test{
		test{
			access: Access{ID: 1, ServerDestination: "a", UserDestination: "b", From: "c", Notes: "d"},
			want:   "a,b,c,d",
		},
		test{
			access: Access{ID: 1, ServerDestination: "a", UserDestination: "b", From: "c"},
			want:   "a,b,c,",
		},
	}

	for _, tt := range tests {
		if got := tt.access.ToCSV(); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_searchByID(t *testing.T) {

	accesses := []Access{
		Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
		Access{ID: 2, ServerDestination: "s2", UserDestination: "u2", From: "f2", Notes: "n2"},
		Access{ID: 3, ServerDestination: "s3", UserDestination: "u3", From: "f3", Notes: "n3"},
	}

	type test struct {
		id       int
		accesses []Access
		wasFound bool
		want     Access
	}

	tests := []test{
		test{
			id:       1,
			accesses: accesses,
			want:     Access{ID: 1, ServerDestination: "s1", UserDestination: "u1", From: "f1", Notes: "n1"},
			wasFound: true,
		},

		test{
			id:       5,
			accesses: accesses,
			want:     Access{},
			wasFound: false,
		},
	}

	for _, tt := range tests {
		if got, found := searchByID(tt.id, &tt.accesses); found != tt.wasFound || got != tt.want {
			t.Errorf("searchByID, expecting [%t], got [%t]", tt.wasFound, found)
		}
	}
}
