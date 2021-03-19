package mgo

// import (
// 	"sparrow/sparrow/micro"
// 	"sparrow/sparrow/types"
// 	"time"

// 	"gitlab.paradise-soft.com.tw/backend/yaitoo/cfg"

// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // [mongo]
// // host=mongodb://10.200.252.219:27017
// // dbname=hb
// // password=
// // user=

// //MgoSessionManager Mongo线程管理器
// var mgoSessionManager mgo.SessionManager

// //GetMgoSessionManager 读取最新Mongo线程管理器
// func GetMgoSessionManager() mgo.SessionManager {
// 	return mgoSessionManager
// }

// var mgoDatbase string

// //GetMgoDatbase 读取Mongo数据库名称
// func GetMgoDatbase() string {
// 	return mgoDatbase
// }

// func initMgo(c *cfg.Config) {
// 	dsn := mgo.NewURI(c.GetValue("mongo", "host", "127.0.0.1:27017"),
// 		c.GetValue("mongo", "database", ""),
// 		c.GetValue("mongo", "login", ""),
// 		c.GetValue("mongo", "passwd", ""))

// 	mgoDatbase = dsn.Database

// 	maxIdleConns := types.Atoi(c.GetValue("mongo", "max_idle_conns", "10"), 10)
// 	maxOpenConns := types.Atoi(c.GetValue("mongo", "max_open_conns", "100"), 100)

// 	waitTimeout := time.Duration(types.Atoi(c.GetValue("mongo", "wait_timeout", "1"), 1)) * time.Second

// 	var err error
// 	mgoSessionManager, err = mgo.NewSessionManager(micro.NewContext(), maxIdleConns, maxOpenConns, waitTimeout, options.Client().ApplyURI(dsn.String()))

// 	if err != nil {
// 		logger.Errorln(err)
// 	}
// }

// func (m *roundStatsManager) GetMemberOrderStats(ctx context.Context, product, number string) ([]*models.MemberOrderStats, error) {
// 	session, err := GetMgoSessionManager().Get()
// 	if session != nil {
// 		defer session.Close()
// 	}
// 	if err != nil {
// 		logger.Errorln(err)
// 		return nil, err
// 	}

// 	productItems := session.Client.Database(cmd.GetMgoDatbase()).Collection(product + "_items")

// 	filter := bson.M{"num": number}

// 	cur, err := productItems.Find(ctx, filter)
// 	if err != nil {
// 		logger.Errorln(err)
// 		return nil, err
// 	}
// 	defer cur.Close(ctx)

// 	list := make([]*models.MemberOrderStats, 0, 10000)

// 	for cur.Next(ctx) {
// 		// To decode into a struct, use cursor.Decode()
// 		it := &models.MemberOrderStats{}

// 		err := cur.Decode(it)
// 		if err != nil {
// 			logger.Errorln(err, cur.Current)
// 			return nil, err
// 		}

// 		list = append(list, it)
// 		// do something with result...

// 		// To get the raw bson bytes use cursor.Current
// 		//raw := cur.Current
// 		// do something with raw...
// 	}
// 	if err := cur.Err(); err != nil {
// 		return nil, err
// 	}

// 	return list, nil
// }

// func (m *roundStatsManager) Get(ctx context.Context, product, number string) (*models.RoundStats, error) {
// 	session, err := GetMgoSessionManager().Get()
// 	if session != nil {
// 		defer session.Close()
// 	}
// 	if err != nil {
// 		logger.Errorln(err)
// 		return nil, err
// 	}

// 	products := session.Client.Database(cmd.GetMgoDatbase()).Collection(product)

// 	rs := &models.RoundStats{}

// 	filter := bson.M{"num": number}

// 	err = products.FindOne(ctx, filter).Decode(rs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return rs, nil
// }

// func (m *roundStatsManager) Add(ctx context.Context, brand string, co *models.ClientOrder) error {
// 	session, err := GetMgoSessionManager().Get()
// 	if session != nil {
// 		defer session.Close()
// 	}
// 	if err != nil {
// 		logger.Errorln(err)
// 		return err
// 	}

// 	code := strings.ToLower(co.Code)

// 	if types.IsEmpty(code) {
// 		code = strings.ToLower(co.PlayCode + "_" + strings.Replace(co.Value, ",", "_", -1))
// 	}

// 	brand = strings.ToLower(brand)

// 	product := strings.ToLower(co.Product)

// 	//1: increase order item
// 	item := bson.M{
// 		"$set": bson.M{
// 			"num":    co.Number,
// 			"member": brand + ":" + co.Member,
// 			"code":   code,
// 			"up":     time.Now(),
// 		},
// 		"$inc": bson.M{
// 			"qty": 1,
// 			"amt": co.Amount,
// 		},
// 	}

// 	productItems := session.Client.Database(cmd.GetMgoDatbase()).Collection(product + "_items")

// 	_, err = productItems.UpdateOne(ctx,
// 		bson.M{"num": co.Number, "member": brand + ":" + co.Member, "code": code},
// 		item, options.Update().SetUpsert(true))
// 	if err != nil {
// 		logger.Errorf("%s %s %s %s", brand, code, item, err)
// 		return err
// 	}

// 	products := session.Client.Database(cmd.GetMgoDatbase()).Collection(product)

// 	//2: upsert round record
// 	insert := bson.M{
// 		"$setOnInsert": bson.M{
// 			"at":  time.Now(),
// 			"num": co.Number,
// 		}}
// 	_, err = products.UpdateOne(ctx,
// 		bson.M{"num": co.Number},
// 		insert, options.Update().SetUpsert(true))
// 	if err != nil {
// 		logger.Errorf("%s %s %s %s", brand, code, insert, err)
// 		return err
// 	}

// 	//3. update total
// 	update := bson.M{
// 		"$set": bson.M{
// 			"orders." + code + ".code":  co.Code,
// 			"orders." + code + ".name":  co.Name,
// 			"orders." + code + ".value": co.Value,
// 			"orders." + code + ".play":  co.PlayCode,
// 		},

// 		"$inc": bson.M{
// 			"stats.amt":                          co.Amount,
// 			"stats.qty":                          1,
// 			"stats.pay":                          co.Payout,
// 			"stats.clients." + brand + ".amt":    co.Amount,
// 			"stats.clients." + brand + ".qty":    1,
// 			"stats.clients." + brand + ".payout": co.Payout,
// 			"orders." + code + ".total.amt":      co.Amount,
// 			"orders." + code + ".total.qty":      1,
// 			"orders." + code + ".total.payout":   co.Payout,

// 			"orders." + code + ".total.clients." + brand + ".amt":    co.Amount,
// 			"orders." + code + ".total.clients." + brand + ".qty":    1,
// 			"orders." + code + ".total.clients." + brand + ".payout": co.Payout,
// 		},
// 		"$addToSet": bson.M{
// 			"orders." + code + ".total.members": brand + ":" + co.Member,
// 		},
// 	}
// 	_, err = products.UpdateOne(ctx,
// 		bson.M{"num": co.Number},
// 		update, options.Update().SetUpsert(true))

// 	if err != nil {
// 		logger.Errorf("%s %s %s %s", brand, code, update, err)
// 		return err
// 	}

// 	return nil
// }
