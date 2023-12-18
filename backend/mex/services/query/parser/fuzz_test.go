package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomQuery() string {
	// We make MEx special characters more likely to create more challenging strings
	probOfSpace := 0.1
	probOfMexChar := 0.2
	maxLength := 25
	spaceSymbol := []rune(` `)
	mexSymbols := []rune(`()+|-*\"`)
	normalChars := []rune(`abcdefghijklmnopquertuxyzäößABCDEFGHIJKLMNOPQUERTUXYZÄÖ0123456789/.:;?'&$@~!'`)

	strLen := randGen.Intn(maxLength) + 1
	str := strings.Builder{}
	var newSym rune
	for i := 0; i < strLen; i++ {
		r := randGen.Float64()
		switch {
		case r < probOfSpace:
			newSym = spaceSymbol[0]
		case r < probOfSpace+probOfMexChar:
			newSym = mexSymbols[randGen.Intn(len(mexSymbols))]
		default:
			newSym = normalChars[randGen.Intn(len(normalChars))]
		}
		str.WriteRune(newSym)
	}
	res := str.String()
	return res
}

/*
These are fuzz tests for the query parser - they generate random input and test that it does not provoke errors.
*/
func Test_ParseSearchQuery_Fuzzing(t *testing.T) {
	numTest := 1000
	for i := 0; i < numTest; i++ {
		matchingOpConfig := MatchingOpsConfig{
			Term: []MatchingFieldConfig{
				{
					FieldName:       "A",
					BoostFactor:     "",
					MaxEditDistance: DefaultMaxTermEditDistance,
				},
				{
					FieldName:       "B",
					BoostFactor:     "1",
					MaxEditDistance: DefaultMaxTermEditDistance,
				},
			},
			Phrase: []MatchingFieldConfig{
				{
					FieldName:       "C",
					BoostFactor:     "",
					MaxEditDistance: 0,
				},
				{
					FieldName:       "D",
					BoostFactor:     "2",
					MaxEditDistance: 0,
				},
			},
		}
		listener, _ := NewListener(matchingOpConfig)
		parser := NewMexParser(listener)
		query := getRandomQuery()
		t.Run(fmt.Sprintf("Fuzz test: query + '%s'", query), func(t *testing.T) {
			_, err := parser.ConvertToSolrQuery(query)
			if err != nil {
				t.Errorf("query error for query '%s': %s", query, strings.Join(err.GetMessages(), ", "))
				return
			}
		})
	}
}
