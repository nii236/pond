package bot

import "fmt"

type StdOutWriter struct {
}

func (w *StdOutWriter) Write(b []byte) (int, error) {
	fmt.Println(b)
	return len(b), nil
}

type SlackWriter struct {
}

func (w *SlackWriter) Write(b []byte) (int, error) {
	return 0, nil
}

type MattermostWriter struct {
}

func (w *MattermostWriter) Write(b []byte) (int, error) {
	return 0, nil
}
