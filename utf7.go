package xxcrypt

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	base64chars      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	base64charsArray = strings.Split(base64chars, "")
)

func tocharcodes(code string, params []string) string {
	var (
		splitChar string
		joinChar  string
		codeArray []string
	)
	if params[0] == "''" {
		splitChar = ""
	} else {
		splitChar = params[0]
	}
	if params[1] == "','" {
		joinChar = ","
	} else {
		joinChar = params[1]
	}
	codeArray = strings.Split(code, splitChar)
	for i, _ := range codeArray {
		codeArray[i] = strconv.Itoa(int(codeArray[i][0]))
	}
	codeArray = append(codeArray, joinChar)
	code = ""
	for _, c := range codeArray {
		code = code + c
	}
	return code
}
func zerofill(code string, param int) string {
	if param > 10000 || param < 0 {
		param = 0
	}
	for {
		if param > len(code) {
			code = "0" + code
		} else {
			break
		}
	}
	return code
}
func UTF7_Encode(istr string) string {
	var code string
	var sixBits []string
	for _, c := range strings.Split(istr, "") {
		strchar := fmt.Sprintf("%b", int(c[0]))
		code = code + zerofill(strchar, 16)
	}
	for i := 0; i < len(code); i += 6 {
		var lastIndex int
		if i+6 > len(code) {
			lastIndex = len(code)
		} else {
			lastIndex = i + 6
		}
		sixBits = append(sixBits, code[i:lastIndex])
	}
	if len(sixBits[len(sixBits)-1]) < 6 {
		tmp_codes3 := strings.Repeat("0", 6-len(sixBits[len(sixBits)-1]))
		sixBits[len(sixBits)-1] = sixBits[len(sixBits)-1] + tmp_codes3
	}
	for i, _ := range sixBits {
		s, _ := strconv.ParseInt(sixBits[i], 2, 10)
		sixBits[i] = base64charsArray[s]
	}
	return "+" + strings.Join(sixBits, "") + "-"
}
func UTF7_Decode(istr string) string {
	fr, _ := regexp.Compile(`\+[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+\/]+-`)
	for _, findStr := range fr.FindAllString(istr, -1) {
		findStr = findStr[1 : len(findStr)-1]
		var decoded string
		for _, ch := range strings.Split(findStr, "") {
			tmp_codes1 := fmt.Sprintf("%b", strings.Index(base64chars, ch))
			decoded += zerofill(tmp_codes1, 6)
		}
		var sixteenBits []string
		for i := 0; i < len(decoded); i += 16 {
			var lastIndex int
			if i+16 > len(decoded) {
				lastIndex = len(decoded)
			} else {
				lastIndex = i + 16
			}
			sixteenBits = append(sixteenBits, decoded[i:lastIndex])
		}
		decoded = ""
		for _, bt := range sixteenBits {
			if len(bt) < 16 {
				if len(bt) > 4 || bt[0:0+1] != "0" {
					log.Fatalln("Invalid UTF-7")
				}
			} else {
				c, _ := strconv.ParseInt(bt, 2, 10)
				decoded += string(c)
			}
		}
		istr = strings.Replace(istr, "+"+findStr+"-", decoded, -1)
	}
	fmt.Println(istr)
	return istr
}
