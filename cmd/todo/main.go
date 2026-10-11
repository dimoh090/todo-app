package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"todo-app/internal/storage"
	"todo-app/internal/todo"
)

// tasksFile — придуманное нами имя константы.
// Константа хранит значение, которое нельзя изменить во время работы программы.
// Все основные команды используют этот файл для хранения задач.
// Относительный путь означает: файл находится в текущей папке терминала.
const tasksFile = "tasks.json"

// main — точка входа: Go начинает выполнение приложения с этой функции.
func main() {
	// Передаём в run аргументы после имени программы.
	// os.Args[0] — имя программы, поэтому начинаем со второго элемента.
	err := run(os.Args[1:])

	// Запрос справки через --help не считаем ошибкой пользователя.
	// errors.Is проверяет, соответствует ли err указанной ошибке.
	if errors.Is(err, flag.ErrHelp) {
		return
	}

	// Если команда завершилась ошибкой, печатаем её и завершаем программу.
	// Код завершения 1 сообщает операционной системе о неудаче.
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

// run — наша функция обработки команд.
// args — срез строк с аргументами пользователя.
// Возвращаем error: nil при успехе либо описание ошибки.
func run(args []string) error {
	// Если пользователь не передал команду, показываем справку.
	if len(args) == 0 {
		printUsage()
		return nil
	}

	// Первый аргумент — название команды: add, list и так далее.
	command := args[0]

	// Выбираем блок кода, соответствующий команде пользователя.
	switch command {
	case "add":
		// Создаём отдельный набор флагов команды add.
		// ContinueOnError позволяет получить ошибку и обработать её самим.
		addCmd := flag.NewFlagSet("add", flag.ContinueOnError)

		// String объявляет текстовый флаг и возвращает указатель на его значение.
		// "desc" — имя флага; "" — значение по умолчанию.
		desc := addCmd.String("desc", "", "описание задачи")

		// Разбираем аргументы после названия команды.
		if err := addCmd.Parse(args[1:]); err != nil {
			return err
		}

		// NArg показывает количество аргументов, которые не были флагами.
		if addCmd.NArg() != 0 {
			return fmt.Errorf("лишние аргументы команды add: %v", addCmd.Args())
		}

		// *desc читает строку по указателю.
		// TrimSpace удаляет пробелы и переносы строк с краёв текста.
		description := strings.TrimSpace(*desc)

		// Запрещаем описание, которое пустое или состояло только из пробелов.
		if description == "" {
			return fmt.Errorf("укажите описание через --desc")
		}

		// Загружаем существующие задачи из основного JSON-файла.
		tasks, err := storage.LoadJSON(tasksFile)
		if err != nil {
			return fmt.Errorf("не удалось загрузить задачи: %w", err)
		}

		// Add создаёт новую задачу и возвращает обновлённый список.
		tasks = todo.Add(tasks, description)

		// Сохраняем обновлённый список, чтобы задача осталась после выхода.
		if err := storage.SaveJSON(tasksFile, tasks); err != nil {
			return fmt.Errorf("не удалось сохранить задачу: %w", err)
		}

		fmt.Println("Задача добавлена")
		return nil

	case "list":
		// Создаём набор флагов команды list.
		listCmd := flag.NewFlagSet("list", flag.ContinueOnError)

		// Если --filter не указан, используем значение all.
		filter := listCmd.String(
			"filter",
			"all",
			"фильтр задач: all, done, pending",
		)

		if err := listCmd.Parse(args[1:]); err != nil {
			return err
		}

		if listCmd.NArg() != 0 {
			return fmt.Errorf("лишние аргументы команды list: %v", listCmd.Args())
		}

		// && означает «И»: ошибка возникает, если значение не равно
		// ни all, ни done, ни pending.
		if *filter != "all" && *filter != "done" && *filter != "pending" {
			return fmt.Errorf("допустимые фильтры: all, done, pending")
		}

		tasks, err := storage.LoadJSON(tasksFile)
		if err != nil {
			return fmt.Errorf("не удалось загрузить задачи: %w", err)
		}

		// Оставляем только задачи, соответствующие выбранному фильтру.
		tasks = todo.List(tasks, *filter)

		// len возвращает количество элементов среза.
		if len(tasks) == 0 {
			fmt.Println("Задачи не найдены")
			return nil
		}

		// range перебирает задачи; _ означает, что индекс нам не нужен.
		for _, task := range tasks {
			// %d — целое число, %s — строка, %t — bool.
			// \n переводит вывод на следующую строку.
			fmt.Printf(
				"ID: %d | Задача: %s | Выполнена: %t\n",
				task.ID,
				task.Description,
				task.Done,
			)
		}

		return nil

	case "complete":
		// Создаём набор флагов для отметки задачи выполненной.
		completeCmd := flag.NewFlagSet("complete", flag.ContinueOnError)

		// Int объявляет целочисленный флаг и возвращает указатель *int.
		// Ноль по умолчанию поможет заметить отсутствие --id.
		id := completeCmd.Int("id", 0, "ID задачи")

		if err := completeCmd.Parse(args[1:]); err != nil {
			return err
		}

		if completeCmd.NArg() != 0 {
			return fmt.Errorf(
				"лишние аргументы команды complete: %v",
				completeCmd.Args(),
			)
		}

		// ID наших задач должен быть положительным числом.
		if *id <= 0 {
			return fmt.Errorf("укажите положительный ID через --id")
		}

		tasks, err := storage.LoadJSON(tasksFile)
		if err != nil {
			return fmt.Errorf("не удалось загрузить задачи: %w", err)
		}

		// Используем =, поскольку tasks и err уже объявлены выше.
		// Complete ищет задачу по ID и устанавливает Done = true.
		tasks, err = todo.Complete(tasks, *id)
		if err != nil {
			return err
		}

		// Изменение списка в памяти нужно записать в файл.
		if err := storage.SaveJSON(tasksFile, tasks); err != nil {
			return fmt.Errorf("не удалось сохранить изменения: %w", err)
		}

		fmt.Printf("Задача с ID %d отмечена выполненной\n", *id)
		return nil

	case "delete":
		// Создаём набор флагов для удаления задачи.
		deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)

		id := deleteCmd.Int("id", 0, "ID задачи")

		if err := deleteCmd.Parse(args[1:]); err != nil {
			return err
		}

		if deleteCmd.NArg() != 0 {
			return fmt.Errorf(
				"лишние аргументы команды delete: %v",
				deleteCmd.Args(),
			)
		}

		if *id <= 0 {
			return fmt.Errorf("укажите положительный ID через --id")
		}

		tasks, err := storage.LoadJSON(tasksFile)
		if err != nil {
			return fmt.Errorf("не удалось загрузить задачи: %w", err)
		}

		// Delete возвращает список без задачи с указанным ID.
		tasks, err = todo.Delete(tasks, *id)
		if err != nil {
			return err
		}

		if err := storage.SaveJSON(tasksFile, tasks); err != nil {
			return fmt.Errorf("не удалось сохранить изменения: %w", err)
		}

		fmt.Printf("Задача с ID %d удалена\n", *id)
		return nil

	case "export":
		// Экспорт сохраняет копию списка в выбранном формате.
		exportCmd := flag.NewFlagSet("export", flag.ContinueOnError)

		format := exportCmd.String("format", "", "формат: json или csv")
		out := exportCmd.String("out", "", "путь к выходному файлу")

		if err := exportCmd.Parse(args[1:]); err != nil {
			return err
		}

		if exportCmd.NArg() != 0 {
			return fmt.Errorf(
				"лишние аргументы команды export: %v",
				exportCmd.Args(),
			)
		}

		// ToLower переводит буквы в нижний регистр.
		// Благодаря этому JSON и json понимаются одинаково.
		exportFormat := strings.ToLower(strings.TrimSpace(*format))
		outputPath := strings.TrimSpace(*out)

		if exportFormat != "json" && exportFormat != "csv" {
			return fmt.Errorf("укажите --format=json или --format=csv")
		}

		if outputPath == "" {
			return fmt.Errorf("укажите выходной файл через --out")
		}

		// Получаем абсолютные пути, чтобы сравнить расположение файлов.
		outputAbsolute, err := filepath.Abs(outputPath)
		if err != nil {
			return fmt.Errorf("не удалось определить выходной путь: %w", err)
		}

		storageAbsolute, err := filepath.Abs(tasksFile)
		if err != nil {
			return fmt.Errorf("не удалось определить путь хранилища: %w", err)
		}

		// Не позволяем экспорту перезаписать основное хранилище.
		// EqualFold сравнивает строки без учёта регистра букв.
		if strings.EqualFold(outputAbsolute, storageAbsolute) {
			return fmt.Errorf("для экспорта выберите файл, отличный от tasks.json")
		}

		tasks, err := storage.LoadJSON(tasksFile)
		if err != nil {
			return fmt.Errorf("не удалось загрузить задачи: %w", err)
		}

		// Выбираем нашу функцию сохранения по формату.
		switch exportFormat {
		case "json":
			err = storage.SaveJSON(outputPath, tasks)
		case "csv":
			err = storage.SaveCSV(outputPath, tasks)
		}

		if err != nil {
			return fmt.Errorf("не удалось экспортировать задачи: %w", err)
		}

		fmt.Printf("Задачи экспортированы в %s\n", outputPath)
		return nil

	case "load":
		// Импорт заменяет текущий список задач списком из другого файла.
		loadCmd := flag.NewFlagSet("load", flag.ContinueOnError)

		file := loadCmd.String("file", "", "путь к файлу JSON или CSV")

		if err := loadCmd.Parse(args[1:]); err != nil {
			return err
		}

		if loadCmd.NArg() != 0 {
			return fmt.Errorf("лишние аргументы команды load: %v", loadCmd.Args())
		}

		inputPath := strings.TrimSpace(*file)
		if inputPath == "" {
			return fmt.Errorf("укажите импортируемый файл через --file")
		}

		// Ext получает расширение: например, .json или .csv.
		extension := strings.ToLower(filepath.Ext(inputPath))
		if extension != ".json" && extension != ".csv" {
			return fmt.Errorf("импорт поддерживает файлы .json и .csv")
		}

		// Проверяем файл заранее.
		// Это не позволит LoadJSON создать пустой файл при ошибке в пути.
		info, err := os.Stat(inputPath)
		if err != nil {
			return fmt.Errorf("не удалось открыть файл импорта: %w", err)
		}

		// IsDir сообщает, является ли указанный путь папкой.
		if info.IsDir() {
			return fmt.Errorf("для импорта нужен файл, а указан путь к папке")
		}

		// Объявляем список, который заполнит одна из функций загрузки.
		var tasks []todo.Task

		switch extension {
		case ".json":
			tasks, err = storage.LoadJSON(inputPath)
		case ".csv":
			tasks, err = storage.LoadCSV(inputPath)
		}

		// При ошибке загрузки основной файл не перезаписываем.
		if err != nil {
			return fmt.Errorf("не удалось импортировать задачи: %w", err)
		}

		// Успешно загруженный список становится основным списком задач.
		if err := storage.SaveJSON(tasksFile, tasks); err != nil {
			return fmt.Errorf("не удалось сохранить импортированные задачи: %w", err)
		}

		fmt.Printf("Импортировано задач: %d\n", len(tasks))
		return nil

	case "help", "--help", "-h":
		// Несколько вариантов команды приводят к одному действию.
		printUsage()
		return nil

	default:
		// default выполняется, если ни одна команда выше не подошла.
		return fmt.Errorf(
			"неизвестная команда %q; используйте help для справки",
			command,
		)
	}
}

// printUsage — наша функция вывода краткой справки.
// Она ничего не принимает и ничего не возвращает.
func printUsage() {
	fmt.Println("To-Do-менеджер для задач по изучению Go")
	fmt.Println()
	fmt.Println("Команды:")
	fmt.Println(`  add --desc="Изучить срезы Go"`)
	fmt.Println("  list --filter=all")
	fmt.Println("  list --filter=done")
	fmt.Println("  list --filter=pending")
	fmt.Println("  complete --id=1")
	fmt.Println("  delete --id=1")
	fmt.Println("  export --format=csv --out=export.csv")
	fmt.Println("  export --format=json --out=export.json")
	fmt.Println("  load --file=export.csv")
	fmt.Println("  help")
	fmt.Println()
	fmt.Println("Основное хранилище: tasks.json в текущей папке.")
	fmt.Println("Важно: load заменяет текущие задачи импортируемыми.")
}
