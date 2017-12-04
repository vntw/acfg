package log

type Container struct {
	InstanceUuid string `json:"instanceUuid"`
	Time         int64  `json:"time"`
	Files        []File `json:"logFiles"`
}

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewLogFile(name string) *File {
	return &File{Name: name}
}

func NewLogContainer(instanceUuid string, time int64) *Container {
	return &Container{
		Files:        []File{},
		InstanceUuid: instanceUuid,
		Time:         time,
	}
}
