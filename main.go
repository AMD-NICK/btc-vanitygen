package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func main() {
	minSequentRepeats := flag.Int("sequentrepeats", 9, "Минимальная длина повторов символов подряд")
	minRepeats := flag.Int("repeats", 14, "Минимальная длина повтора символа в любом месте")
	maxUnique := flag.Int("unique", 10, "Максимальное количество уникальных символов")
	numWorkers := flag.Int("workers", 16, "Количество потоков")

	flag.Parse()

	attempts := 0
	var wg sync.WaitGroup
	resultChan := make(chan string, *numWorkers)

	fmt.Printf("Запуск с настройками: sequentrepeats=%d, repeats=%d, unique=%d, workers=%d\n", *minSequentRepeats, *minRepeats, *maxUnique, *numWorkers)

	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(*minSequentRepeats, *minRepeats, *maxUnique, resultChan, &wg, &attempts)
	}

	// Закрываем канал после завершения всех воркеров
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Выводим результаты по мере их поступления
	for result := range resultChan {
		fmt.Println(result)
	}
}

func worker(minSequentRepeats, minRepeats, maxUnique int, resultChan chan string, wg *sync.WaitGroup, attempts *int) {
	defer wg.Done()

	for {
		*attempts++
		if *attempts%100000000 == 0 {
			fmt.Printf("Попыток (млн): %d...\n", *attempts/1000000)
		}

		privKey, err := btcec.NewPrivateKey()
		if err != nil {
			fmt.Println("Ошибка генерации ключа:", err)
			return
		}

		pubKey := privKey.PubKey()
		address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
		if err != nil {
			fmt.Println("Ошибка создания адреса:", err)
			return
		}

		addrStr := address.EncodeAddress()
		pretty := isPrettyAddress(addrStr, minSequentRepeats, minRepeats, maxUnique)
		if pretty != "" {
			result := fmt.Sprintf("%s Адрес: %s, ключ: %s", pretty, addrStr, privKeyToWIF(privKey))
			resultChan <- result
		}
	}
}

func isPrettyAddress(addr string, minSequentRepeats, minRepeats, maxUnique int) string {
	if hasRepeatingChars(addr, minSequentRepeats) {
		return "🔁"
	}
	if hasFrequentChar(addr, minRepeats) {
		return "💠"
	}
	if countUniqueChars(addr) <= maxUnique {
		return "🌈"
	}
	return ""
}

func hasRepeatingChars(s string, minSequentRepeats int) bool {
	s = strings.ToLower(s)
	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= minSequentRepeats {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

func hasFrequentChar(s string, minRepeats int) bool {
	s = strings.ToLower(s)
	charCount := make(map[rune]int)
	for _, char := range s {
		charCount[char]++
		if charCount[char] >= minRepeats {
			return true
		}
	}
	return false
}

func countUniqueChars(s string) int {
	unique := make(map[rune]bool)
	for _, char := range strings.ToLower(s) {
		unique[char] = true
	}
	return len(unique)
}

func privKeyToWIF(privKey *btcec.PrivateKey) string {
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
	if err != nil {
		panic(err)
	}
	return wif.String()
}
