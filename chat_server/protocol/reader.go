package protocol

import (
	"bufio"
	"errors"
	"io"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read() (interface{}, error) {
	commandName, err := r.reader.ReadString(' ')
	if err != nil {
		return nil, err
	}

	switch commandName {
	case "MESSAGE ":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		message, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return MessageCommand{
			Name:    user[:len(user)-1],
			Message: message[:len(message)-1],
		}, nil
	case "SEND ":
		message, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return SendCommand{
			Message: message[:len(message)-1],
		}, nil

	case "NAME: ":
		name, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return NameCommand{
			Name: name[:len(name)-1],
		}, nil
	}

	return nil, errors.New("Unknown Command")
}
