package preprocessor

import "testing"

func TestEnvSubstitutionStep(t *testing.T) {
	env := map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin",
		"USER": "tester",
	}

	step := &EnvSubstitutionStep{Env: env}
	pre := NewPreprocessor(step)

	result, err := pre.Process("echo $HOME ${PATH} $UNDEFINED $USER")
	if err != nil {
		t.Fatalf("ожидался успех, получили ошибку: %v", err)
	}

	expected := "echo /home/user /usr/bin $UNDEFINED tester"
	if result.Value != expected {
		t.Fatalf("ожидалось %q, получили %q", expected, result.Value)
	}
}

func TestPreprocessor_NoSteps(t *testing.T) {
	pre := NewPreprocessor()

	result, err := pre.Process("echo test")
	if err != nil {
		t.Fatalf("ошибка при обработке без шагов: %v", err)
	}

	if result.Value != "echo test" {
		t.Fatalf("значение должно совпадать: %s", result.Value)
	}
}
