package web

type Echo struct {
}

/*func JwtConfig() echojwt.Config {
	secretKey := config.GetConfig().AppJwt.AccessSecret
	config := echojwt.Config{
		SigningKey: []byte(secretKey),
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			err := tokenValid(c.Request())
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	}
	return config
}*/
