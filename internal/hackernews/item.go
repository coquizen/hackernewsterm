package hackernews

// Item describes an interface of the generic fetchCanonicalItem item. It contains all the fields used by the various posting types (each with overlapping fields)
type Item interface {
	By() string
	ID() int
	Kids() []int
	Parent() int
	Parts() []int
	Score() int
	Text() string
	Time() int
	Title() string
	Type() string
	URL() string
	Deleted() bool
	Dead() bool
	Descendants() int
}

// item implements Item interface for convenient conversion to the different posting types.
type item map[string]interface{}

// By returns the author of said item.
func (i item) By() string {
	str, _ := i["by"].(string)
	return str
}

// ID returns the item's unique ID for convenient reference
func (i item) ID() int {
	num, _ := i["id"].(int)
	return num
}

// Kids returns a list of ID's the item's comments
func (i item) Kids() []int {
	kids := i["kids"]
	iKids := kids.([]interface{})
	parsedArray := make([]int, len(iKids))
	for index, kid := range iKids {
		parsedArray[index] = kid.(int)
	}
	return parsedArray
}

// Parent returns the ID of the comment's parent or relevant story
func (i item) Parent() int {
	parentID, _ := i["parent"].(int)
	return parentID
}

// Parts returns a list of related pollopts in display order
func (i item) Parts() []int {
	parts := i["parts"]
	iParts := parts.([]interface{})
	parsedArray := make([]int, len(iParts))
	for index, part := range iParts {
		parsedArray[index] = part.(int)
	}
	return parsedArray
}

// Score returns the story's, poll's, or job's score
func (i item) Score() int {
	score, _ := i["score"].(int)
	return score
}

// Text returns the the comment,story, or polltext in HTML
func (i item) Text() string {
	text, _ := i["text"].(string)
	return text
}

// Time return the creation date of said item in Unix Time
func (i item) Time() int {
	time, _ := i["time"].(int)
	return time
}

// Title returns the story's, poll's , or job's title in HTML
func (i item) Title() string {
	title, _ := i["title"].(string)
	return title
}

// Type returns the item's type: one of "job", "story", "comment", "poll", or "pollopt"
func (i item) Type() string {
	iType, _ := i["type"].(string)
	return iType
}

func (i item) URL() string {
	addr, _ := i["url"].(string)
	return addr
}

func (i item) Descendants() int {
	numOfDescendants, _ := i["descendants"].(int)
	return numOfDescendants
}

func (i item) Deleted() bool {
	s, _ := i["deleted"].(bool)
	return s
}

func (i item) Dead() bool {
	s, _ := i["dead"].(bool)
	return s
}