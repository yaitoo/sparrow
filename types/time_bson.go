package types

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

//func (t *Time) UnmarshalBSONValue(b bsontype.Type, data []byte) error {
//	x, y, z := bsoncore.ReadDateTime(data)
//
//	fmt.Println(x, y, z)
//	return nil
//}

func (t *Time) MarshalBSON() ([]byte, error) {
	return nil, nil
}

func (t *Time) UnmarshalBSON(data []byte) error {
	if data == nil || len(data) == 0 {
		t.valid = false
		return nil
	}

	tm, _, ok := bsoncore.ReadTime(data)
	if ok {
		t.Time = tm
		t.valid = true
	} else {
		t.valid = false
	}

	return nil
}
