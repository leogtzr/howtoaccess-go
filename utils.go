package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

// String ...
func (a Access) String() string {
	if len(a.Notes) > 0 {
		return fmt.Sprintf("%d: %s@%s from %s (%s)", a.ID, a.ServerDestination, a.UserDestination, a.From, a.Notes)
	}
	return fmt.Sprintf("%d: %s@%s from %s", a.ID, a.ServerDestination, a.UserDestination, a.From)
}

func extractAccessesFromFile(file io.Reader) ([]Access, error) {
	records := make([]Access, 0)

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return []Access{}, err
	}

	for i, record := range lines {
		if i == 0 {
			continue
		}
		if len(record) != requiredNumberOfFields {
			return []Access{}, fmt.Errorf("wrong number of fields in line: %d", i)
		}
		records = append(records, Access{
			ID:                i,
			ServerDestination: record[0],
			UserDestination:   record[1],
			From:              record[2],
			Notes:             record[3],
		})
	}

	return records, nil
}

// Equal ...
func Equal(a, b []Access) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func getNextIndex(accesses *[]Access) int {
	if len(*accesses) == 0 {
		return -1
	}
	lastAccessElement := (*accesses)[len(*accesses)-1]
	return lastAccessElement.ID + 1
}

func removeElementByID(id int, accesses *[]Access) {
	id = getIndexByID(id, accesses)
	if id == -1 {
		return
	}
	*accesses = append((*accesses)[:id], (*accesses)[id+1:]...)
}
