package nbformat

type Notebook struct {
	Metadata      Metadata `json:"metadata"`
	NBFormatMinor int      `json:"nbformat_minor"`
	NBFormat      int      `json:"nbformat"`
	Cells         []Cell   `json:"cells"`
}

type Metadata struct {
	Kernelspec   Kernelspec   `json:"kernelspec"`
	LanguageInfo LanguageInfo `json:"language_info"`
	OrigNbformat int          `json:"orig_nbformat,omitempty"`
	Title        string       `json:"title,omitempty"`
	Authors      []Author     `json:"authors,omitempty"`
}

type Kernelspec struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Language    string `json:"language,omitempty"`
}

type LanguageInfo struct {
	Name           string      `json:"name"`
	CodemirrorMode interface{} `json:"codemirror_mode,omitempty"` // Could be a string or an object
	FileExtension  string      `json:"file_extension,omitempty"`
	Mimetype       string      `json:"mimetype,omitempty"`
	PygmentsLexer  string      `json:"pygments_lexer,omitempty"`
}

type Author struct {
	Name string `json:"name"`
}

type Cell struct {
	ID             string        `json:"id,omitempty"`
	CellType       string        `json:"cell_type"` // "code"
	Metadata       *CellMetadata `json:"metadata"`
	Source         []string      `json:"source"`
	ExecutionCount *int          `json:"execution_count,omitempty"`
	Outputs        []Output      `json:"outputs,omitempty"`
}

type CellMetadata struct {
	Jupyter JupyterMetadata `json:"jupyter,omitempty"`
}

type JupyterMetadata struct {
	SourceHidden  bool `json:"source_hidden,omitempty"`
	OutputsHidden bool `json:"outputs_hidden,omitempty"`
}

type Output struct {
	OutputType     string                 `json:"output_type"`
	Name           string                 `json:"name,omitempty"`
	Text           []string               `json:"text,omitempty"`
	ExceptionName  string                 `json:"ename,omitempty"`
	Exceptionvalue string                 `json:"evalue,omitempty"`
	Traceback      []string               `json:"traceback,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
