package mask

type Mask interface {
	Mask(value interface{}, fieldType FieldType) interface{}
}

type FieldMask[C any] interface {
	Mask(C) C
}

type StringMask struct{}

func (mask *StringMask) Mask(input string) string {
	return "x"
}

type Float64Mask struct{}

func (mask *Float64Mask) Mask(input float64) float64 {
	return 0
}

type BooleanMask struct{}

func (mask *BooleanMask) Mask(input bool) bool {
	return false
}

type JSONMask struct {
	String  FieldMask[string]
	Float64 FieldMask[float64]
	Boolean FieldMask[bool]
}

func NewJSONMask() *JSONMask {
	return &JSONMask{
		String:  &StringMask{},
		Float64: &Float64Mask{},
		Boolean: &BooleanMask{},
	}
}

func (jsonMask *JSONMask) Mask(input interface{}, fieldType FieldType) interface{} {
	switch fieldType {
	case FieldTypeString:
		return jsonMask.String.Mask(input.(string))
	case FieldTypeFloat64:
		return jsonMask.Float64.Mask(input.(float64))
	case FieldTypeBool:
		return jsonMask.Boolean.Mask(input.(bool))
	}

	return input
}
