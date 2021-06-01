package main

type Keys struct {
}

type Position struct {
	X int64
	Y int64
}

func NewMap() *map[int64]Position {
	m := make(map[int64]Position)

	m[27] = Position{
		X: 0,
		Y: 1,
	}
	m[192] = Position{
		X: 1,
		Y: 1,
	}
	m[9] = Position{
		X: 2,
		Y: 1,
	}
	m[20] = Position{
		X: 3,
		Y: 1,
	}
	m[160] = Position{
		X: 4,
		Y: 1,
	}
	m[168] = Position{
		X: 5,
		Y: 1,
	}

	m[49] = Position{
		X: 1,
		Y: 2,
	}
	m[80] = Position{
		X: 2,
		Y: 2,
	}
	m[65] = Position{
		X: 3,
		Y: 2,
	}
	m[226] = Position{
		X: 4,
		Y: 2,
	}
	m[91] = Position{
		X: 5,
		Y: 2,
	}

	m[112] = Position{
		X: 0,
		Y: 3,
	}
	m[50] = Position{
		X: 1,
		Y: 3,
	}
	m[87] = Position{
		X: 2,
		Y: 3,
	}
	m[83] = Position{
		X: 3,
		Y: 3,
	}
	m[90] = Position{
		X: 4,
		Y: 3,
	}
	m[164] = Position{
		X: 5,
		Y: 3,
	}

	m[113] = Position{
		X: 0,
		Y: 4,
	}
	m[51] = Position{
		X: 1,
		Y: 4,
	}
	m[69] = Position{
		X: 2,
		Y: 4,
	}
	m[68] = Position{
		X: 3,
		Y: 4,
	}
	m[88] = Position{
		X: 4,
		Y: 4,
	}

	m[114] = Position{
		X: 0,
		Y: 5,
	}
	m[52] = Position{
		X: 1,
		Y: 5,
	}
	m[82] = Position{
		X: 2,
		Y: 5,
	}
	m[70] = Position{
		X: 3,
		Y: 5,
	}
	m[67] = Position{
		X: 4,
		Y: 5,
	}

	m[115] = Position{
		X: 0,
		Y: 6,
	}
	m[53] = Position{
		X: 1,
		Y: 6,
	}
	m[84] = Position{
		X: 2,
		Y: 6,
	}
	m[71] = Position{
		X: 3,
		Y: 6,
	}
	m[86] = Position{
		X: 4,
		Y: 6,
	}

	m[116] = Position{
		X: 0,
		Y: 7,
	}
	m[54] = Position{
		X: 1,
		Y: 7,
	}
	m[89] = Position{
		X: 2,
		Y: 7,
	}
	m[71] = Position{
		X: 3,
		Y: 7,
	}
	m[66] = Position{
		X: 4,
		Y: 7,
	}
	m[32] = Position{
		X: 5,
		Y: 7,
	}

	m[117] = Position{
		X: 0,
		Y: 8,
	}
	m[55] = Position{
		X: 1,
		Y: 8,
	}
	m[85] = Position{
		X: 2,
		Y: 8,
	}
	m[74] = Position{
		X: 3,
		Y: 8,
	}
	m[78] = Position{
		X: 4,
		Y: 8,
	}

	m[118] = Position{
		X: 0,
		Y: 9,
	}
	m[56] = Position{
		X: 1,
		Y: 9,
	}
	m[73] = Position{
		X: 2,
		Y: 9,
	}
	m[75] = Position{
		X: 3,
		Y: 9,
	}
	m[77] = Position{
		X: 4,
		Y: 9,
	}

	m[119] = Position{
		X: 0,
		Y: 10,
	}
	m[57] = Position{
		X: 1,
		Y: 10,
	}
	m[79] = Position{
		X: 2,
		Y: 10,
	}
	m[76] = Position{
		X: 3,
		Y: 10,
	}
	m[188] = Position{
		X: 4,
		Y: 10,
	}

	m[120] = Position{
		X: 0,
		Y: 11,
	}
	m[48] = Position{
		X: 1,
		Y: 11,
	}
	m[80] = Position{
		X: 2,
		Y: 11,
	}
	m[86] = Position{
		X: 3,
		Y: 11,
	}
	m[190] = Position{
		X: 4,
		Y: 11,
	}
	m[165] = Position{
		X: 5,
		Y: 11,
	}

	m[121] = Position{
		X: 0,
		Y: 12,
	}
	m[189] = Position{
		X: 1,
		Y: 12,
	}
	m[219] = Position{
		X: 2,
		Y: 12,
	}
	m[222] = Position{
		X: 3,
		Y: 12,
	}
	m[191] = Position{
		X: 4,
		Y: 12,
	}

	m[122] = Position{
		X: 0,
		Y: 13,
	}
	m[187] = Position{
		X: 1,
		Y: 13,
	}
	m[221] = Position{
		X: 2,
		Y: 13,
	}
	m[220] = Position{
		X: 3,
		Y: 13,
	}
	m[93] = Position{
		X: 5,
		Y: 13,
	}

	m[123] = Position{
		X: 0,
		Y: 14,
	}
	m[8] = Position{
		X: 1,
		Y: 14,
	}
	m[13] = Position{
		X: 3,
		Y: 14,
	}
	m[161] = Position{
		X: 4,
		Y: 14,
	}
	m[163] = Position{
		X: 5,
		Y: 14,
	}

	m[44] = Position{
		X: 0,
		Y: 15,
	}
	m[45] = Position{
		X: 1,
		Y: 15,
	}
	m[46] = Position{
		X: 2,
		Y: 15,
	}
	m[37] = Position{
		X: 5,
		Y: 15,
	}

	m[145] = Position{
		X: 0,
		Y: 16,
	}
	m[36] = Position{
		X: 1,
		Y: 16,
	}
	m[35] = Position{
		X: 2,
		Y: 16,
	}
	m[38] = Position{
		X: 4,
		Y: 16,
	}
	m[40] = Position{
		X: 5,
		Y: 16,
	}

	m[19] = Position{
		X: 0,
		Y: 17,
	}
	m[33] = Position{
		X: 1,
		Y: 17,
	}
	m[34] = Position{
		X: 2,
		Y: 17,
	}
	m[39] = Position{
		X: 5,
		Y: 17,
	}

	m[144] = Position{
		X: 1,
		Y: 18,
	}
	m[103] = Position{
		X: 2,
		Y: 18,
	}
	m[100] = Position{
		X: 3,
		Y: 18,
	}
	m[97] = Position{
		X: 4,
		Y: 18,
	}

	m[111] = Position{
		X: 1,
		Y: 19,
	}
	m[104] = Position{
		X: 2,
		Y: 19,
	}
	m[101] = Position{
		X: 3,
		Y: 19,
	}
	m[98] = Position{
		X: 4,
		Y: 19,
	}
	m[45] = Position{
		X: 5,
		Y: 19,
	}

	m[106] = Position{
		X: 1,
		Y: 20,
	}
	m[105] = Position{
		X: 2,
		Y: 20,
	}
	m[102] = Position{
		X: 3,
		Y: 20,
	}
	m[99] = Position{
		X: 4,
		Y: 20,
	}
	m[46] = Position{
		X: 5,
		Y: 20,
	}

	m[109] = Position{
		X: 1,
		Y: 21,
	}
	m[107] = Position{
		X: 2,
		Y: 21,
	}
	m[13] = Position{
		X: 4,
		Y: 21,
	}

	return &m
}
