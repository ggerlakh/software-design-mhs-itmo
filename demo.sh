#!/bin/bash

# Скрипт для демонстрации возможностей go-cli
# Использование: ./demo.sh

echo "=== Демонстрация возможностей go-cli ==="
echo

# Проверяем, что проект собран
if [ ! -f "./go-cli" ]; then
    echo "Сборка проекта..."
    go build -o go-cli ./cmd/go-cli
    if [ $? -ne 0 ]; then
        echo "Ошибка сборки!"
        exit 1
    fi
fi
echo "✓ Проект готов"
echo

echo "🎯 Демонстрация 1: Базовые команды"
echo "----------------------------------------"
echo "Команда: echo 'Привет, мир!'"
echo "echo 'Привет, мир!'" | ./go-cli
echo

echo "Команда: pwd"
echo "pwd" | ./go-cli
echo

echo "🎯 Демонстрация 2: Работа с пайпами"
echo "----------------------------------------"
echo "Команда: echo 'hello world' | wc"
echo "echo 'hello world' | wc" | ./go-cli
echo

echo "Команда: echo 'hello world' | wc -w (только слова)"
echo "echo 'hello world' | wc -w" | ./go-cli
echo

echo "Команда: pwd | wc (подсчет символов в пути)"
echo "pwd | wc" | ./go-cli
echo

echo "🎯 Демонстрация 3: Подстановка переменных окружения"
echo "----------------------------------------"
echo "Команда: echo \$HOME"
echo 'echo $HOME' | ./go-cli
echo

echo "Команда: echo \${PATH}"
echo 'echo ${PATH}' | ./go-cli
echo

echo "Команда: echo Пользователь \$USER живет в \$HOME"
echo 'echo Пользователь $USER живет в $HOME' | ./go-cli
echo

echo "Команда: echo \$UNDEFINED (несуществующая переменная)"
echo 'echo $UNDEFINED' | ./go-cli
echo

echo "🎯 Демонстрация 4: Работа с файлами"
echo "----------------------------------------"
# Создаем тестовый файл
echo "Создание тестового файла..."
cat > demo_file.txt << EOF
Строка 1
Строка 2
Строка 3
EOF

echo "Команда: cat demo_file.txt"
echo "cat demo_file.txt" | ./go-cli
echo

echo "Команда: cat demo_file.txt | wc"
echo "cat demo_file.txt | wc" | ./go-cli
echo

echo "🎯 Демонстрация 5: Интерактивный режим"
echo "----------------------------------------"
echo "Запуск интерактивного режима с несколькими командами:"
cat > interactive_demo.txt << EOF
echo Добро пожаловать в go-cli!
pwd
echo Текущая директория выше
echo hello world | wc
exit
EOF

echo "Выполнение команд из файла:"
cat interactive_demo.txt
echo
echo "Результат:"
./go-cli < interactive_demo.txt
echo

# Очистка
rm -f demo_file.txt interactive_demo.txt

echo "🎉 Демонстрация завершена!"
echo "=========================="
echo "go-cli поддерживает:"
echo "✓ Базовые команды (echo, pwd, cat, wc)"
echo "✓ Пайпы (command1 | command2)"
echo "✓ Подстановку переменных окружения (\$VAR, \${VAR})"
echo "✓ Интерактивный режим"
echo "✓ Обработку ошибок"
echo "✓ Высокую производительность"
echo
echo "Для запуска интерактивного режима:"
echo "  ./go-cli"
echo
echo "Для выполнения команд из stdin:"
echo "  echo 'pwd' | ./go-cli"
echo "  echo 'echo hello | wc' | ./go-cli"
