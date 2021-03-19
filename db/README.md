# 如何判断缺少参数?参数不存在的化，我要忽略该段SQL，而不需要报错要怎么处理？
``` golang
    query = ctx.NewQuery(context.TODO()).
    For("deposit").
    Select("account_number", "card_number", "card_name", "deposit_amount").
    Where("id={id} and status=0")
//query.Var("id", "1")

	d := struct {
		AccountNumber string
		CardNumber    string
		CardName      string
		DepositAmount float64
	}{}

	if err := query.Find(&d); err != nil {
		if errors.Is(err, parser.ErrParameterMissing) {
			//TODO: 因为没有给query添加id参数值，所以执行后会报错
		}
    }
```

但是当做查询接口的时候，使用者并不会存入全部参数，那我们应该如何处理？使用以下的写法，将避免执行时抛出错误

``` golang
    query = ctx.NewQuery(context.TODO()).
    For("deposit").
    Select("account_number", "card_number", "card_name", "deposit_amount").
    Where("id={id} and status={status}","id","status")//使用第二个参数varNames， 当query当Vars里没有id或者status的话，这段Where不会被包含在最后执行的SQL里。从而避免抛出ErrParameterMissing



	d := struct {
		AccountNumber string
		CardNumber    string
		CardName      string
		DepositAmount float64
	}{}

	if err := query.Find(&d); err != nil {
		if errors.Is(err, parser.ErrParameterMissing) {
			//TODO: 不会抛出ErrParameterMissing， 因为Where已经被跳过了
		}
    }
```