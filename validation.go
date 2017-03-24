package main

type ValidationError struct {
	in               Input
	CitiesCountIsNil bool
}

func (err ValidationError) Error() string {
	return "Validation error : no cities indicated !"
}

func (input *Input) Validate() error {
	var err ValidationError
	if input.CitiesCount == 0 {
		err.CitiesCountIsNil = true
		return err
	}
	return nil
}
