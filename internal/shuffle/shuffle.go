package shuffle

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
)

func ShuffleEntries(entries []model.Entry, seed int64) []model.Entry {
	result := make([]model.Entry, len(entries))
	copy(result, entries)

	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

func DailySeed(command string) int64 {
	today := time.Now().Format("2006-01-02")
	hash := sha256.Sum256([]byte(command + today))
	return int64(binary.BigEndian.Uint64(hash[:8]))
}

func RandomSeed() int64 {
	return time.Now().UnixNano()
}

func SaveSeed(cacheDir, command string, seed int64) error {
	dir := filepath.Join(cacheDir, "seeds")
	os.MkdirAll(dir, 0o755)
	return os.WriteFile(filepath.Join(dir, command), []byte(fmt.Sprintf("%d", seed)), 0o644)
}

func LoadSeed(cacheDir, command string) (int64, error) {
	data, err := os.ReadFile(filepath.Join(cacheDir, "seeds", command))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}
