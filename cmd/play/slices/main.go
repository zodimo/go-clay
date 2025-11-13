package main

import "fmt"

func main() {

	a := make([]int, 10)
	a[0] = 1
	a[1] = 2
	a[2] = 3
	a[3] = 4
	a[4] = 5
	a[5] = 6
	a[6] = 7
	a[7] = 8
	a[8] = 9
	a[9] = 10

	b := a[0:5]
	fmt.Println(b)
	b = a[0:6]
	fmt.Println(b)
	c := a[5:10]

	b = c

	d := make([]int, 10)
	d[0] = 1
	d[1] = 2
	d[2] = 3
	d[3] = 4
	d[4] = 5
	d[5] = 6
	d[6] = 7
	d[7] = 8
	d[8] = 9
	d[9] = 10

	fmt.Println("d length", len(d))
	d = b

	fmt.Println("d length", len(d))

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(b)
	fmt.Println(d)
}
