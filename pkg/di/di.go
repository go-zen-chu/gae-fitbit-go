package di

type DI map[string]interface{}

func NewDI() DI {
	return map[string]interface{}{}
}
