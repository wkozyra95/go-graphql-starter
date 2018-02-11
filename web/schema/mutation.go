package schema

type updatePatch struct {
	Field string
	Value any
}

type any struct {
	value interface{}
}

func (a any) string() (string, bool) {
	val, ok := a.value.(string)
	return val, ok
}

func (a any) bool() (bool, bool) {
	val, ok := a.value.(bool)
	return val, ok
}

func (a any) int() (int32, bool) {
	val, ok := a.value.(int32)
	return val, ok
}

func (a any) float() (float32, bool) {
	val, ok := a.value.(float32)
	return val, ok
}

func (a any) ImplementsGraphQLType(name string) bool {
	return name == "Any"
}

func (a *any) UnmarshalGraphQL(input interface{}) error {
	a.value = input
	return nil
}
