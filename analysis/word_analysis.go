package analysis

import (
	"bufio"
	"fmt"
	"io"
	"lin/util"
	"lin/util/words"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var badWords = words.NewBadWords()

const step = int64(5)

type worldFrequency struct {
	frequency map[int64]int //时间 次数 todo 可以考虑使用优先队列或者最小堆
}

func (w *worldFrequency) addTimes() {
	t := time.Now().Unix()
	w.check(t)
	w.frequency[t] += 1
}

func (w *worldFrequency) check(t int64) {
	for o := range w.frequency { //最多循环step次数,这里5秒也就是5次
		if o < (t - step) {
			delete(w.frequency, o)
		}
	}
}

var worldsFrequency = make(map[string]*worldFrequency)

func init() {
	p, _ := util.GetCurrentPath()
	p = path.Join(p, "../data/badwords.txt")
	fmt.Println("读取:" + p)
	file, err := os.Open(p) //for normal
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		} else {
			line = strings.Trim(line, "\n")
			line = strings.Trim(line, "\r")
			badWords.AddBadWord(line)
		}
	}
}

func ReplaceBadWord(text string, replaceChar rune) string {
	return badWords.ReplaceBadWord(text, replaceChar)
}

func AddWorldFrequency(text string) {
	var words = strings.Split(text, " ")
	for _, word := range words {
		if word == "" {
			continue
		}
		if _, ok := worldsFrequency[word]; !ok {
			worldsFrequency[word] = &worldFrequency{frequency: make(map[int64]int)}
		}
		worldsFrequency[word].addTimes()
	}
}
func GetWorldsFrequencyStatistics() []string {
	t := time.Now().Unix()
	out := make([]string, 0)
	for k, world := range worldsFrequency {
		world.check(t)
		times := 0
		for _, o := range world.frequency {
			times += o
		}
		if times == 0 {
			delete(worldsFrequency, k)
			continue
		}
		out = append(out, k+":"+strconv.Itoa(times))
	}
	if len(out) == 0 {
		out = append(out, fmt.Sprintf("最近%ds无流行词汇", step))
	}
	return out
}
