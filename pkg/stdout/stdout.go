package stdout

import (
	"fmt"

	"github.com/nii236/pond/pkg/pond"
)

type Adaptor struct {
	AdaptorIn chan *pond.Message

	BotIn  chan *pond.Message
	BotOut chan *pond.Message
}

func (s *Adaptor) Init(botIn chan *pond.Message, botOut chan *pond.Message) error {
	s.BotIn = botIn
	s.BotOut = botOut

	return nil
}

func (s *Adaptor) Run() {
	for {
		select {
		case out := <-s.BotOut:
			fmt.Println(out.Message)
			// fmt.Printf("%+v\n", out.Meta)
		case in := <-s.AdaptorIn:
			s.BotIn <- in
		}
	}
}
