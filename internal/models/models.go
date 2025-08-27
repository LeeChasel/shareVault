package models

func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&File{},
	}
}
