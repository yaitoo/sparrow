
1) Select/SelectModel
 1.1)
    query.For("User").Select("id","name").Where("id={id}", "id").Where("name={name}")
    query.Var("id",id)
    query.Var("name",name)
 1.2)
    query.For("User").SelectModel(&user).Where("id={id}", "id").Where("name={name}")
    query.Var("id",id)
    query.Var("name",name)
  1.3)
    query.RawSQL(sql).Where("id={id}", "id").Where("name={name}")
    query.Var("id",id)
    query.Var("name",name)

query.Find()/query.Query()/query.QueryRow()/query.First()


2) Insert  -Wed

   ctx := databases.Open()

   cmd := ctx.NewCommand();
   
    2.1)
        cmd.For("User").Insert("id",id).Insert("name",name)

    2.2)
        cmd.For("User").Inserts(map[string]interface{}).Insert("name",name)
   


   result,err := cmd.Exec()

 

cmd1.Where("id={id}","id")
cmd1.Var("id",1)
cmd1.For("User").Update("name","dave")


3) Update -Tue

   ctx := databases.Open()

   cmd := ctx.NewCommand();
   
    2.1)

        cmd.For("User").Where("id={id}", "id").Update("name",name)

    2.2)
        cmd.For("User").Where("id={id}", "id").Updates(map[string]interface{}).Update("name",name)
   

    cmd.Var("id",id)

   result,err := cmd.Exec()

4) Delete -Fri
    ctx := databases.Open()

    cmd := ctx.NewCommand();
   

    cmd.For("User").Where("id={id}", "id").Delete()


    cmd.Var("id",id)

   result,err := cmd.Exec()


5) Begin/Commit/Rollback -Fri
    ctx := databases.Open()

    if err := ctx.Begin(); err != nil {
//...
   }
   defer ctx.Rollback()

    cmd1 := ctx.NewCommand();
    id := NextID()

    cmd1.Var("id",id)
    //...
    cmd1.Exec()

    cmd2 := ctx.NewCommand();
    id2 := NewSubNextID(id)
    cmd2.Var("id",id2)
    //...
    cmd2.Exec()

    if err := ctx.Commit(); err != nil {
//...
   }



6) Multi-Databases/Tables
  //... config/dbs/tables




databases := db.NewDatabase(ctx,conns);

type Context struct {
    tx *sql.tx
    connName string
}

func AddUser(){



//cmd1.Updates(map[string]interfacec{}{"id":1, "name":"name"})
  
  
if err := cmd1.Begin(); err != nil {

   }
//ctx.tx =  tx
//ctx.conn =  conn

   defer ctx.Rollback()
   cmd1.Exec()

  //ctx.tx =  
//ctx.conn =  


cmd2 := ctx.NewCommand();
cmd2.RawSQL()

// if ctx.tx != nil {
return err
}
cmd2.Begin()
defer cmd2.Rollback()

ctx.conn != "" && ctx.conn != conn2 {
return err
}
cmd2.Exec()

  if err := ctx.Commit(); err != nil {

  }







query.Var("id",1)

id := ctx.Vars(map)

if id != "" {
    query.Where("id={id}")
}

if col


}

}


1008:
1)  去掉For的Tag變數
2） NewID組合進database，並增加tag變數， For("order").NewID(ctx, tag)/NewSubID(ctx, tag, parentID)
3） config默認位置，NewDatabase(ctx)
4)  columns名稱的規範化， 小寫/移除_
5) 特殊的Error使用const


1017

1) QueryRow/QueryRows
2) Find
   2.1)
   var max int
   query.Find(&max)

   2.2)
   var user User
   query.Find(&user)

   2.3)
   var list []string
   query.Find(&list)

   2.4)
   var users []*User
   query.Take(pageIndex, pageSize)
   query.Find(&users)

   2.5)

   query.Where("id={id",true)
   query.Var("id",id)

   var sumary Sumary{}
   
   query.RawSQL("SELECT count(id) as cnt,sum(amount) as total FROM Users").Find(&sumary)


3) log
    sparrow/sparrow/log

    logger = log.NewLogger("db")

    1) RawSQL (DEBUG)
    -2) SlowSQL ()

1023)
   *1) shardingId.InstanceID 设定
   *2）实际使用场景 
   *3）connection pool 

    





   

