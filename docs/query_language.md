# The MEx query language and search logic

## The MEx query language

MEx offers a relatively simply query language that allows user to do more structured searches on the MEx platform using wildcards and Boolean operators.
The formal definition of the language is given by the ANTLR 4 grammar in the file _backend/mex/services/query/parser/MexQueryGrammar.g4_.
A search query consists of any list of basic search terms, quoted phrases, and Boolean expressions using the supported operators (with grouping provided by parentheses).

### Simple search terms

A basic search term is any string which

- does not contain space characters, and
- does not contain MEx special symbols other than the wildcard `*` (see below), and
- is not enclosed in double quotes

The simplest query is a single word, e.g. `hand`.
If multiple search terms separated by spaces are given, only items matching all terms will be returned - that is, the terms are combined using the logical AND operator.
For instance, the query `hello goodbye` triggers a search for items containing both _hello_ and _goodbye_.
An OR search requires the explicit use ot the OR-operator (see below).

If the search query is empty (or contains only spaces), no search constraint is imposed and all indexed items are returned (or, more precisely: all items with entries in at least one of the fields being searched).

The wildcard character `*` matches zero or more characters at position where it occurs.
For instance, `co*d` matches e.g. "cod", "cord", and "covid".
If an asterisk should _not_ be interpreted as a wildcard match, it should be escaped with a backslash (`\`).
Some example of search terms would be `hand`, `ha*d`, and `escaped\(example\"`.
Escape symbols in front of non-special characters (including spaces) are discarded.

### Search phrases (double-quoted terms)

Double-quoted terms are treated as phrases during search, i.e. items only match if they contain exactly this phrase (including whitespace).
For instance, `"hand pain"` triggers a search for items containing the phrase _hand pain_, not just _hand_ and _pain_ separately.
The non-whitespace strings in a quoted phrase must be valid search terms.
Quote characters can be escaped if they should not be treated a phrase groupers, e.g. the input `"the \"funny\" joke"` triggers a search for items containing the single phrase _the "funny" joke_.

MEx search operators (see below) inside phrases are treated as normal characters, i.e. `"hand | foot"` triggers a search for the phrase _hand | foot_, not an OR-search.

### Boolean and grouping operators

The following operators that combine or modify individual search terms are supported.

| Operator | Meaning                                                          | Example                   | Comment                                                                                                    |
| -------- |------------------------------------------------------------------| ------------------------- |------------------------------------------------------------------------------------------------------------|
| &dash;   | NOT: term must be absent from matching item                      | &dash;hand                | Must be prefixed to the term without intervening whitespace. Can also be used in front of a quoted phrase. |
| &#43;    | AND: both terms must be present in matching record               | hand &#43; foot           | '+' is also interpreted as AND when not surrounded as spaces (as in 'hand+foot').                          |
| &#124;   | OR: at least one of the terms must be present in matching record | hand &#124; foot          | '&#124;' is also interpreted as OR when not surrounded as spaces (as in 'hand&#124;foot').                 |
| (...)    | grouping of terms                                                | head + (hand &#124; foot) | Used to ensure a particular evaluation order                                                               |

As an example, to search for items containing either (or both) of the words "hand" and "foot", but not the word, "nose", we would use the search query `(hand | foot ) + -foot`.

Note that MEx symbols that cannot be assigned to a valid MEx expression - e.g. unmatched brackets or binary Boolean operators with one argument missing - are ignored (dropped) even if they are attached to a search term (prevents breaking the query).
Hence, `hello)`, `hello"`, and `+hello` will all give the same results as the query `hello`.
The exception is a dash (NOT-operator) occurring inside or at the end of a term.
Such internal or post-fixed dashes are interpreted as part of the term, allowing search for hyphenated terms like `hard-working` without triggering a negation.

#### Operator precedence and handling of ambiguous Boolean expressions

In the absence of grouping parentheses, converting a query with multiple binary Boolean operators to a specific search requires a precedence ordering of the operators.
For instance, without such an ordering the input `hand + foot | nose` might mean `(hand AND foot) OR nose` or `hand AND (foot OR nose)`.
As per the general guiding principle, the aim is to try to capture what the user most likely meant.
MEx uses the following precedences, given here in order of _decreasing_ precedence:

1. `-` (NOT)
2. `+` (_explicit_ AND)
3. `|` (OR)
4. whitespace between terms (_implicit_ AND)

Brackets can always be used to enforce a different order of evaluation.
Below are some example to clarify the rules.

| Search query                             | Interpretation                               | Comment                                                                                                     |
| ---------------------------------------- | -------------------------------------------- |-------------------------------------------------------------------------------------------------------------|
| hand + foot &#124; nose                  | (hand AND foot) OR nose                      | AND has higher precedence than OR                                                                           |
| hand + (foot &#124; nose)                | hand AND (foot OR nose)                      | Brackets override normal precedence                                                                         |
| hand foot &#124; nose                    | hand AND (foot OR nose)                      | Implicit AND (whitespace) has _lower_ precedence than OR                                                    |
| hand + foot &#124; nose + ear            | (hand AND foot) OR (nose AND ear)            | AND has higher precedence than OR                                                                           |
| hand + (foot &#124; nose) + ear          | hand AND (foot OR nose) AND ear              | Brackets override normal precedence                                                                         |
| hand &#124; foot nose &#124; ear + mouth | (hand OR foot) AND (nose OR (ear AND mouth)) | Multiple expressions (with explicit operators) separated by space are evaluated separately and then combined |

Treating an _explicit_ AND (`+`) and an _implicit_ AND (whitespace) differently ensures that if the search string is a mixture of words/phrases and Boolean expressions, the Boolean operator precedence rules will only be imposed within the Boolean expression, not across the whole search string (cf. last example in table).

## Fuzzy search in MEx

Fuzzy search refers to matching a search term with an item term even when the two are not quite identical.
The wildcard operator (`*`) discussed above explicitly triggers a particular kind of fuzzy search.
However, MEx offers two further forms of fuzzy matching that do not require explicit triggering in the query: prefix matching and soft matching of individual terms.

### Prefix matching

Prefix matching means that a search term will match an item term that _starts_ with the search term, e.g. searching for `round` will return items containing the term `roundabout`.
Note that prefix matching is not applied to very short words, to avoid flooding the search results with matches to common prefixes.
Likewise, only prefixes up to a certain maximal length are searched.

### Soft matching of individual terms

If configured or explicitly requested in a request to the MEx search API, terms that differ by e.g. only one character are matched ("soft matching").
This allows matching even in the presence of small typos, e.g. "cosnider" will match "consider" (swapped neighboring letters is considered a single difference).
As for prefix matching, soft matching is not applied to very short words to avoid too many spurious matches.
It is not recommended to allow matching in the presence of more than a single difference: apart from possibly producing a lot of unwanted matches, it also dramatically slows down the search.

## Language-specific search behavior

Content in field of the MEx kind `text` which is explicitly classified as being in one of the supported langauges (currently English and German) is also stored in language-aware fields in the search index.
This allows optimizing the search in language-specific ways.

### Stemming

Stemming is the process of reducing words to a simpler forms using simple heuristic rules that are specific to the (natural) language in question.
It does not aim to capture the full complexity of the language, but merely aims to reduce related words to the same basic form in the most common cases.
For instance, an English stemmer would reduce both of the words "runs" and "running" to the basic form "run".
In language-aware fields, stemming is applied to both stored texts and search queries, in order to allow matches of related words.
Thus, following the above example, querying for "runs" would return a record containing "running" in a stemmed field since stemming reduces both to the basic form "run".
Since it allows matching terms that are not identical, stemming is a specific form of fuzzy matching.

### Stopwords

In every language, certain short words occur very often and largely independently of the topic of the text.
Examples from English would be e.g. "a:, "an", "to", and "on".
Because of their ubiquity and lack of specificity, matches on such terms bring little values when it comes to returning the most relevant items to the user.
To avoid skewing the search result in irrelevant ways, the  MEx search engine therefore ignores a language-specific set of such words (the so-called _stopwords_) when searching fields in a supported language.
This will typically primarily affect the ranking of documents, but may also affect the set of search results e.g. if the user issues a query consisting of only stopwords.

## Fields targeted when querying

As discussed in the detail in the [documentation of the search configuration](metadata_config.md), a search is always directed at a specific set of fields in the underlying search index.
Such a set is called a _search focus_ in MEx and typically represents a particular theme that the client might want to focus on in a specific context (general, main descriptive labels, people and organizations etc.).
A search focus called "default" must always be defined; it is used if the client does not explicitly provide a search focus in a search query.
Corresponding to each search focus is a set of backing fields in the search index, each of which is optimized to support a particular kind of matching:

1. Fields optimized for fuzzy matching text in a specific language
2. Fields optimized for fuzzy matching text in an unspecified language
3. Fields optimized for matching prefixes of stored text
4. Fields optimized for exact matching of text

Exactly which combination of these fields should be searched for a given query differs not only from query to query, but also between individual parts of a single query.
Consider, for instance, the query `"back pain" + swimming`.
The phrase `"back pain"` should only match if this exact phrase is found, so we should only look for matches in the search focus backing fields optimized for exact matching (no. 4 in the list above).
On the other hand, for the non-phrase term `swimming`, we expect the different kinds of fuzzy matching to apply: prefix (to match e.g. "swimmingly"), English stemming (to match e.g. "swim"), and generic fuzzy matching (to match e.g. an entry with a typo like "swiming").
However, we also expect an exact match to score higher.

The MEx query engine handles this by defining two so-called _matching operators_, one for single (non-phrase) terms and one for phrases.
The matching operator is a sub-query that specifies exactly which backing fields should be targeted when searching for a term or phrase, respectively.
When processing a MEx query, each term or phrase is replaced with the corresponding matching operator.
For phrases, the matching operator for a given phrase is simply a search for the phrase in the backing field for exact matching.
For non-phrase terms, the matching operator is an OR-combination of search in all the backing fields.

Thus, the query `"back pain" + swimming` effectively gets transformed into a query that has the following structure

```
("back pain" HAS MATCH IN ExactMatchField)
AND
(
  ("swimming" HAS MATCH IN EnglishLanguageField)
  OR
  ("swimming" HAS MATCH IN GermanLanguageField)
  OR
  ("swimming" HAS MATCH IN GenericLanguageField)
  OR
  ("swimming" HAS MATCH IN PrefixField)
  OR
  ("swimming" HAS MATCH IN ExactField)
)
```

The logical structure above fixes which items match the query.
However, we also expect these matches to be ranked in a sensible way so that the best matches can be returned first.
For instance, we would expect an item which contains the _exact_ term "swimming" to be ranked higher than an item that only has a fuzzy match (e.g. "swimmingly").
To configure this, the matching operators assign a weight (often referred to as a _boost factor_) to each included field.
This weight indicates how strongly a match in this field should affect the match ranking.
Thus, in the non-phrase term matching operator, the exact match field is assigned a higher weight than the other fields to boost exact matches in the match ranking.
