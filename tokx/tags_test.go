package tokx

import (
	"fmt"
	"testing"
)

func TestParseTags(t *testing.T) {
	obj := ParseTags("INT;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL;sdfds:")
	fmt.Printf("%v\n", obj)
}
