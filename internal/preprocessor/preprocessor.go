// Package preprocessor предоставляет шаги предобработки пользовательского ввода
// перед передачей строк в парсер.
package preprocessor

import (
	"regexp"
)

// Result описывает состояние строки после выполнения всех шагов препроцессинга.
type Result struct {
	Original string
	Value    string
}

// Step описывает отдельный шаг препроцессинга.
type Step interface {
	Apply(input Result) (Result, error)
}

// Preprocessor выполняет последовательность шагов обработки пользовательского ввода.
type Preprocessor struct {
	steps []Step
}

// NewPreprocessor создает препроцессор с указанными шагами.
func NewPreprocessor(steps ...Step) *Preprocessor {
	return &Preprocessor{steps: steps}
}

// Process последовательно применяет шаги к исходной строке.
func (p *Preprocessor) Process(input string) (Result, error) {
	result := Result{
		Original: input,
		Value:    input,
	}

	var err error
	for _, step := range p.steps {
		result, err = step.Apply(result)
		if err != nil {
			return Result{}, err
		}
	}

	return result, nil
}

// EnvSubstitutionStep выполняет подстановку переменных окружения.
type EnvSubstitutionStep struct {
	Env map[string]string
}

var (
	bracedPattern  = regexp.MustCompile(`\$\{([^}]+)\}`)
	defaultPattern = regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
)

// Apply реализует шаг подстановки переменных окружения.
func (s *EnvSubstitutionStep) Apply(input Result) (Result, error) {
	value := bracedPattern.ReplaceAllStringFunc(input.Value, func(match string) string {
		varName := match[2 : len(match)-1]
		if envValue, ok := s.Env[varName]; ok {
			return envValue
		}
		return match
	})

	value = defaultPattern.ReplaceAllStringFunc(value, func(match string) string {
		varName := match[1:]
		if envValue, ok := s.Env[varName]; ok {
			return envValue
		}
		return match
	})

	return Result{
		Original: input.Original,
		Value:    value,
	}, nil
}
