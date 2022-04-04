package main

import (
	"image"
	"image/color"
)

func drawLines(xBase, yBase int, line []string, m *image.RGBA) {
	for y := 0; y < len(line); y++ {
		for x := 0; x < len(line[y]); x++ {
			if line[y][x] == '*' {
				x1 := float64(xBase + 2*x)
				x2 := x1 + 1.0
				y1 := float64(yBase - 7 + 2*y)
				y2 := y1 + 1.0
				drawlines([4]float64{x1, y1, x2, y2}, m, color.Black)
			}
		}
	}

}

func GWei(xBase, yBase int, m *image.RGBA) {
	gwei := []string{
		"  *****                   ",
		" *     *                  ",
		" *                       *",
		" *                        ",
		" *        *     *  ***   *",
		" *   ***  *     * *   *  *",
		" *     *  *  *  * *****  *",
		" *     *  *  *  * *      *",
		"  *****    ** **   ***   *",
	}
	drawLines(xBase, yBase, gwei, m)
}

func Hours(xBase, yBase int, m *image.RGBA) {
	gwei := []string{
		" *     *                                                         ",
		" *     *                                                         ",
		" *     *                                                         ",
		" *     *   ***   *   *  ***   ***    * *     *  *******   ****  *",
		" * *****  *   *  *   * *  *  *      *  *     *     *     *    *  *",
		" *     *  *   *  *   * *  *   ***   *  *     *     *     *       *",
		" *     *  *   *  *   * ***       *  *  *     *     *     *       *",
		" *     *  *   *  *   * *  *      *  *  *     *     *     *    *  *",
		" *     *   ***    ***  *  *   ***    *  *****      *      ****  *",
	}
	drawLines(xBase, yBase, gwei, m)
}

func drawNumbers(number, xBase, yBase int, m *image.RGBA) {
	// rpainter := Rectangle{length: maxY, width: maxX} // Colored Mask Layer
	// cr := rpainter.drawShape()
	numbers := [][]string{
		{
			"  *** ",
			" *   *",
			" *   *",
			" *   *",
			" *   *",
			" *   *",
			"  *** ",
		},
		{
			"  * ",
			" ** ",
			"  * ",
			"  * ",
			"  * ",
			"  * ",
			" ***",
		},
		{
			"  *** ",
			" *   *",
			"    * ",
			"  *   ",
			" *    ",
			" *    ",
			" *****",
		},
		{
			"  *** ",
			" *   *",
			"     *",
			"  *** ",
			"     *",
			" *   *",
			"  *** ",
		},
		{
			"    ** ",
			"   * * ",
			"  *  * ",
			" ******",
			"     * ",
			"     * ",
			"    ***",
		},
		{
			" *****",
			" *    ",
			" *    ",
			"  *** ",
			"     *",
			" *   *",
			"  *** ",
		},
		{
			"  ****",
			" *    ",
			" *    ",
			" **** ",
			" *   *",
			" *   *",
			"  *** ",
		},
		{
			" *****",
			"     *",
			"     *",
			"    * ",
			"   *  ",
			"  *   ",
			" *    ",
		},
		{
			"  *** ",
			" *   *",
			" *   *",
			"  *** ",
			" *   *",
			" *   *",
			"  *** ",
		},
		{
			"  *** ",
			" *   *",
			" *   *",
			"  ****",
			"     *",
			"     *",
			"  *** ",
		},
	}
	line := []string{"", "", "", "", "", "", ""}
	for number != 0 || len(line[0]) == 0 {
		digit := number % 10
		//fmt.Println(digit, numbers[digit])
		for pos := 0; pos < len(numbers[digit]); pos++ {
			//fmt.Println(pos, numbers[digit][pos])
			line[pos] = numbers[digit][pos] + line[pos]
		}
		number /= 10
	}
	//draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{xBase, yBase}, 2}, image.ZP, draw.Over)

	for y := 0; y < len(line); y++ {
		for x := 0; x < len(line[y]); x++ {
			if line[y][x] == '*' {
				x1 := float64(xBase + 2*x)
				x2 := x1 + 1.0
				y1 := float64(yBase - 7 + 2*y)
				y2 := y1 + 1.0
				drawlines([4]float64{x1, y1, x2, y2}, m, color.Black)
			}
		}
	}
}
