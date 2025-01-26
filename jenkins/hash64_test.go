package jenkins

import (
	"fmt"
	"testing"
)

func TestHash64(t *testing.T) {
	type test_case struct {
		String   string
		Expected uint64
	}

	test_cases := []test_case{
		{"test", 0xa1c5526862861a7a},
		{"test0", 0xf5fa834a110a00d6},
		{"test1", 0xb67649914a45cba3},
		{"test2", 0x42f771e171b6e47f},
		{"test3", 0x1ad5f392b358a33b},
		{"test4", 0x84e261d44dd8108e},
		{"test5", 0x8b76d034389cdc82},
		{"test6", 0x67c88e7bdac3350b},
		{"test7", 0x56a531cd29e0c07f},
		{"test8", 0xa11ea616715fbc18},
		{"test9", 0x307abdfeda9bb620},
		{"test10", 0x26aeae8e0142b3fd},
		{"test11", 0xd2b6595564679de3},
		{"test12", 0xdf2467b0c7d585a9},
		{"test13", 0x6c6f5693bf61f68d},
		{"test14", 0x53f60f7d7b0bf2cb},
		{"test15", 0x3488a1495ae0d6f1},
		{"test16", 0xf7335137c69aa0d4},
		{"test17", 0x09885b733bd12916},
		{"test18", 0xc8f407be1966647d},
		{"test19", 0xb6e5d612efef2dad},
		{"test20", 0x66db6b1cff86e83b},
		{"test21", 0xe3a0b24a9d7af1c5},
		{"test22", 0x7685acc637098107},
		{"test23", 0x08fdd06bbfdf59de},
		{"test24", 0x59b89475614a76f5},
		{"test25", 0x4bf291958a865854},
		{"test26", 0xf3f32b134af30f89},
		{"test27", 0x2c5d1ab492727f09},
		{"test28", 0xf2b864101600677a},
		{"test29", 0xc37c94dec60ed0a8},
		{"test30", 0x7990b5aba1c4c435},
		{"test31", 0x65dc850ca9f8c757},
		{"test32", 0x15e63787c0859f34},
		{"test33", 0x33a4ddb1290a3309},
		{"test34", 0xe27a9cefc51b91d0},
		{"test35", 0x7b46f5740c4dc5c1},
		{"test36", 0x507bcd5fc4bfd0d4},
		{"test37", 0xf745ee90540ac88c},
		{"test38", 0xeea77a088c2742e3},
		{"test39", 0x93e5a5fc90336fe4},
		{"test40", 0x837fab401eef4dbc},
		{"test41", 0x3335c61e9ba8536b},
		{"test42", 0x53baf5aeabfe8337},
		{"test43", 0x1ba05e9f1b69cbfb},
		{"test44", 0x23865d8ba26a3e25},
		{"test45", 0xd7ce4d57acfd60c9},
		{"test46", 0x31d83a42d7ee57fd},
		{"test47", 0x5425c2f0cd42e5c6},
		{"test48", 0x18950f93b871c220},
		{"test49", 0x8dd0e84921c36032},
		{"test50", 0xda9a30f7520d5469},
		{"test51", 0xd5038c74127cc6ea},
		{"test52", 0x27008618de491543},
		{"test53", 0x214ebefb4eae488c},
		{"test54", 0xae334afe7c9ac824},
		{"test55", 0xb14fa49b301e6cab},
		{"test56", 0xb9f354ee5bd40499},
		{"test57", 0xc4a378054fc7d48e},
		{"test58", 0x9af18aca95b609be},
		{"test59", 0x496711399b49978e},
		{"test60", 0x0774f806626ed615},
		{"test61", 0xe929378c8807dbdb},
		{"test62", 0x724a94f9152b1203},
		{"test63", 0x4e13740653cc4ed7},
		{"(LISTFILE)", 0xd8598604e289efef},
		{"SOUND\\AMBIENCE\\GHOSTSTATE.OGG", 0xa55f36ff5b11d658},
		{"sound\\music\\cataclysm\\mus_arathihighlandsb_gd01.mp3", 0xad9b081e753d99d2},
		{"sound\\music\\cataclysm\\mus_winterspring_gu05.mp3", 0xfb83690c2acd51a4},
	}

	for i, tc := range test_cases {
		result := Hash64([]byte(tc.String))
		if result != tc.Expected {
			t.Fatal("result", i, "was", fmt.Sprintf("%016x", result), "expected", fmt.Sprintf("%016x", tc.Expected))
		}
	}
}
