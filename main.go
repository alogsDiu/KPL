package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var codeRunes []rune

var variables = make(map[string]*WrapperForKPLTypes)

type WrapperForKPLTypes struct {
	value        any
	kplTypeIndex int
}

// coreKPLTypes = ["string","float","int","any","function","object"] // "null"
var coreKPLTypes = []string{"мәтін", "сан", "бүтін_сан", "белгісіз", "қасиет", "зат"} // "ештеңе"
var KPLFunctions = []string{}
var row = 0

func main() {

	the_file_name := "test.kpl"

	the_file, err := os.ReadFile(the_file_name)
	if err != nil {
		fmt.Println("'" + the_file_name + "' файлы жоқ!")
	}
	entire_code := strings.Trim(string(the_file), " ")

	startCompiler(&entire_code)
}

func startCompiler(codePointer *string) {
	code := *codePointer

	codeRunes = []rune(code)

	comment_started := false

	firstItem := ""

	for i := 0; i < int(len(codeRunes)); i++ {
		if codeRunes[i] == '\n' {
			comment_started = false
			row++
		} else if comment_started {
			continue
		} else if codeRunes[i] == '/' {
			comment_started = true
		} else if codeRunes[i] == ' ' || codeRunes[i] == '\n' {
			if len(firstItem) != 0 {
				i = choseWhatWayToTake(i, firstItem)
				firstItem = ""
			}
			if codeRunes[i] == '\n' {
				row++
			}
		} else {
			firstItem += string(codeRunes[i])
		}
	}

	for k, v := range variables {
		fmt.Println(k, v)
	}
}

func choseWhatWayToTake(currentIndex int, theFirstItem string) int {
	//Check if variable with that name exists
	for k, v := range variables {
		if k == theFirstItem {
			return changeVariable(currentIndex, v)
		}
	}

	//Check if it is not the function name
	for _, fName := range KPLFunctions {
		if fName == theFirstItem {
			fmt.Println("Сіз '"+fName+"'ны айнымалы атау ретінде пайдалана алмайсыз! қатар:", row)
			os.Exit(0)
		}
	}

	//Check if it is not the function name
	for _, typeName := range coreKPLTypes {
		if typeName == theFirstItem {
			fmt.Println("Сіз '"+typeName+"'ны айнымалы атау ретінде пайдалана алмайсыз! қатар:", row)
			os.Exit(0)
		}
	}

	theSecondStr := ""
	for i := currentIndex; i < len(codeRunes); i++ {
		if codeRunes[i] == ' ' || codeRunes[i] == '\n' || codeRunes[i] == '=' {
			if codeRunes[i] == '\n' {
				row++
			}
			if len(theSecondStr) != 0 {
				futureValue := &WrapperForKPLTypes{}
				variables[theFirstItem] = futureValue
				return createVariable(i, futureValue, theSecondStr)
			}
		} else {
			theSecondStr += string(codeRunes[i])
		}
	}
	return 0
}

func changeVariable(currentIndex int, variablesObject *WrapperForKPLTypes) int {
	index := 0
	variablesObject.value, index = getTheValueAndReturnThefarthestParsedIndex(variablesObject.kplTypeIndex, currentIndex)
	return index
}

func createVariable(currentIndex int, variablesObject *WrapperForKPLTypes, typeName string) int {
	for index, tn := range coreKPLTypes {
		if typeName == tn {
			variablesObject.kplTypeIndex = index
			variablesObject.value, index = getTheValueAndReturnThefarthestParsedIndex(index, currentIndex)
			return index
		}
	}
	fmt.Println("Мұндай айнымалы түрі жоқ қатар:", row)
	os.Exit(0)
	return 0
}

func getTheValueAndReturnThefarthestParsedIndex(typeIndex int, index int) (any, int) {

	switch typeIndex {
	case 0:
		return stringValueParser(index)
	case 1:
		return float64ValueParser(index)
	case 2:
		return intValueParser(index)
	case 3:
		return anyValueParser(index)
	case 4:
		return functionValueParser(index)
		// case 5:
		// 	objectValueParser(index)
	}

	return nil, 0
}

func stringValueParser(index int) (string, int) {
	theValue := ""
	stringStart := false

	for i := index; i < len(codeRunes); i++ {
		if codeRunes[i] == '=' {
			continue
		}

		if stringStart {
			if codeRunes[i] == '"' {
				return theValue, i
			} else {
				theValue += string(codeRunes[i])
			}
		} else if codeRunes[i] == ' ' || codeRunes[i] == '\n' {
			if codeRunes[i] == '\n' {
				row++
			}
			if len(theValue) != 0 {
				return getExistingValue(theValue).(string), i
			}
		} else if codeRunes[i] == '"' {
			stringStart = true
		} else {
			theValue += string(codeRunes[i])
		}
	}
	return "", 0
}

func getExistingValue(varName string) any {

	if variables[varName] == nil {
		fmt.Println("'", varName, " 'деген айнымалы жоқ қатар:", row-1)
		os.Exit(0)
	} else {
		return variables[varName].value
	}
	return ""
}

func float64ValueParser(index int) (float64, int) {
	theValue := ""
	theDigitStart := false
	theVarNameStart := false
	for i := index; i < len(codeRunes); i++ {
		if codeRunes[i] == '=' {
			continue
		}
		if !theDigitStart && !theVarNameStart {
			if codeRunes[i] == ' ' {
				continue
			}
			if codeRunes[i] >= '0' && codeRunes[i] <= '9' {
				theDigitStart = true
			} else {
				theVarNameStart = true
			}
		}
		if codeRunes[i] == ' ' || codeRunes[i] == '\n' {
			if codeRunes[i] == '\n' {
				row++
			}
			if len(theValue) != 0 && theVarNameStart {
				return getExistingValue(theValue).(float64), i
			} else if len(theValue) != 0 {
				f, err := strconv.ParseFloat(theValue, 64)
				if err != nil {
					fmt.Println(theValue, "- сан емес ! қатар:", row)
					os.Exit(0)
				}
				return f, i
			}
		} else {
			theValue += string(codeRunes[i])
		}
	}
	return float64(0), 0
}

func intValueParser(index int) (int, int) {
	theValue := ""
	theDigitStart := false
	theVarNameStart := false
	for i := index; i < len(codeRunes); i++ {
		if codeRunes[i] == '=' {
			continue
		}
		if !theDigitStart && !theVarNameStart {
			if codeRunes[i] == ' ' {
				continue
			}
			if codeRunes[i] >= '0' && codeRunes[i] <= '9' {
				theDigitStart = true
			} else {
				theVarNameStart = true
			}
		}
		if codeRunes[i] == ' ' || codeRunes[i] == '\n' {
			if codeRunes[i] == '\n' {
				row++
			}
			if len(theValue) != 0 && theVarNameStart {
				return getExistingValue(theValue).(int), i
			} else if len(theValue) != 0 {
				d, err := strconv.ParseInt(theValue, 0, 64)
				if err != nil {
					fmt.Println(theValue, "- бүтін_сан емес ! қатар:", row)
					os.Exit(0)
				}
				return int(d), i
			}
		} else {
			theValue += string(codeRunes[i])
		}
	}
	return int(0), 0
}

func anyValueParser(index int) (any, int) {
	index++
	return "", 0
}

func functionValueParser(index int) (any, int) {
	index++
	return func() {
		fmt.Println("shit")
	}, 0
}

// func objectValueParser(index int) (object, int) {

// }
