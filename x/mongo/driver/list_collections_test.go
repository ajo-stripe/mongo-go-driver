// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"testing"
)

func addRegexFitler(filter bsonx.Doc) bsonx.Doc {
	regexFilter := bsonx.Doc{
		{"name", bsonx.Regex("^[^$]*$", "")},
	}
	if filter == nil {
		return regexFilter
	}

	arr := bsonx.Arr{bsonx.Document(regexFilter), bsonx.Document(filter)}
	return bsonx.Doc{
		{"$and", bsonx.Array(arr)},
	}
}

func TestListCollections(t *testing.T) {
	dbName := "db"
	noNameFilter := bsonx.Doc{
		{"foo", bsonx.String("bar")},
	}
	nonStringFilter := bsonx.Doc{
		{"name", bsonx.Int32(1)},
	}
	nameFilter := bsonx.Doc{
		{"name", bsonx.String("coll")},
	}
	modifiedFilter := bsonx.Doc{
		{"name", bsonx.String(dbName + ".coll")},
	}

	t.Run("TestTransformFilter", func(t *testing.T) {
		testCases := []struct {
			name           string
			filter         bsonx.Doc
			expectedFilter bsonx.Doc
			err            error
		}{
			{"TestNilFilter", nil, addRegexFitler(nil), nil},
			{"TestNoName", noNameFilter, addRegexFitler(noNameFilter), nil},
			{"TestNonStringName", nonStringFilter, nil, ErrFilterType},
			{"TestName", nameFilter, addRegexFitler(modifiedFilter), nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				newFilter, err := transformFilter(tc.filter, dbName)
				require.Equal(t, tc.err, err)
				require.Equal(t, tc.expectedFilter, newFilter)
			})
		}
	})
}