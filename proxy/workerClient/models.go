package workerClient

type SummRequest struct {
	Numbers []int `form:"numbers"`
}

type SummRepsponse struct {
	Result int `form:"result"`
}

type SummRepsponseError struct {
	Error string `form:"error"`
}
