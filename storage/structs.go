package Storage


type StorageStruct struct {
	basepath string
}

type Task struct {
	id int
	date string
	title string
	comment string
	repeat int32
}