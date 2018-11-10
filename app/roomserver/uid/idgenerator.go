package uid

import "strconv"

var nextid int
const prefix = "uid"

func init(){
	nextid = 0
}

func GeneratorID() string{
	str := strconv.Itoa(nextid)
	switch 6 - len(str) {
	case 1:str = "0"+str
	case 2:str = "00"+str
	case 3:str = "000"+str
	case 4:str = "0000"+str
	case 5:str = "00000"+str
	}
	nextid++
	return str
}
