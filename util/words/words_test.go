package words

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

var badWorlds = make([]string, 0)

func init() {
	fmt.Println("读取badwords.txt")
	file, err := os.Open(`../../data/badwords.txt`)
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		} else {
			badWorlds = append(badWorlds, line)
		}
	}
}

func TestAddBadWord(t *testing.T) {
	bw := NewBadWords()
	if !bw.AddBadWord("A") {
		t.Error("AddBadWord")
	}
	if !bw.AddBadWord("sb") {
		t.Error("AddBadWord")
	}
	if bw.AddBadWord("sb") {
		t.Error("AddBadWord")
	}
	if !bw.AddBadWord("FaLunDaFaHao") {
		t.Error("AddBadWord")
	}
}

func TestReplace(t *testing.T) {
	bw := NewBadWords()
	words := []string{"Sb", "c", "FaLunDaFaHao"}
	for _, word := range words {
		bw.AddBadWord(word)
	}

	var testData = []struct {
		str    string
		result string
	}{
		{"C", "*"},
		{"S b", "* *"},
		{"bb", "bb"},
		{"FaLunDaFaHa  ddd", "FaLunDaFaHa  ddd"},
		{"FaLunDaFaHaoFaLunDaFaHao M", "************************ M"},
	}

	for _, v := range testData {
		if bw.ReplaceBadWord(v.str, '*') != v.result {
			t.Error("ReplaceBadWord")
		}
	}
}

func TestContainsBadWord(t *testing.T) {
	bw := NewBadWords()
	words := []string{"sb", "C", "sb", "FaLunDaFaHao"}
	for _, word := range words {
		bw.AddBadWord(word)
	}

	var testData = []struct {
		str    string
		result bool
	}{
		{"C", true},
		{"S b", true},
		{"bb", false},
		{"FaLunDaFaHa  ddd", false},
		{"FaLunDaFaHaoFaLunDaFaHao M", true},
	}

	for _, v := range testData {
		if bw.ContainsBadWord(v.str) != v.result {
			t.Error("ContainsBadWord", v)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	//bw := NewBadWords()

	for i := 0; i < b.N; i++ {
		//bw.AddBadWord(table.Words.Data[i%table.Words.Count()].Word)

		BadWords := NewBadWords()
		for _, v := range badWorlds {
			BadWords.AddBadWord(v)
		}
	}
}

func BenchmarkReplace(b *testing.B) {
	bw := NewBadWords()

	for _, v := range badWorlds {
		bw.AddBadWord(v)
	}
	str := "色动日师地胺摄☺a☺a☻☺b☺色动日师地胺摄毛泽东色动日师十七大摄☺a☺a☻☺b☺色动日师地胺摄☺a☺a☻☺b☺色动日师地胺摄☺a☺a☻☺"
	c := '*'
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bw.ReplaceBadWord(str, c)
	}
}

func BenchmarkContains(b *testing.B) {
	bw := NewBadWords()

	for _, v := range badWorlds {
		bw.AddBadWord(v)
	}

	str := "色动日师地胺摄"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bw.ContainsBadWord(str)
	}
}

func TestChatManager_Emoj(t *testing.T) {
	bw := NewBadWords()
	bw.AddBadWord("<anger>")

	s := []rune("a<<anger>>b<ang er>>b<anger>a")
	p := bw.DeleteBadWord(s)
	if strings.Compare(string(p), "a<>b>ba") != 0 {
		t.Error("DeleteBadWord")
	}
}

func BenchmarkChatManager_Emoj(b *testing.B) {
	bw := NewBadWords()
	for _, v := range badWorlds {
		bw.AddBadWord(v)
	}

	s := []rune("a<<anger>>b<ang er>>b<anger>a")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bw.DeleteBadWord(s)
	}
}
