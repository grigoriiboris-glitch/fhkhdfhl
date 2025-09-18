# Справочник по часто используемым функциям в Go

## Основные пакеты и их функции

### fmt - Форматированный ввод/вывод
```go
import "fmt"

fmt.Print("Hello")           // Простой вывод
fmt.Println("Hello")         // Вывод с новой строкой
fmt.Printf("Hello %s", name) // Форматированный вывод

// Ввод
var name string
fmt.Scan(&name)
fmt.Scanln(&name)
fmt.Scanf("%s", &name)
```

### strings - Работа со строками
```go
import "strings"

strings.Contains("hello", "he")     // true - проверка содержания
strings.HasPrefix("hello", "he")    // true - проверка начала
strings.HasSuffix("hello", "lo")    // true - проверка конца
strings.ToUpper("hello")            // "HELLO" - верхний регистр
strings.ToLower("HELLO")            // "hello" - нижний регистр
strings.Trim(" hello ", " ")        // "hello" - обрезка пробелов
strings.Split("a,b,c", ",")         // ["a", "b", "c"] - разделение
strings.Join([]string{"a","b"}, ",")// "a,b" - объединение
strings.Replace("hello", "l", "L", 2) // "heLLo" - замена
```

### strconv - Конвертация строк
```go
import "strconv"

strconv.Atoi("123")          // 123, error - строка в int
strconv.Itoa(123)            // "123" - int в строку
strconv.ParseFloat("3.14", 64) // 3.14, error
strconv.FormatFloat(3.14, 'f', 2, 64) // "3.14"
strconv.ParseBool("true")    // true, error
```

### os - Взаимодействие с ОС
```go
import "os"

os.Getenv("PATH")            // Получить переменную окружения
os.Setenv("KEY", "value")    // Установить переменную окружения
os.Args                      // Аргументы командной строки
os.Exit(1)                   // Выход с кодом ошибки

// Файлы
os.Create("file.txt")        // Создать файл
os.Open("file.txt")          // Открыть файл
os.Remove("file.txt")        // Удалить файл
os.ReadFile("file.txt")      // Прочитать весь файл
os.WriteFile("file.txt", data, 0644) // Записать файл
```

### io/ioutil - Утилиты ввода/вывода
```go
import "io/ioutil"

ioutil.ReadFile("file.txt")  // Прочитать файл (устарело в Go 1.16)
ioutil.WriteFile("file.txt", data, 0644) // Записать файл
ioutil.ReadDir(".")          // Список файлов в директории

// Используйте os вместо ioutil в новых версиях
os.ReadFile("file.txt")
os.WriteFile("file.txt", data, 0644)
```

### path/filepath - Работа с путями
```go
import "path/filepath"

filepath.Join("dir", "file.txt")    // "dir/file.txt"
filepath.Base("/dir/file.txt")      // "file.txt" - имя файла
filepath.Dir("/dir/file.txt")       // "/dir" - директория
filepath.Ext("file.txt")            // ".txt" - расширение
filepath.IsAbs("/path")             // true - абсолютный путь?
filepath.Glob("*.go")               // Поиск файлов по шаблону
```

### time - Работа со временем
```go
import "time"

time.Now()                          // Текущее время
time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) // Конкретная дата

// Форматирование
time.Now().Format("2006-01-02 15:04:05")
time.Parse("2006-01-02", "2023-01-01")

// Таймеры и задержки
time.Sleep(2 * time.Second)
time.AfterFunc(2*time.Second, func() {})
ticker := time.NewTicker(1 * time.Second)

// Измерение времени
start := time.Now()
// выполнение кода
elapsed := time.Since(start)
```

### math - Математические функции
```go
import "math"

math.Abs(-5.5)               // 5.5 - абсолютное значение
math.Ceil(3.14)              // 4 - округление вверх
math.Floor(3.14)             // 3 - округление вниз
math.Round(3.14)             // 3 - округление
math.Max(3, 5)               // 5 - максимум
math.Min(3, 5)               // 3 - минимум
math.Pow(2, 3)               // 8 - степень
math.Sqrt(16)                // 4 - квадратный корень
math.Sin(math.Pi/2)          // 1 - синус
```

### sort - Сортировка
```go
import "sort"

// Сортировка срезов
nums := []int{3, 1, 2}
sort.Ints(nums)              // [1, 2, 3]
sort.Strings([]string{"c", "a", "b"})

// Проверка отсортированности
sort.IntsAreSorted(nums)     // true

// Поиск
sort.SearchInts(nums, 2)     // индекс 2

// Кастомная сортировка
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

### encoding/json - JSON
```go
import "encoding/json"

// Маршалинг (Go → JSON)
data, _ := json.Marshal(user)
json.MarshalIndent(user, "", "  ") // Красивый вывод

// Демаршалинг (JSON → Go)
json.Unmarshal(jsonData, &user)

// Работа с потоками
encoder := json.NewEncoder(file)
encoder.Encode(user)

decoder := json.NewDecoder(file)
decoder.Decode(&user)
```

### net/http - HTTP клиент и сервер
```go
import "net/http"

// Клиент
resp, err := http.Get("http://example.com")
resp, err := http.Post("http://example.com", "application/json", bytes)

// Чтение ответа
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()

// Сервер
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World")
})
http.ListenAndServe(":8080", nil)

// Параметры запроса
r.URL.Query().Get("param")   // GET параметры
r.FormValue("field")         // POST форма данные
```

### sync - Примитивы синхронизации
```go
import "sync"

var mu sync.Mutex
mu.Lock()
// критическая секция
mu.Unlock()

var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // работа
}()
wg.Wait()

// Однократное выполнение
var once sync.Once
once.Do(func() { /* выполнится один раз */ })

// Пулы объектов
var pool = sync.Pool{
    New: func() interface{} { return new(Buffer) },
}
buf := pool.Get().(*Buffer)
pool.Put(buf)
```

### context - Контекст
```go
import "context"

// Создание контекста
ctx := context.Background()
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

ctx, cancel := context.WithCancel(context.Background())
ctx, cancel := context.WithDeadline(context.Background(), time)

// Проверка контекста
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // работа
}

// Передача значений
ctx = context.WithValue(ctx, "key", "value")
value := ctx.Value("key")
```

### errors - Ошибки
```go
import "errors"

// Создание ошибок
err := errors.New("error message")
err := fmt.Errorf("error: %v", details)

// Обертывание ошибок
err = fmt.Errorf("context: %w", originalErr)

// Проверка ошибок
if errors.Is(err, io.EOF) {
    // Ошибка EOF
}
if errors.As(err, &customErr) {
    // Ошибка определенного типа
}
```

### reflect - Рефлексия
```go
import "reflect"

// Получение типа и значения
t := reflect.TypeOf(obj)
v := reflect.ValueOf(obj)

// Работа со структурами
if t.Kind() == reflect.Struct {
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
    }
}

// Вызов методов
method := v.MethodByName("MethodName")
method.Call([]reflect.Value{})
```

### flag - Аргументы командной строки
```go
import "flag"

// Определение флагов
var port = flag.Int("port", 8080, "port number")
var name = flag.String("name", "", "user name")

// Парсинг
flag.Parse()

// Использование
fmt.Printf("Port: %d, Name: %s\n", *port, *name)
```

## Полезные шаблоны

### Дефер с ошибкой
```go
func deferExample() (err error) {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer func() {
        closeErr := file.Close()
        if err == nil {
            err = closeErr
        }
    }()
    
    // работа с файлом
    return nil
}
```

### Кастомный тип ошибки
```go
type CustomError struct {
    Msg string
    Code int
}

func (e *CustomError) Error() string {
    return fmt.Sprintf("%s (code: %d)", e.Msg, e.Code)
}

func someFunction() error {
    return &CustomError{Msg: "something went wrong", Code: 500}
}
```

### Генератор уникальных ID
```go
import "crypto/rand"

func GenerateID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return fmt.Sprintf("%x", b)
}
```

Этот справочник покрывает большинство часто используемых функций в Go. Сохраните его для быстрого доступа!