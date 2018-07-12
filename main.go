// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2018-06-14 01:17:40 D828E8                       go-experiments/[main.go]
// -----------------------------------------------------------------------------

package main

import (
	"fmt"
	/// "reflect"
)

func main() {
	fmt.Println("running go-experiments")

	/// serverDemo()
	tlsServerDemo()

	///fmt.Println("type of AddFunc:", reflect.TypeOf(AddFunc))
	///fmt.Println("type of DeleteFunc:", reflect.TypeOf(DeleteFunc))
}

type AddFuncT func(a int, b string)

var AddFunc AddFuncT = func(a int, b string) {
}

type DeleteFuncT func(a int, b string)

var DeleteFunc DeleteFuncT = func(a int, b string) {
}

//end
