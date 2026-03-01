package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"os/exec"
)

// Executor выполняет Python скрипты и Jupyter Notebooks
type Executor struct{}

// NewExecutor создает новый исполнитель
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute выполняет файл скрипта по расширению: .py или .ipynb
func (e *Executor) Execute(filePath string, timeout int) JobResult {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".ipynb":
		return e.ExecuteNotebook(filePath, timeout)
	case ".py":
		return e.ExecutePythonScript(filePath, timeout)
	default:
		return JobResult{
			Success: false,
			Error:   fmt.Sprintf("неподдерживаемый формат: %s (поддерживаются .py и .ipynb)", ext),
		}
	}
}

// ExecutePythonScript выполняет Python скрипт с заданным таймаутом
func (e *Executor) ExecutePythonScript(filePath string, timeout int) JobResult {
	if timeout == 0 {
		timeout = 300
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", filePath)
	output, err := cmd.CombinedOutput()

	result := JobResult{
		Success: err == nil,
		Output:  string(output),
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}

// ipynb структуры для парсинга Jupyter Notebook (только нужные поля)
type ipynb struct {
	Cells []ipynbCell `json:"cells"`
}

type ipynbCell struct {
	Outputs []ipynbOutput `json:"outputs"`
}

type ipynbOutput struct {
	OutputType string   `json:"output_type"`
	Name       string   `json:"name"`
	Text       []string `json:"text"`
}

// ExecuteNotebook выполняет Jupyter Notebook с заданным таймаутом
func (e *Executor) ExecuteNotebook(filePath string, timeout int) JobResult {
	if timeout == 0 {
		timeout = 300
	}

	dir := filepath.Dir(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	outputPath := filepath.Join(dir, base+"_executed.ipynb")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// python -m jupyter nbconvert ... — использует jupyter из того же venv, что и PYTHON_PATH
	// cmd.Dir = dir, поэтому --output-dir "." — пишем в текущую (dir) директорию процесса
	cmd := exec.CommandContext(ctx, "python3", "-m", "jupyter", "nbconvert",
		"--to", "notebook",
		"--execute",
		fmt.Sprintf("--ExecutePreprocessor.timeout=%d", timeout),
		"--output", filepath.Base(outputPath),
		"--output-dir", ".",
		filepath.Base(filePath),
	)
	cmd.Dir = dir
	nbconvertOut, err := cmd.CombinedOutput()

	// Удаляем временный исполненный файл при выходе
	defer os.Remove(outputPath)

	if err != nil {
		return JobResult{
			Success: false,
			Output:  string(nbconvertOut),
			Error:   err.Error(),
		}
	}

	// Парсим исполненный notebook и извлекаем вывод ячеек
	output, parseErr := e.extractNotebookOutput(outputPath)
	if parseErr != nil {
		return JobResult{
			Success: true,
			Output:  string(nbconvertOut) + "\n[Вывод ячеек не удалось извлечь: " + parseErr.Error() + "]",
		}
	}

	return JobResult{
		Success: true,
		Output:  strings.TrimSpace(output),
	}
}

// extractNotebookOutput читает исполненный .ipynb и возвращает вывод только последней ячейки
func (e *Executor) extractNotebookOutput(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var nb ipynb
	if err := json.Unmarshal(data, &nb); err != nil {
		return "", err
	}

	if len(nb.Cells) == 0 {
		return "", nil
	}

	// Берём только последнюю ячейку
	lastCell := nb.Cells[len(nb.Cells)-1]
	var out strings.Builder
	for _, o := range lastCell.Outputs {
		if o.OutputType != "stream" {
			continue
		}
		for _, line := range o.Text {
			out.WriteString(line)
		}
	}
	return out.String(), nil
}
