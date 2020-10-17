package words

import (
	"bytes"
)

// MaxUTF8 utf8 最大值
const MaxUTF8 int = 0xffff

func hashCode(s []rune) int {
	var h = 0
	for _, v := range s {
		h = 31*h + int(toLower(v))
	}
	return h
}

func toLower(c rune) rune {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

type wordList [][]rune

func (w *wordList) contains(s []rune) bool {
	for _, v := range *w {
		if len(v) == len(s) {
			i := 0
			for ; i < len(v) && toLower(v[i]) == toLower(s[i]); i++ {
			}
			if i >= len(v) {
				return true
			}
		}
	}

	return false
}

func (w *wordList) add(s []rune) {
	*w = append(*w, s)
}

// BadWords 屏蔽字
type BadWords struct {
	hashSet        map[int]*wordList // 脏字集合
	fastCharCheck  [MaxUTF8]byte
	fastCharLength [MaxUTF8]byte
	lastCharCheck  *BitArray
	oneCharCheck   *BitArray
	maxLength      int

	temp   []rune
	buffer bytes.Buffer
}

// NewBadWords new
func NewBadWords() *BadWords {
	return &BadWords{
		hashSet:       make(map[int]*wordList),
		lastCharCheck: NewBitArray(MaxUTF8),
		oneCharCheck:  NewBitArray(MaxUTF8),
		maxLength:     0,
	}
}

// AddBadWord 增加屏蔽字
func (b *BadWords) AddBadWord(word string) bool {
	if len(word) == 0 {
		return false
	}

	runeWord := []rune(word)
	for i, v := range runeWord {
		runeWord[i] = toLower(v)
	}

	if int(runeWord[0]) >= MaxUTF8 {
		return false
	}

	if len(runeWord) == 1 {
		b.oneCharCheck.Set(int(runeWord[0]), true)
		b.fastCharCheck[int(runeWord[0])] |= 1
	} else {
		h := hashCode(runeWord)
		if words, ok := b.hashSet[h]; ok {
			if words.contains(runeWord) {
				return false
			}
		} else {
			b.hashSet[h] = &wordList{}
		}
		b.hashSet[h].add(runeWord)

		for i, c := range runeWord {
			if i == len(runeWord)-1 {
				b.lastCharCheck.Set(int(c), true)
			} else if i < 7 {
				if int(c) >= MaxUTF8 {
					continue
				}
				b.fastCharCheck[int(c)] |= 1 << uint32(i)
			} else {
				if int(c) >= MaxUTF8 {
					continue
				}
				b.fastCharCheck[int(c)] |= 0x80
			}
		}
		m := uint32(len(runeWord) - 2)
		if m > 7 {
			m = 7
		}
		b.fastCharLength[runeWord[0]] |= 1 << m
		b.fastCharLength[toLower(runeWord[0])] |= 1 << m
	}

	if len(runeWord) > b.maxLength {
		b.maxLength = len(runeWord)
	}
	return true
}

func (b *BadWords) getTemp() []rune {
	if b.temp == nil {
		b.temp = make([]rune, 0, b.maxLength)
	}
	return b.temp
}

// ReplaceBadWord 替换屏蔽字为*
func (b *BadWords) ReplaceBadWord(text string, replaceChar rune) string {
	var runeText = []rune(text)
	var charCount = len(runeText)
	sub := b.getTemp()
	var find = false
	for index := 0; index < charCount; index++ {
		firstChar := toLower(runeText[index])
		if int(firstChar) >= MaxUTF8 || b.fastCharCheck[int(firstChar)]&1 == 0 {
			continue
		}

		if b.oneCharCheck.Get(int(firstChar)) {
			runeText[index] = replaceChar
			find = true
			continue
		}

		sub = sub[:0]
		sub = append(sub, firstChar)
		hash := int(firstChar)
		spaceCount := 0
		for j := 1; j < (b.maxLength+spaceCount) && j < charCount-index; j++ {
			currentChar := toLower(runeText[index+j])
			if b.isJumpChar(currentChar) {
				spaceCount++
				continue
			}

			m := uint32(j - spaceCount - 1)
			if m > 7 {
				m = 7
			}
			sub = append(sub, currentChar)
			hash = 31*hash + int(currentChar)
			if b.fastCharLength[firstChar]>>m&1 == 1 && b.lastCharCheck.Get(int(currentChar)) {
				if words, ok := b.hashSet[hash]; ok && words.contains(sub) {
					find = true
					for i := index; i <= index+j; i++ {
						if !(b.isJumpChar(runeText[i])) {
							runeText[i] = replaceChar
						}
					}
					index += j
					break
				}
			}

			k := uint32(j - spaceCount)
			if k > 7 {
				k = 7
			}
			if int(currentChar) >= MaxUTF8 || b.fastCharCheck[int(currentChar)]&(1<<k) == 0 {
				break
			}
		}
	}
	if find {
		return string(runeText)
	}
	return text
}

func (b *BadWords) DeleteBadWord(runeText []rune) []rune {
	var charCount = len(runeText)
	sub := b.getTemp()
	var find = false
	b.buffer.Reset()
	for index := 0; index < charCount; index++ {
		firstChar := toLower(runeText[index])
		if int(firstChar) >= MaxUTF8 || b.fastCharCheck[int(firstChar)]&1 == 0 {
			b.buffer.WriteRune(firstChar)
			continue
		}

		if b.oneCharCheck.Get(int(firstChar)) {
			find = true
			continue
		}

		sub = sub[:0]
		sub = append(sub, firstChar)
		hash := int(firstChar)
		spaceCount := 0
		for j := 1; j < (b.maxLength+spaceCount) && j < charCount-index; j++ {
			currentChar := toLower(runeText[index+j])
			if b.isJumpChar(currentChar) {
				spaceCount++
				continue
			}

			m := uint32(j - spaceCount - 1)
			if m > 7 {
				m = 7
			}
			sub = append(sub, currentChar)
			hash = 31*hash + int(currentChar)
			if b.fastCharLength[firstChar]>>m&1 == 1 && b.lastCharCheck.Get(int(currentChar)) {
				if words, ok := b.hashSet[hash]; ok && words.contains(sub) {
					find = true
					index += j
					break
				}
			}

			k := uint32(j - spaceCount)
			if k > 7 {
				k = 7
			}
			if int(currentChar) >= MaxUTF8 || b.fastCharCheck[int(currentChar)]&(1<<k) == 0 {
				b.buffer.WriteRune(currentChar)
				break
			}
		}
	}
	if find {
		return []rune(b.buffer.String())
	}
	return runeText
}

// ContainsBadWord 是否含有屏蔽字
func (b *BadWords) ContainsBadWord(text string) bool {
	var runeText = []rune(text)
	var charCount = len(runeText)
	sub := b.getTemp()
	for index := 0; index < charCount; index++ {
		firstChar := toLower(runeText[index])
		if int(firstChar) >= MaxUTF8 || b.fastCharCheck[int(firstChar)]&1 == 0 {
			continue
		}

		if b.oneCharCheck.Get(int(firstChar)) {
			return true
		}

		sub = sub[:0]
		sub = append(sub, firstChar)
		hash := int(firstChar)
		spaceCount := 0
		for j := 1; j < b.maxLength+spaceCount && j < charCount-index; j++ {
			currentChar := toLower(runeText[index+j])
			if b.isJumpChar(currentChar) {
				spaceCount++
				continue
			}

			m := uint32(j - spaceCount - 1)
			if m > 7 {
				m = 7
			}
			sub = append(sub, currentChar)
			hash = 31*hash + int(currentChar)
			if b.fastCharLength[firstChar]>>m&1 == 1 && b.lastCharCheck.Get(int(currentChar)) {
				if words, ok := b.hashSet[hash]; ok && words.contains(sub) {
					return true
				}
			}

			k := uint32(j - spaceCount)
			if k > 7 {
				k = 7
			}
			if int(currentChar) >= MaxUTF8 || b.fastCharCheck[int(currentChar)]&(1<<k) == 0 {
				break
			}
		}
	}
	return false
}

func (b *BadWords) isJumpChar(c rune) bool {
	return c == ' ' || c == '\t'
}
