package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"bufio"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Stats struct {
	attempts  uint64
	startTime time.Time
}

func (s *Stats) incrementAttempts() {
	atomic.AddUint64(&s.attempts, 1)
}

func (s *Stats) getStats() (uint64, float64) {
	attempts := atomic.LoadUint64(&s.attempts)
	duration := time.Since(s.startTime).Seconds()
	return attempts, float64(attempts) / duration
}

const helpText = `
Ethereum Address Generator with Pattern Matching

Usage: genaddr [options]

Options:
  -pattern string
        Pattern to match (required). Examples:
        - "123*"    : address starts with "123"
        - "*123"    : address ends with "123"
        - "*123*"   : address contains "123"
        - "123*321" : address starts with "123" and ends with "321"
        - "*123*456*": address contains "123" followed by "456"
        Special characters:
        - "#" : any digit (0-9)
        - "@" : any letter (a-f)
        Examples with special characters:
        - "###*"   : starts with any 3 digits
        - "@@@*"   : starts with any 3 letters
        - "#@#@*"  : alternating digits and letters
        Multiple patterns can be specified using commas:
        - "123*,*456,*789*" : matches any of the patterns
  -workers int
        Number of worker goroutines (default 4)
  -continue
        Continue searching after finding a match
  -output string
        Save found addresses to file
  -help
        Show this help message

Examples:
  genaddr -pattern "123*"
  genaddr -pattern "dead*beef"
  genaddr -pattern "*cafe*,*babe*" -workers 8 -continue -output results.txt
  genaddr -pattern "1*2*3*4" -workers 4

GitHub Repository:
  https://github.com/grom42kem/genaddr

Support the Project:
  If you find this tool useful, you can support its development by sending donations to:
  0x77777777b487e2FD60F3C60B080E03e7247338f6
`

func main() {
	// Определяем флаги
	patterns := flag.String("pattern", "", "")
	numWorkers := flag.Int("workers", 4, "")
	continueSearch := flag.Bool("continue", false, "")
	outputFile := flag.String("output", "", "")
	help := flag.Bool("help", false, "")

	// Переопределяем Usage для вывода собственной справки
	flag.Usage = func() {
		fmt.Print(helpText)
	}

	flag.Parse()

	// Показываем справку
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *patterns == "" {
		fmt.Println("Error: -pattern flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Разделяем паттерны по запятой
	patternList := strings.Split(*patterns, ",")
	// Убираем пробелы
	for i := range patternList {
		patternList[i] = strings.TrimSpace(patternList[i])
	}

	// Открываем файл для записи, если указан
	var writer *bufio.Writer
	var file *os.File
	if *outputFile != "" {
		var err error
		file, err = os.OpenFile(*outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	}

	// Инициализируем статистику
	stats := &Stats{
		startTime: time.Now(),
	}

	// Создаем каналы для параллельной обработки
	found := make(chan bool)
	addresses := make(chan string)
	var wg sync.WaitGroup

	// Запускаем горутину для отображения статистики
	stopStats := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				attempts, speed := stats.getStats()
				fmt.Printf("\rAddresses checked: %d | Speed: %.2f addr/sec", attempts, speed)
			case <-stopStats:
				return
			}
		}
	}()

	// Запускаем горутину для обработки найденных адресов
	go func() {
		for addr := range addresses {
			fmt.Printf("\n\nFound matching address!\n%s\n", addr)
			if writer != nil {
				_, err := writer.WriteString(addr + "\n")
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
				}
				writer.Flush()
			}
			if !*continueSearch {
				found <- true
				return
			}
		}
	}()

	// Запускаем горутины для генерации адресов
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			generateAddresses(patternList, addresses, stats)
		}()
	}

	// Ждем нахождения подходящего адреса или завершения по Ctrl+C
	<-found
	close(stopStats)
	close(addresses)

	// Выводим финальную статистику
	attempts, speed := stats.getStats()
	fmt.Printf("\n\nFinal Statistics:\n")
	fmt.Printf("Total addresses checked: %d\n", attempts)
	fmt.Printf("Average speed: %.2f addr/sec\n", speed)
	fmt.Printf("Search time: %.2f sec\n", time.Since(stats.startTime).Seconds())
}

func generateAddresses(patterns []string, addresses chan<- string, stats *Stats) {
	for {
		// Генерируем новую пару ключей
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			continue
		}

		stats.incrementAttempts()

		// Получаем публичный ключ
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			continue
		}

		// Получаем адрес
		address := crypto.PubkeyToAddress(*publicKeyECDSA)
		addressHex := address.Hex()

		// Проверяем соответствие хотя бы одному паттерну
		for _, pattern := range patterns {
			if checkAddress(addressHex[2:], pattern) { // Пропускаем "0x" префикс при проверке
				privateKeyBytes := crypto.FromECDSA(privateKey)
				privateKeyHex := hexutil.Encode(privateKeyBytes)
				result := fmt.Sprintf("Address: %s\nPrivate Key: %s", addressHex, privateKeyHex)
				addresses <- result
				break
			}
		}
	}
}

func checkAddress(addr, pattern string) bool {
	// Приводим к нижнему регистру для регистронезависимого сравнения
	addr = strings.ToLower(addr)
	pattern = strings.ToLower(pattern)

	// Разбиваем паттерн на части по звездочке
	parts := strings.Split(pattern, "*")
	
	// Если паттерн не содержит звездочек, проверяем точное совпадение
	if len(parts) == 1 {
		return matchExactPattern(addr, pattern)
	}

	// Текущая позиция в адресе
	pos := 0

	// Проверяем каждую часть паттерна
	for i, part := range parts {
		if part == "" {
			continue
		}

		// Для первой части проверяем префикс
		if i == 0 {
			if !matchPrefix(addr, part) {
				return false
			}
			pos = len(part)
			continue
		}

		// Для последней части проверяем суффикс
		if i == len(parts)-1 {
			return matchSuffix(addr[pos:], part)
		}

		// Для промежуточных частей ищем подстроку с учетом специальных символов
		found := false
		for j := pos; j <= len(addr)-len(part); j++ {
			if matchExactPattern(addr[j:j+len(part)], part) {
				pos = j + len(part)
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// Функция для проверки точного совпадения с учетом специальных символов
func matchExactPattern(s, pattern string) bool {
	if len(s) != len(pattern) {
		return false
	}
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '#':
			if !isDigit(s[i]) {
				return false
			}
		case '@':
			if !isHexLetter(s[i]) {
				return false
			}
		default:
			if s[i] != pattern[i] {
				return false
			}
		}
	}
	return true
}

// Функция для проверки префикса с учетом специальных символов
func matchPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return matchExactPattern(s[:len(prefix)], prefix)
}

// Функция для проверки суффикса с учетом специальных символов
func matchSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return matchExactPattern(s[len(s)-len(suffix):], suffix)
}

// Проверка, является ли символ шестнадцатеричной цифрой
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Проверка, является ли символ шестнадцатеричной буквой
func isHexLetter(c byte) bool {
	return (c >= 'a' && c <= 'f')
} 