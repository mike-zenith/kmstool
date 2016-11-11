package kmstool

type KMSItem struct {
	c string
	m map[string]string
}

func NewKMSItem(content string) *KMSItem {
	return &KMSItem{
		c: content,
	}
}

func (s *KMSItem) Content() string {
	return s.c
}

func (s *KMSItem) Parse(parser func(content string) (parseMap map[string]string, err error)) error {
	parsed, err := parser(s.c)
	if err != nil {
		return err
	}
	s.m = parsed
	return nil
}

func (s *KMSItem) RunOnParsed(command func(parsedHit string, parsedContent string) (content string, err error)) error {
	newParsed := make(map[string]string, len(s.m))
	for hit, content := range s.m {
		parsed, err := command(hit, content)
		if err != nil {
			return err
		}
		newParsed[hit] = parsed

	}
	s.m = newParsed
	return nil
}

func (s *KMSItem) Replace(replacer func(content string, parsedHit string, parsedContent string) (newContent string, err error)) error {
	var content string
	var err error
	var newContent string

	content = s.c
	for hit, key := range s.m {
		newContent, err = replacer(content, hit, key)
		if err != nil {
			return err
		}
		content = newContent
	}
	s.c = content
	return nil
}
