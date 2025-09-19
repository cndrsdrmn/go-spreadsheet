package spreadsheet

import (
	"bufio"
	"encoding/json"
	"os"
)

type Reader interface {
	BatchStream(*os.File) (<-chan Worksheet, error)
	Close() error
	Options() any
	Read(*os.File) (Worksheet, error)
}

func ExecuteReader(r Reader, file *os.File, stream bool) error {
	ch, err := getWorksheetChannel(r, file, stream)
	if err != nil {
		return err
	}

	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()

	for ws := range ch {
		if ws.Err != nil {
			return ws.Err
		}
		if err := json.NewEncoder(buf).Encode(ws); err != nil {
			return err
		}
		if err := buf.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func getWorksheetChannel(r Reader, file *os.File, stream bool) (<-chan Worksheet, error) {
	if stream {
		return r.BatchStream(file)
	}

	ws, err := r.Read(file)
	if err != nil {
		return nil, err
	}

	c := make(chan Worksheet, 1)
	c <- ws
	close(c)

	return c, nil
}
