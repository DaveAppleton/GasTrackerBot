package main

import (
	"image/png"
	"time"

	"os"

	"testing"

	"github.com/shopspring/decimal"
)

func TestPng(t *testing.T) {
	data := []GasData{
		{BlockNum: 10252175, SafeLow: decimal.NewFromFloat(28)},
		{BlockNum: 10252435, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10252716, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10253000, SafeLow: decimal.NewFromFloat(18.6)},
		{BlockNum: 10253256, SafeLow: decimal.NewFromFloat(18)},
		{BlockNum: 10253554, SafeLow: decimal.NewFromFloat(17.2)},
		{BlockNum: 10253813, SafeLow: decimal.NewFromFloat(17.7)},
		{BlockNum: 10254080, SafeLow: decimal.NewFromFloat(17.7)},
		{BlockNum: 10254347, SafeLow: decimal.NewFromFloat(17.9)},
		{BlockNum: 10254626, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10254892, SafeLow: decimal.NewFromFloat(17.3)},
		{BlockNum: 10255164, SafeLow: decimal.NewFromFloat(19.1)},
		{BlockNum: 10255422, SafeLow: decimal.NewFromFloat(23)},
		{BlockNum: 10255682, SafeLow: decimal.NewFromFloat(18.4)},
		{BlockNum: 10255965, SafeLow: decimal.NewFromFloat(24.1)},
		{BlockNum: 10256263, SafeLow: decimal.NewFromFloat(19)},
		{BlockNum: 10256525, SafeLow: decimal.NewFromFloat(18.1)},
		{BlockNum: 10256799, SafeLow: decimal.NewFromFloat(23.2)},
		{BlockNum: 10257067, SafeLow: decimal.NewFromFloat(21.6)},
		{BlockNum: 10257335, SafeLow: decimal.NewFromFloat(18.9)},
		{BlockNum: 10257581, SafeLow: decimal.NewFromFloat(24)},
		{BlockNum: 10257874, SafeLow: decimal.NewFromFloat(23)},
		{BlockNum: 10258148, SafeLow: decimal.NewFromFloat(21)},
		{BlockNum: 10258405, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10258688, SafeLow: decimal.NewFromFloat(15)},
		{BlockNum: 10258977, SafeLow: decimal.NewFromFloat(14.3)},
		{BlockNum: 10259254, SafeLow: decimal.NewFromFloat(13)},
		{BlockNum: 10259546, SafeLow: decimal.NewFromFloat(13)},
		{BlockNum: 10259803, SafeLow: decimal.NewFromFloat(13)},
		{BlockNum: 10260092, SafeLow: decimal.NewFromFloat(12)},
		{BlockNum: 10260344, SafeLow: decimal.NewFromFloat(13)},
		{BlockNum: 10260617, SafeLow: decimal.NewFromFloat(12.5)},
		{BlockNum: 10260904, SafeLow: decimal.NewFromFloat(13.2)},
		{BlockNum: 10261196, SafeLow: decimal.NewFromFloat(13.3)},
		{BlockNum: 10261452, SafeLow: decimal.NewFromFloat(15.8)},
		{BlockNum: 10261718, SafeLow: decimal.NewFromFloat(18.2)},
		{BlockNum: 10262006, SafeLow: decimal.NewFromFloat(12)},
		{BlockNum: 10262298, SafeLow: decimal.NewFromFloat(12.1)},
		{BlockNum: 10262579, SafeLow: decimal.NewFromFloat(14)},
		{BlockNum: 10262831, SafeLow: decimal.NewFromFloat(15.5)},
		{BlockNum: 10263104, SafeLow: decimal.NewFromFloat(16)},
		{BlockNum: 10263390, SafeLow: decimal.NewFromFloat(12.2)},
		{BlockNum: 10263644, SafeLow: decimal.NewFromFloat(14)},
		{BlockNum: 10263911, SafeLow: decimal.NewFromFloat(14)},
		{BlockNum: 10264189, SafeLow: decimal.NewFromFloat(13.9)},
		{BlockNum: 10264458, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10264735, SafeLow: decimal.NewFromFloat(18.3)},
		{BlockNum: 10265031, SafeLow: decimal.NewFromFloat(18.1)},
		{BlockNum: 10265293, SafeLow: decimal.NewFromFloat(18)},
		{BlockNum: 10265544, SafeLow: decimal.NewFromFloat(19)},
		{BlockNum: 10265819, SafeLow: decimal.NewFromFloat(17)},
		{BlockNum: 10266082, SafeLow: decimal.NewFromFloat(10.8)},
		{BlockNum: 10266361, SafeLow: decimal.NewFromFloat(11)},
		{BlockNum: 10266646, SafeLow: decimal.NewFromFloat(12.1)},
		{BlockNum: 10266933, SafeLow: decimal.NewFromFloat(11)},
		{BlockNum: 10267180, SafeLow: decimal.NewFromFloat(12)},
		{BlockNum: 10267454, SafeLow: decimal.NewFromFloat(13)},
		{BlockNum: 10267714, SafeLow: decimal.NewFromFloat(16.1)},
		{BlockNum: 10267970, SafeLow: decimal.NewFromFloat(23.1)},
		{BlockNum: 10268244, SafeLow: decimal.NewFromFloat(23)},
		{BlockNum: 10268532, SafeLow: decimal.NewFromFloat(26)},
		{BlockNum: 10268789, SafeLow: decimal.NewFromFloat(50)},
		{BlockNum: 10269036, SafeLow: decimal.NewFromFloat(50)},
		{BlockNum: 10269294, SafeLow: decimal.NewFromFloat(35)},
		{BlockNum: 10269575, SafeLow: decimal.NewFromFloat(40)},
		{BlockNum: 10269852, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10270142, SafeLow: decimal.NewFromFloat(25)},
		{BlockNum: 10270407, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10270677, SafeLow: decimal.NewFromFloat(34)},
		{BlockNum: 10270952, SafeLow: decimal.NewFromFloat(43)},
		{BlockNum: 10271238, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10271530, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10271788, SafeLow: decimal.NewFromFloat(23)},
		{BlockNum: 10272084, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10272365, SafeLow: decimal.NewFromFloat(26)},
		{BlockNum: 10272644, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10272919, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10273159, SafeLow: decimal.NewFromFloat(22)},
		{BlockNum: 10273446, SafeLow: decimal.NewFromFloat(23.3)},
		{BlockNum: 10273718, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10274008, SafeLow: decimal.NewFromFloat(20.4)},
		{BlockNum: 10274294, SafeLow: decimal.NewFromFloat(22)},
		{BlockNum: 10274540, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10274783, SafeLow: decimal.NewFromFloat(33)},
		{BlockNum: 10275055, SafeLow: decimal.NewFromFloat(29)},
		{BlockNum: 10275347, SafeLow: decimal.NewFromFloat(21)},
		{BlockNum: 10275599, SafeLow: decimal.NewFromFloat(23.1)},
		{BlockNum: 10275866, SafeLow: decimal.NewFromFloat(6)},
		{BlockNum: 10276149, SafeLow: decimal.NewFromFloat(31)},
		{BlockNum: 10276433, SafeLow: decimal.NewFromFloat(29)},
		{BlockNum: 10276664, SafeLow: decimal.NewFromFloat(38)},
		{BlockNum: 10276941, SafeLow: decimal.NewFromFloat(33)},
		{BlockNum: 10277211, SafeLow: decimal.NewFromFloat(45)},
		{BlockNum: 10277484, SafeLow: decimal.NewFromFloat(39)},
		{BlockNum: 10277766, SafeLow: decimal.NewFromFloat(31)},
		{BlockNum: 10278060, SafeLow: decimal.NewFromFloat(23.4)},
		{BlockNum: 10278345, SafeLow: decimal.NewFromFloat(22.9)},
		{BlockNum: 10278626, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10278880, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10279146, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10279389, SafeLow: decimal.NewFromFloat(23)},
		{BlockNum: 10279660, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10279946, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10280227, SafeLow: decimal.NewFromFloat(21.3)},
		{BlockNum: 10280488, SafeLow: decimal.NewFromFloat(21.3)},
		{BlockNum: 10280766, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10281020, SafeLow: decimal.NewFromFloat(28)},
		{BlockNum: 10281300, SafeLow: decimal.NewFromFloat(29)},
		{BlockNum: 10281538, SafeLow: decimal.NewFromFloat(35)},
		{BlockNum: 10281805, SafeLow: decimal.NewFromFloat(34)},
		{BlockNum: 10282079, SafeLow: decimal.NewFromFloat(31)},
		{BlockNum: 10282353, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10282625, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10282915, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10283187, SafeLow: decimal.NewFromFloat(35)},
		{BlockNum: 10283481, SafeLow: decimal.NewFromFloat(36)},
		{BlockNum: 10283760, SafeLow: decimal.NewFromFloat(42)},
		{BlockNum: 10284003, SafeLow: decimal.NewFromFloat(48)},
		{BlockNum: 10284270, SafeLow: decimal.NewFromFloat(35)},
		{BlockNum: 10284532, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10284804, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10285081, SafeLow: decimal.NewFromFloat(21)},
		{BlockNum: 10285329, SafeLow: decimal.NewFromFloat(25)},
		{BlockNum: 10285613, SafeLow: decimal.NewFromFloat(19.4)},
		{BlockNum: 10285881, SafeLow: decimal.NewFromFloat(25)},
		{BlockNum: 10286158, SafeLow: decimal.NewFromFloat(22.1)},
		{BlockNum: 10286400, SafeLow: decimal.NewFromFloat(26)},
		{BlockNum: 10286695, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10286988, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10287270, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10287557, SafeLow: decimal.NewFromFloat(24.1)},
		{BlockNum: 10287822, SafeLow: decimal.NewFromFloat(24.2)},
		{BlockNum: 10288066, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10288358, SafeLow: decimal.NewFromFloat(33)},
		{BlockNum: 10288627, SafeLow: decimal.NewFromFloat(6)},
		{BlockNum: 10288887, SafeLow: decimal.NewFromFloat(28)},
		{BlockNum: 10289154, SafeLow: decimal.NewFromFloat(32)},
		{BlockNum: 10289431, SafeLow: decimal.NewFromFloat(22.9)},
		{BlockNum: 10289694, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10289963, SafeLow: decimal.NewFromFloat(31)},
		{BlockNum: 10290235, SafeLow: decimal.NewFromFloat(37)},
		{BlockNum: 10290485, SafeLow: decimal.NewFromFloat(37)},
		{BlockNum: 10290768, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10291041, SafeLow: decimal.NewFromFloat(28)},
		{BlockNum: 10291327, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10291604, SafeLow: decimal.NewFromFloat(20)},
		{BlockNum: 10291874, SafeLow: decimal.NewFromFloat(21.7)},
		{BlockNum: 10292126, SafeLow: decimal.NewFromFloat(24.6)},
		{BlockNum: 10292396, SafeLow: decimal.NewFromFloat(26)},
		{BlockNum: 10292655, SafeLow: decimal.NewFromFloat(25)},
		{BlockNum: 10292940, SafeLow: decimal.NewFromFloat(20.9)},
		{BlockNum: 10293232, SafeLow: decimal.NewFromFloat(21)},
		{BlockNum: 10293492, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10293771, SafeLow: decimal.NewFromFloat(30)},
		{BlockNum: 10294028, SafeLow: decimal.NewFromFloat(31)},
		{BlockNum: 10294318, SafeLow: decimal.NewFromFloat(23)},
		{BlockNum: 10294597, SafeLow: decimal.NewFromFloat(31)},
		{BlockNum: 10294849, SafeLow: decimal.NewFromFloat(35)},
		{BlockNum: 10295128, SafeLow: decimal.NewFromFloat(28)},
		{BlockNum: 10295429, SafeLow: decimal.NewFromFloat(32)},
		{BlockNum: 10295700, SafeLow: decimal.NewFromFloat(65)},
		{BlockNum: 10295974, SafeLow: decimal.NewFromFloat(55)},
		{BlockNum: 10296247, SafeLow: decimal.NewFromFloat(33)},
		{BlockNum: 10296516, SafeLow: decimal.NewFromFloat(37)},
		{BlockNum: 10296760, SafeLow: decimal.NewFromFloat(38)},
		{BlockNum: 10297025, SafeLow: decimal.NewFromFloat(40)},
		{BlockNum: 10297332, SafeLow: decimal.NewFromFloat(27)},
		{BlockNum: 10297587, SafeLow: decimal.NewFromFloat(30)},
	}
	tim := time.Now()

	for pos, dat := range data {
		dat.DateAdded = tim
		data[pos] = dat
		tim = tim.Add(30 * time.Minute)
	}
	w, _ := os.Create("blogmap.png")
	defer w.Close()
	m := buildMap(&data)
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.
	t.Fail()
}

func TestDecimalRound(t *testing.T) {
	ten := decimal.NewFromInt(10)

	num := decimal.NewFromInt(75)
	res := num.DivRound(ten, 0)
	t.Log(res.String())
	t.Fail()
}

// cmd/fontgen/fontgen -x=1 -y=20 -h=8 -a="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" -img examples/minecraftia.png > minecraftia.txt