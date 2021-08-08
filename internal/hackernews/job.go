package hackernews

type Job struct {
	By    string
	ID    int
	Score int
	Time  int
	Title string
	Type  string
	URL   string
}

func (i item) ToJob() *Job {
	return &Job{
		i.By(),
		i.ID(),
		i.Score(),
		i.Time(),
		i.Title(),
		i.Type(),
		i.URL(),
	}
}
