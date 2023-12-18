package items

func (item *Item) Schema() []string {
	schema := make([]string, 0)
	dict := make(map[string]struct{})

	for _, v := range item.Values {
		if _, ok := dict[v.FieldName]; !ok {
			schema = append(schema, v.FieldName)
			dict[v.FieldName] = struct{}{}
		}
	}

	return schema
}
