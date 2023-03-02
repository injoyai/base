package bytes

import "testing"

func TestBytes_HEXToInt(t *testing.T) {
	t.Log(Entity{0x0a, 0x0b}.HEXToInt())
}

func TestBytes_Sub0x33ReverseHEXToFloat64(t *testing.T) {
	t.Log(Entity{0x3b, 0x3a}.Sub0x33ReverseHEXToInt())
}

func TestBytes_Int64(t *testing.T) {
	{
		x := Entity{0x01, 0x10}.Int64()
		t.Log(x)
		if x != 272 {
			t.Log("错误")
		}
	}
	{
		//读取前8字节,后面的会舍弃
		x := Entity{0x00, 0, 00, 0, 0, 0, 0x01, 0x10, 0, 0}.Int64()
		t.Log(x)
		if x != 272 {
			t.Log("错误")
		}
	}

}

func TestEntity_Int64(t *testing.T) {
	t.Log(Entity{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}.Int64())
	t.Log(Entity{1, 2}.Int64())
}
