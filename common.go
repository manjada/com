package mjd

import (
	"crypto/rand"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/manjada/com/dto"
	"time"
)

func Copier(to interface{}, source interface{}) {
	copier.CopyWithOption(&to, &source, copyOption())
}

func copyOption() copier.Option {
	const copyInt64 int64 = 0
	return copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)

					if !ok {
						err := fmt.Errorf("src type not matching")
						return nil, err
					}

					return s.Format(dto.DateFormatMediumDash), nil
				},
			},
			{
				SrcType: copier.Int,
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(int)

					if !ok {
						err := fmt.Errorf("src type not matching")
						Error(err)
						return nil, err
					}

					return fmt.Sprintf("%d", s), nil
				},
			},
			{
				SrcType: copier.Float64,
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(float64)

					if !ok {
						err := fmt.Errorf("src type not matching")
						Error(err)
						return nil, err
					}

					return fmt.Sprintf("%.2f", s), nil
				},
			},
			{
				SrcType: copyInt64,
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(int64)

					if !ok {
						err := fmt.Errorf("src type not matching")
						Error(err)
						return nil, err
					}

					return fmt.Sprintf("%d", s), nil
				},
			},
		},
	}
}

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(dto.OtpChar)
	for i := 0; i < length; i++ {
		buffer[i] = dto.OtpChar[int(buffer[i])%otpCharsLength]
	}
	fmt.Println(string(buffer))
	return string(buffer), nil
}
