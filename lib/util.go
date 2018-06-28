package lib

// Nullable types with UnmarshalYAML interface implemented

type NullableBool struct {
	Value bool
	Valid bool
}

func (x *NullableBool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var bl bool
	err := unmarshal(&bl)
	if err != nil {
		return err
	}
	x.Value = bl
	x.Valid = true
	return nil
}

type NullableString struct {
	Value string
	Valid bool
}

func (x *NullableString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	x.Value = s
	x.Valid = true
	return nil
}

type NullableUint struct {
	Value uint
	Valid bool
}

func (x *NullableUint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var u uint
	err := unmarshal(&u)
	if err != nil {
		return err
	}
	x.Value = u
	x.Valid = true
	return nil
}
