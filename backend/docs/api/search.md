# The MEx Search API

This document describes how to use the search endpoint (`POST v0/query/search`) to search the MEx data.
This endpoint uses the underlying [Apache Solr](https://solr.apache.org/) search engine to perform the search.
However, the API is designed so that the clients do not need to know anything about Solr, and that the API can remain  unchanged even if the underlying search index is changed.

The search format is intimately connected with the structure of the search index, as defined by the MEx configuration of supported fields, and of the groups of fields to search together (_search foci_) or to use for faceting/ordering (_ordinal/ hierarchy axes_).
See the detailed [documentation of the search index configuration](../../../docs/metadata_config.md) for details.

## Request format

A search is triggered by a POST request to the `/search` endpoint.
The query is specified entirely in the body of the request body.
An example of q query illustrating many (but not all) available options is the following:

```json
{
  "query": "disease | influence",
  "axisConstraints": [
    {
      "type": "exact",
      "axis": "author",
      "values": ["John Doe"]
    }
  ],
  "limit": 5,
  "offset": 10,
  "fields": ["author", "title"],
  "sorting": {
    "axis": "author",
    "order": "asc"
  },
  "searchFocus": "title",
  "facets": [
    {
      "type": "exact",
      "axis": "author",
      "limit": 3,
      "offset": 0
    }
  ],
  "highlightFields": ["title"],
  "maxEditDistance": 1
}
```

The search query includes search criteria and a description of the information to be returned - the meaning of the individual properties are explained below.
The above query would find items which either (or both) of the terms "disease" and "influence" is found in a title-relevant fields, restricted the search items for which John Doe is named as author.
For the matching items, the fields "author" and "title" will be returned (in addition to the item ID).
A total of at most 5 matches, starting with match index 10 (zero-based), will be returned (the 10 first matches are skipped).
The results will be sorted in ascending order by author.
In addition, facet buckets of items counts based on author will be returned (up to a max of 3), along with highlight-snippets for matches in the title field.

### User-entered search query: `query`

In a search UI, this is what the user enters into the search box.
It should be sent as entered by the users (apart from optional trivial changes like trimming empty spaces from start and end).
An empty search string (or one consisting only of spaces) or the query `*` return all searchable items (or, to be exact: all items with entries in at least one of the fields being searched).

The MEx search query uses a MEx-specific query language, supporting phrase search, wildcard search, grouping and basic Boolean operators (AND, OR, and NOT).
The two simplest features of this language is wildcard search using `*` (e.g. `super*`) and search for phrases enclosed by double quotes (e.g. `"back pain"`).
See the [MEx query language document](../../../docs/query_language.md) for details about the language.
The following characters are assigned special meaning in the MEx  query language:

```
+ | - " * ) (
```

The MEx query engine will attempt to parse every search query and use the parsed result to build the query for the underlying search engine (currently Solr).
However, in line with the design goal of making MEx useful to non-expert users, it is not necessary to know the MEx query language, e.g. one can simply search for a set of separated words (`back pain chronic`) to find items containing all these terms.
Users can only do structured search as far as allowed by the MEx language.
There is no "raw mode" allowing them to formulate queries directly in the language used by the underlying search engine (currently Solr), in line with the principle of decoupling the search API from the underlying technology stack.
This separation is kept for both usability , engineering, and security reasons:

- it sets a clear scope for allowed queries , helping to ensure that all features are working and documented
- it ensures that the underlying search engine can be replaced without any working queries breaking
- it prevents injection attacks that aim to trigger problematic/slow searches in the underlying search engine

The MEx query engine will sanitize input to ensure that characters with special meaning to the underlying search engine are escaped or removed before being passed on.
The sanitation will aim to err on the side of caution, e.g. cases of multiple escape characters may be simplified to ensure that no unwanted search is triggered.

#### Handling of invalid queries

Rather than enforcing strict syntax checks based on the MEx language, the search engine is tolerant of errors.
For instance, special MEx symbols are simply ignored if they do not form part of valid expression, as may happen if the user makes a typo (e.g. accidentally adds a plus at the end of term, as in `hand+`).
However, if the user types in a valid MEx expression - even if by accident - it will be interpreted as such and the corresponding search executed.
This should only very rarely lead to confusion but may cause unexpected behavior on some rare occasions.
If, for instance, a user enters a query with a chemical formula for an ion like `Ca(2+) transport`, the query engine would interpret the brackets and `+` as search operators.
Note that in such cases, placing a phrase in double quotes (`"Ca(2+)"`) will ensure that special symbols are not interpreted as search operators.
Conversely, the tolerant handling of "dangling symbols" implies that advanced users intending to do a structured search will not be explicitly told if their query did not fit the MEx syntax.
However, the return diagnostic information (see below) allows clients to know how a given query was interpreted by the MEx query engine and whether characters were ignored.
The diagnostic information also indicates whether the query parsing failed and, in that case, includes parsing error messages.
Regardless of whether the parsing succeeded or not, a string indicating the detected logical structure of the query submitted to the search engine will also be returned.
See the section on diagnostics below for details.

### Result ordering: `sorting`

By default, the returned results are sorted by the match ranking (in descending order, i.e. best match first), using the match-score of the underlying search engine.
To sort the results by other values, one or more so-called _ordinal axis_ need to be defined.
An ordinal axis is basically a named set of fields that are used for sorting and faceting (see the [search index configuration documentation](../../../docs/metadata_config.md) for details).
In the simplest and most common case, an ordinal axis consists of just one field; sorting by this axis then simply correspond to sorting by this field.
Suppose there is an ordinal axis called "author" to sort by, well, author.
You can then sort in ascending order by the content of the author field as follows:

```json
"sorting": {
  "axis": "author",
  "order": "asc"
}
```

The entry in `axis` is required and must be the name of ordinal axis.
`order` must be "asc" (ascending) or "desc" (descending), respectively - anything else causes as error.
If only the axis name is given, the results are returned in ascending order.
Typically, ordinal axis will be defined so that, for any given item (1) exactly one of the included fields (if there is more than one) has an entry, and (2) the field has a single entry.
This ensures that each item is sorted based on a single value, ensuring an unambiguous sorting.
However, it is allowed to define sort axis that picks out multiple values for single item to use for sorting (e.g. axes based on multi-valued fields or on several fields that are not mutually exclusive)
In such cases, the minimum value will be used for ascending ordering and the maximum value for descending ordering.
For instance, if the ordinal axis `keyword` is based on the multivalued field of the same name and a given item contains the values `[age, trend, nose]` in this field, the item will occur in the position appropriate for `age` in an ascending sort and in the position appropriate the word `trend` in a descending sort.

For text fields, the search order is the alphabetic order.
The specific order of characters differs between languages.
For instance, in German common accented characters are sorted with the base character (e.g. "ü" with "u") whereas Danish treats the character "å" as unrelated to "a".
Likewise, while both and Danish and Swedish use the character "å", they put it in different places in the alphabetic ordering.
MEx does a "best effort" sort that gives the correct ordering for the currently supported languages (English and German) and should perform well for other Roman or Germanic langauges.
However, words containing unknown characters or non-word symbols (e.g. brackets) may not occur in the place the user expects.

### Set of fields to be searched (search focus): `searchFocus`

The search is always carried out on a specific set of fields.
Such a set of fields to be searched are called a _search focus_ (see the [search index configuration documentation](../../../docs/metadata_config.md) for details).
To set the search focus, add the name of the focus in the `searchFocus` property.
Every MEx system must always have a search focus called "default" configured.
This search focus is used if no other focus is specified explicitly in the `searchFocus` property.
Hence, it will typically contain a broad spectrum of fields that are assumed to be relevant to users.
However, it will usually not include _all_ fields since some fields may be useful for display but not for search.
For instance, items could have a field that stores the full text of the license terms that apply to the underlying data.
This may be useful to show users, but is unlikely to be helpful to search since they are standardized texts and large groups of items may have identical licenses.
Likewise, there may be fields that are only relevant for search in a very specific context that is better covered by a targeted search focus (e.g. for names of contact persons).

One can configure further search foci to target the search on sets of fields relevant to a specific aspect of the data.
For instance, one might want a search focus for searching only in titles or abbreviated titles of the metadata entities.

### Fields to be returned: `fields`

This is an array with fields to be returned.
Each field is referenced using the name it has in the configuration.
Since the field `id` is needed to structure the response (and connect results to data in DB), the field `id` is automatically added to the list of fields if not already included.
Hence, if the `fields` property is present but contains an empty array, only the `id`-values are returned.
There is currently no way to request that all fields be returned.

### Pagination: `limit` and `offset`

- `limit` is the number of items to be returned. It defaults to 0. The limit is capped at 1000 - larger values will be silently reset to this value.
- `offset` is the index of the first item to be returned (used for paging) - it defaults to 0.

### Constraints on axes: `axisConstraints`

Axis constraints fix the allowed values along a specific ordinal axis to either one or more exact values, or a range of value.
They are typically imposed when the user chooses to constrain to one or more buckets or ranges given by a facet (see below).
There are currently two types of constraints, exact and string/datetime ranges - the type has to be indicated in the `type` property.

#### Constraining to exact values

These are constraints that require the value in a given field to be exactly one of a set of allowed value.
This is used when the possible field values are discrete (e.g. words) and can be used e.g. to constrain the results to those in one or more selected buckets from an exact facet.

The constraint object looks as follows (a further property can be added for hierarchy axes, see below).

```json
"axisConstraints": [
  {
    "type": "exact",
    "axis": "author",
    "values": [
      "Does",
      "Jones"
    ],
    "combineOperator": "or"
  }
]
```

Here

- `axis` gives the ordinal axis to which the constraint is to be applied (see documentation of the `sorting` parameter for documentation of what an ordinal axis is)
- `values` is an array of values to constrain to
- `combineOperator` is the name of the binary boolean operator that should be used for combining the constraints for
  individual values (optional - defaults to `or`)

Currently, the only supported combination operators are `or` and `and`.
If set to `and`, all the entries in `values` have to occur in an item simultaneously to get a hit.
If set to `or`, only one of the entries has to be present.

#### Constraining on a hierarchy axis

For hierarchy axes, setting an axis constraint to a given value (code) will restrict the returned results to those items that have this value _or a child value_ (in the code hierarchy) in the relevant fields.
That is, if we imagine the hierarchy as a tree, restricting to the code for a specific node X will pick out all items carrying the code for X or a code for any of the nodes in the subtree below X.
Suppose, for instance, that projects are assigned to an organizational unit that is part of a hierarchy of units and that we have a corresponding hierarchy field configured.
We have configured a hierarchical axis called `orgUnit` with just this field.
In the organization hierarchy, the health research unit (code `HR`) has the sub-units epidemiology (code `EP`) and clinical research (code `CR`).
Consider now the following constraint:

```json
"axisConstraints": [
  {
    "type": "exact",
    "axis": "orgUnit",
    "values": [ "HR" ],
    "combineOperator": "or"
  }
]
```

This would restruct to all projects assigned to either health research directly OR epidemiology OR clinical research (or multiple of them).

Selecting a subtree (node plus all child nodes) is usually the desired behavior when dealing with hierarchies.
E.g. if we search for specific disease, we are usually also interested in the more specific forms of that diseases (child nodes of the node for the main disease).
However, there may be cases where sub a subtree selection is not desired.
For instance, the health research unit may have its own office, employees etc. distinct from any of its subunits.
In that case, a project could be managed directly by the health research unit and not any of its subunits.
To pick out such cases by search, we want to pick out only projects that are assigned _directly_ to the health research unit.
This is done using the `singleNodeValues` property: it works just like `values` property, but does not include child nodes (this property is ignored for non-hierarchy axes).
Hence, to pick out projects managed directly by the health research unit and not a child unit, we would use the following constraint:

```json
"axisConstraints": [
  {
    "type": "exact",
    "axis": "orgUnit",
    "singleNodeValues": [ "HR" ],
    "combineOperator": "or"
  }
]
```

As with the normal `values` property, multiple values can be passed to `singleNodeValues`, and like `value` entries, they are combined as specified by `combineOperator`.
Hierarchical and single node values can also be combined in a single constraint clause, in which case the resulting constraints are also combined as specified by `combineOperator`.

#### Constraining to string or datetime ranges

These constrain forces a string or datetime field value to fall in a particular interval.
It is used e.g. when the user selects a range of bins returned by a "yearRange" facet.
The format is

```json
"axisConstraints": [
  {
    "type": "stringRange",
    "axis": "createdAt",
    "stringRanges": [
      {
        "min": "2014-03-16T17:33:18Z",
        "max": "2018-07-22T06:12:26Z"
      },
      {
        "min": "2021-03-16T17:33:18Z"
      }
    ],
    "combineOperator": "or"
  }
]
```

The `field` and `combineOperator` works as described for the exact constraints e.g. separate range constraints are combined with the specified operator (defaults to `or`).
The array property `stringRanges` describe ranges on strings and datetimes (given as strings).
Each range object can have a `min` and `max` string, at least one of which must be given.
If only a `min` value is given, the range has no upper limit (includes all values above `min`), if only a `max` value is given, it has no lower limit (includes all item below `max`).
The intervals are inclusive at both the upper and lower boundaries, i.e. if we have

```json
"axisConstraints": [
  {
    "type": "stringRange",
    "axis": "journalName",
    "stringRanges": [
      {
        "min": "Nature",
        "max": "Science"
      }
    ]
  }
]
```

then items with `journalName` axis values _equal to_ "Nature" or "Science" will satisfy the constraint and will be returned.

_NOTE: The ordering used when applying range constraints to string fields is currently **case-sensitive**.
Specifically, it places all uppercase characters before all lowercase character.
Hence, for instance, the filter above would filter out any item with `journalName` set to e.g. "nature" or "research today" (because these start with a lowercase and hence both come after "Science")._

#### Combinations of constraints on different fields

If there are constraints on multiple fields, the constraint on each field must be satisfied by an item for it to match (i.e. constraints on different fields are combined with a logical AND).
Hence, if we set the following field constraints,

```json
"axisConstraints": [
  {
    "type": "stringRange",
    "axis": "journalName",
    "stringRanges": [
      {
        "min": "Nature",
        "max": "Science"
      }
    ]
  },
  {
    "type": "exact",
    "axis": "author",
    "values": [
      "Doe",
      "Jones"
    ]
  }
]
```

we get all items with `journalName` in the (inclusive) range from "Nature" to "Science" for which the author list also includes either "Doe" or "Jones" (or both).

### Faceting: `facets`

This specifies an array of objects representing facets to be returned in addition to the matches, i.e. match counts (or other statistics) for different categories.
There are three kinds of facets currently available in MEx:

1. exact (categorical) facets
2. year-range facets,
3. statistical string facets

The type of facet is specified by the `type` field in the facet object.
Depending on the type, different options can be specified - options not relevant for the facet type will be ignored.

#### Exact (categorical) facets

Exact facets divides the matches into categories (buckets) depending on the value in an ordinal or hierarchy axis (cf. the discussing of sorting).
Faceting parameters does not affect the returned list of matches.
The format is this:

```json
"facets": [
  {
    "type": "exact",
    "axis": "author",
    "limit": 5,
    "offset": 10
  }
]
```

`axis` indicates on which ordinal axis to facet by (see documentation of the `sorting` parameter to see what an ordinal axis is) - this property is required
`limit` and `offset` are optional paging parameters, setting the number of buckets to return (ranked by count) and the index of the first bucket, respectively.
_NOTE: `limit` defaults to 0, so it _must_ be set if returned buckets are required_.
It is capped at 1000 - larger limits will be silently reset to this value.

Note that MEx uses what is known as _multi-select faceting_.
That is, if you have an axis constraint on a field X and also facet on X, then the axis constraint on X will be ignored for the purpose of this specific facet, but will still be applied everywhere else (e.g. when finding matching items or when faceting other fields).
This ensures that the returned facets always contain all options available given the query and other constraints, not just the values or ranges to which the field was restricted.
For a discussion of multi-select faceting, see [this blog post discussing it in the context of the Solr search engine](https://yonik.com/multi-select-faceting/).

The facet section of the response will look like the following.

```json
"facets": [
   {
    "type": "exact",
    "axis": "author",
    "bucketNo": 16,
    "buckets": [
      {
        "value": "Carol Greene",
        "count": 1
      },
      {
        "value": "Carrie McCracken",
        "count": 1
      },
      {
        "value": "Cynthia Bearer",
        "count": 1
      }
    ]
  }
]
```

#### Year-range facets

These are facets that sort matches into bins spanning exactly one calendar year (01 Jan - 31 Dec), based on the value in a specified datetime/timestamp field.
The MEx search engine will automatically set the range so that it exactly covers all matches returned by the query.
Since the range is detected by an extra query run in the background, using year-range facets will make the overall search slightly slower, though the difference should usually be negligible.
Note that for date fields covering either a very small time-range (one or a few years) or a very long time-spans (say, a century), binning by year can lead to facets with either very few or very many bins.

The format for a year range facet is this:

```json
"facets": [
  {
    "type": "yearRange",
    "axis": "createdAtAxis"
  }
]
```

Like the exact facets, we use multi-select faceting and hence ignore any axis constraints on the field being faceted (see above).

The returned data will have the same format as for exact facets, except that "bucketNo" will be missing or default to 0.
The `value` property gives the starting date of the bin (start of year).
The `type` property will be set to "yearRange".

### Statistical facet for string/date fields

These are facets that do not return bins (buckets), but rather single values computed based on the matching items for the query.
Currently, this is only implemented for string and timestamp fields.
The format is this.

```json
"facets": [
  {
    "type": "stringStat",
    "axis": "publicationDate",
    "statName": "min_publicationDate",
    "statOp": "min"
  }
]
```

Here, `statName` is a client-chosen name that is used to identify the result in the returned data.
It is required and must be unique for each facet in a single request.
`statOp` gives the operator that should be applied to the chosen field to compute the result.
Currently, only two are available.

- `min`: minimum value across items (earliest datetime, first string when alphabetically sorted)
- `max`: maximum value across items (latest datetime, last string when alphabetically sorted)

_NOTE: Currently, the `min` and `max` operators are **case sensitive**, placing all uppercase characters before all lowercase character.
Hence, given the set of terms ["Zone", "art"], `min` would return "Zone" and `max` "art" since uppercase "Z" is sorted before lowercase "a"._

Unlike the exact and string-range facets, statistical facet do not use multi-select faceting (cf. description of how this works for exact facets).
That is, if you have a statistical facet on the field X, axis constraints on X will _not_ be excluded when doing the faceting.
Hence, if you e.g. restrict the range of a date field X with an axis constraint and then add a statistical facet for getting the min value in X, the returned value will only take into account the values that fall in the allowed range.

The returned result contain the field name, the specified `statName`, and the result in the property `stringStatResult`

```json
"facets": [
  {
    "type": "stringStat",
    "axis": "publicationDate",
    "statName": "min_publicationDate",
    "stringStatResult": "2014-03-16T17:33:18Z"
  }
]
```

### Highlighting: `highlightFields` and `autoHighlight`

MEx can highlight the matched terms to allow the client to understand the match context.
Highlights are returned as snippets with inserted highlighting tags around search terms of the user query.
Snippets are not returned in a particular order.
The terms are tagged with the unicode characters `\ue000` before and `\ue001` after the term, respectively.
These unicode code points from the 'Private Use Area' are used to allow unambiguous identification in front end.

`highlightFields` specifies an array of fields for which snippets will be returned.

If `autoHighlight` is set to `true`, snippets will be returned for the fields belonging to the current search focus (hence the focus called "default" if none is given explicitly in the `searchFocus` property).

```json
"highlightFields": [
"title",
"abstract"
],
"autoHighlight": false
```

If `highlightFields` and `autoHighlight` are both omitted, no highlight snippets are returned.
If both parameters are set, `highlightFields` is ignored.
The format of the returned snippets is described below.

### Fuzzy search and prefix search settings: `maxEditDistance` and `useNgramField`

Fuzzy search allows search results to be retrieved even if a query term does not exactly match the term in the search index.
`maxEditDistance` defines by how many edits (addition/deletion of a character or swap of two neighboring characters) a query and an indexed term can _at most_ differ while still allowing a match.
The only allowed values are 0, 1, and 2.
_NOTE: Setting the max edit distance to 2 is discouraged as it significantly slows down the search_.However, how many differences are allowed also depends on the _length_ of a word.
Query terms with three characters or fewer must be matched exactly; for query terms with 7 or more, an edit distance of `maxEditDistance` is applied.
For words with lengths between these two boundaries, the effective edit increases linearly (with rounding) between 0 and the max.
This ensures that we do not get too many spurious matches due to fuzzy matching of short terms, but still allow typos, spelling differences etc. in longer words.
For instance, with `maxEditDistance = 2` searching for `raw` would require an exact match while searching for `roaming` would  match e.g. both `roamjng` and `froamin`
The edit distance is applied to each term individually, i.e. each term can deviate from the matched terms by `maxEditDistance`.
The edit distance is not applied to phrase searches.
To disable fuzzy search for a specific term, a quoted term can be used.
If `maxEditDistance` is omitted, fuzzy search is disabled (same as setting `maxEditDistance` to 0).

The Boolean flag `useNgramField` (default: false) controls whether we should also query indexed prefixes.
Turning this on will enable matches on indexed words that start with the query term, regardless of the edit-distance.
For instance, if we query for `friend`, we would find `friendly`, `friendliness` etc.
Note that we only search for matches that start with the _full_ query term (prefixes are generated only for the stored data, not for the query terms).
Hence, if we query for `friends`, we would not find `friendliness` since the former is not an exact prefix of latter (although they share the prefix `friend`).
Only prefixes consisting of 5 characters or more will be matched, ensuring that we do not get spurious matches from short words.
`useNgramField` can be combined with a non-zero edit distance to allow typos in prefix matches.
Indeed, the lower bound on the prefix length is set to be optimal in combination with fuzzy matching (e.g. allowing prefix matches of length four)

Example values (these are the recommended settings):

```json
"maxEditDistance": 1,
"useNgramField": true
```

## Search query response format

The below is a possible response to the query given at the top of this item.

```json
{
  "numFound": 2,
  "numFoundExact": true,
  "start": 0,
  "maxScore": 0,
  "diagnostics": {
    "parsingSucceeded": true,
    "parsingErrors": [],
    "cleanedQuery": "disease | influence",
    "queryWasCleaned": false
  },
  "items": [
    {
      "itemId": "022f1379-f4cd-4326-91ee-72ab8373f534",
      "entityType": "",
      "values": [
        {
          "fieldName": "author",
          "fieldValue": "Carol Greene"
        },
        {
          "fieldName": "author",
          "fieldValue": "Carrie McCracken"
        },
        {
          "fieldName": "author",
          "fieldValue": "Cynthia Bearer"
        },
        {
          "fieldName": "author",
          "fieldValue": "Dustin Olley"
        },
        {
          "fieldName": "author",
          "fieldValue": "J Allen Baron"
        },
        {
          "fieldName": "author",
          "fieldValue": "Victor Felix"
        },
        {
          "fieldName": "title",
          "fieldValue": "The Human Disease Ontology 2022 update.",
          "language": "en"
        }
      ]
    },
    {
      "itemId": "a8ad4b55-1b84-41d2-a585-e22c302c6542",
      "entityType": "",
      "values": [
        {
          "fieldName": "author",
          "fieldValue": "Michael P Fisher"
        },
        {
          "fieldName": "title",
          "fieldValue": "Politicized disease surveillance: A theoretical lens for understanding sociopolitical influence on the surveillance of disease epidemics.",
          "language": "en"
        }
      ]
    }
  ],
  "facets": [
    {
      "type": "exact",
      "axis": "author",
      "bucketNo": 16,
      "buckets": [
        {
          "value": "Carol Greene",
          "count": 1
        },
        {
          "value": "Carrie McCracken",
          "count": 1
        },
        {
          "value": "Cynthia Bearer",
          "count": 1
        }
      ]
    }
  ],
  "highlights": [
    {
      "itemId": "022f1379-f4cd-4326-91ee-72ab8373f534",
      "matches": [
        {
          "fieldName": "title",
          "language": "en",
          "snippets": ["The Human \ue000Disease\ue001 Ontology 2022 update."]
        }
      ]
    },
    {
      "itemId": "a8ad4b55-1b84-41d2-a585-e22c302c6542",
      "matches": [
        {
          "fieldName": "title",
          "language": "en",
          "snippets": [
            "Politicized \ue000disease\ue001 surveillance: A theoretical lens for understanding sociopolitical \ue000influence\ue001 on"
          ]
        }
      ]
    }
  ]
}
```

Note that text field entries in matching items and highlights are assigned language tags if such tags where assigned in the original data _and_ the assigned language is supported (for search purposes, an unknown language is treated like an unspecified language

### Returned diagnostic information: `diagnostics`

This element is used to return further information about how the query was handled (assuming that the query could be carried out).
The field `cleanedQuery` contains a cleaned query string that indicates how the MEx parser understood the logical structure of the user search query without potential dangling operators (no Solr-specific details like sanitization is included).
For instance, if the user inputs `-hello ) | there+here`, the cleaned query would be `-hello | (there + here)`.
If the user query was cleaned (i.e. if the executed query differed from the original user input), `queryWasCleaned` will be set to `true`.
If the input was _not_ a valid MEx query, the Boolean flag `parsingSucceeded` will be set to `false` and the string array field `parsingErrors` will contain error messages for the corresponding parsing errors (if available).

If the overall MEx system is configured to use relaxed error handling, some errors inside the query engine will be ignored and will not cause an overall search error.
This is the case for specific kind of errors (such as facets on unknown axes or certain configuration inconsistencies) for which the query can still be carried out, but where parts of the response (e.g. a specific facet) may be missing or empty.
In such cases, the `diagnostics` object also contains a string array in the property `ignoredErrors` which lists codes for the errors that were ignored.
For systems using the relaxed error handling, clients must therefore check for this field to understand whether the returned search result may diverge from the standard format.
If using the strict error policy, the `ignoredErrors` property should always be absent or empty, as any error otherwise reported here should lead to a 500 error without a response body.
