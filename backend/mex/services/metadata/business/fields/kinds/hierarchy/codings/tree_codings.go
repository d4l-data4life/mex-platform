package codings

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/d4l-data4life/mex/mex/shared/coll/forest"
)

type PostgresCodings struct {
	DB *pgxpool.Pool

	cache map[string]Hull
}

// Note: language is hard-coded to 'en' for resolve until we fixed the API calls for it.
// Get all items (latest version) of the required type, along with the ID of their parent linked through the specified field & display information.
const sqlTreeNodes = `
select
	liwbi.business_id,
	coalesce(iwl.target_business_id, '') as parent_id,
	coalesce(civ.field_value, '') as display,
	coalesce(civ."language", 'en') as language
from
	latest_items_with_business_id liwbi
left outer join items_with_links(array['%s']::TEXT[]) iwl
	on iwl.source_business_id = liwbi.business_id
left outer join current_item_values civ
	on civ.item_id = liwbi.item_id and civ.field_name = $2
where
	liwbi.entity_name = $1
`

// Language -> Coding
type Hull = map[string][]Coding

type myNode struct {
	id       string
	display  string
	language string
	parentID string
}

func NewPostgresCodings(db *pgxpool.Pool) PostgresCodings {
	return PostgresCodings{DB: db, cache: make(map[string]Hull)}
}

func (codings *PostgresCodings) Reset() {
	codings.cache = make(map[string]Hull)
}

func (codings *PostgresCodings) TransitiveClosure(ctx context.Context, entityType string, linkType string, displayFieldName string, code string) (map[string][]Coding, error) {
	key := fmt.Sprintf("%s-%s-%s-%s", entityType, linkType, displayFieldName, code)
	if h, ok := codings.cache[key]; ok {
		return h, nil
	}

	rows, err := codings.DB.Query(ctx, fmt.Sprintf(sqlTreeNodes, linkType), entityType, displayFieldName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	f := forest.NewForestWriter[myNode]()
	for rows.Next() {
		var n myNode
		if err := rows.Scan(&n.id, &n.parentID, &n.display, &n.language); err != nil {
			return nil, fmt.Errorf("could not compute transitive full: %w", err)
		}

		f.Add(n.id, n.parentID, n)
	}

	r := f.Seal()
	if r == nil {
		return nil, fmt.Errorf("invalid forest")
	}

	// Determine path to root
	hulls := make(Hull)
	path, err := r.RootPath(code)
	if err != nil {
		return nil, err
	}
	path = append([]string{code}, path...)

	for _, id := range path {
		payload, err := r.GetByID(id)
		if err != nil {
			return nil, err
		}

		c := Coding{
			Code:       payload.id,
			Display:    payload.display,
			Language:   payload.language,
			Depth:      r.MustDepth(payload.id),
			ParentCode: r.MustParent(payload.id),
		}

		bucket, ok := hulls[payload.language]
		if ok {
			hulls[payload.language] = append(bucket, c)
		} else {
			hulls[payload.language] = []Coding{c}
		}
	}

	codings.cache[key] = hulls

	return hulls, nil
}

// Note: language is hard-coded to 'en' for resolve until we fixed the API calls for it.
const sqlResolve = `select '' as parent_code, coalesce(civ.field_value, '') as display, 0 as depth
from latest_items_with_business_id liwbi
left outer join current_item_values civ
  on civ.item_id = liwbi.item_id and civ.field_name = $1 and coalesce(civ.language, 'en') = 'en'
where liwbi.business_id = $2
limit 1
`

func (codings *PostgresCodings) Resolve(ctx context.Context, codeSystemName string, linkType string, displayFieldName string, code string) ([]Coding, error) {
	rows, err := codings.DB.Query(ctx, sqlResolve, displayFieldName, code)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	codes := []Coding{}
	for rows.Next() {
		var c Coding
		if err := rows.Scan(&c.ParentCode, &c.Display, &c.Depth); err != nil {
			return nil, err
		}

		codes = append(codes, c)
	}

	return codes, nil
}

func (codings *PostgresCodings) GetCodeSystemNames(ctx context.Context) ([]string, error) {
	return []string{"<dynamic>"}, nil
}

func (codings *PostgresCodings) Close() error {
	codings.DB.Close()
	return nil
}
