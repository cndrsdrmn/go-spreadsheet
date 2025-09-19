package spreadsheet

type Sheet struct {
	Index   int      `json:"index"`
	Name    string   `json:"name"`
	Rows    *[]Row   `json:"rows,omitempty"`
	Batches *[]Batch `json:"batches,omitempty"`
}

type Batch struct {
	RowStart int   `json:"row_start"`
	RowEnd   int   `json:"row_end"`
	Rows     []Row `json:"rows"`
}

type Row struct {
	Sheet int      `json:"sheet"`
	Index int      `json:"index"`
	Cells []string `json:"cells"`
}

type Worksheet struct {
	Sheets []Sheet `json:"sheets"`
	Err    error   `json:"error,omitempty"`
}

func CreateBatch(start, end int, rows []Row) Batch {
	return Batch{
		RowStart: start,
		RowEnd:   end,
		Rows:     rows,
	}
}

func CreateSheet(index int, name string, rows []Row, batches []Batch) Sheet {
	var rowPtr *[]Row
	if len(rows) > 0 {
		rowPtr = &rows
	}

	var batchPtr *[]Batch
	if len(batches) > 0 {
		batchPtr = &batches
	}

	return Sheet{
		Index:   index,
		Name:    name,
		Rows:    rowPtr,
		Batches: batchPtr,
	}
}

func WrapWorksheet(sheets ...Sheet) Worksheet {
	return Worksheet{
		Sheets: sheets,
	}
}
