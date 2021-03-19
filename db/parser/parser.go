package parser

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xwb1989/sqlparser"
	"github.com/yaitoo/sparrow/db/model"
	configmanager "github.com/yaitoo/sparrow/db/model"
	idgenerator "github.com/yaitoo/sparrow/db/shardingId"
	"github.com/yaitoo/sparrow/db/util"
	"github.com/yaitoo/sparrow/micro"
)

var timeNow = time.Now

var (
	regexSQLVarToken    = regexp.MustCompile(`{([\w\d_-])+?}`)
	regexSQLPlaceholder = regexp.MustCompile(`:v\d{1,9}`)
)

type KeyAndValue struct {
	Key   string
	Value idgenerator.IdStruct
}

func (kv *KeyAndValue) IsNullObject() bool {
	return kv.Value.IsNullObject()
}

type AliasEntityKeyPair struct {
	EntityName string
	Alias      string
	Key        KeyAndValue
}

type SqlStringAndEnitityData struct {
	SqlString   string
	PrepareArgs []interface{}
	EntityData  EntityKeyPair
}

func (ss *SqlStringAndEnitityData) GetTarget() AliasEntityKeyPair {
	for e := range ss.EntityData {
		if ss.EntityData[e].Key.Value.IsNullObject() == false {
			return ss.EntityData[e]
		}
	}
	return AliasEntityKeyPair{}
}

func (p *AliasEntityKeyPair) IsNullObject() bool {
	return p.EntityName == ""
}

type EntityKeyPair map[string]AliasEntityKeyPair

// 對entity key paire 做歷尋; 看看是否有entity 或 alias相同的物件
func (ekp *EntityKeyPair) aliasEntityExist(identName string) AliasEntityKeyPair {
	for _, v := range *ekp {
		if v.EntityName == identName || v.Alias == identName {
			return v
		}
	}
	return AliasEntityKeyPair{}
}

// 確認identity(entityName或aliasName) 和key是屬於哪個entity table(反找用)
func (ekp *EntityKeyPair) checkEntityByKeyAndAlias(identName, key string) (bool, AliasEntityKeyPair) {
	e := ekp.aliasEntityExist(identName)

	if e.IsNullObject() {
		for _, v := range *ekp {
			if v.Key.Key == key {
				return true, v
			}
		}
		return false, AliasEntityKeyPair{}
	} else {
		if e.Key.Key == key {
			return true, e
		}
	}

	return false, e
}

func (ekp *EntityKeyPair) CheckAllEntitiesInSameDatabase(ruleConfig configmanager.Config) error {
	var involveSlice = make([]configmanager.Table, 0)
	// 逐個確定該table有沒有在分表規則內
	// 把有的參與的放到一個slice
	//  有多個參與則reject
	for outer := range *ekp {
		outerObejct := (*ekp)[outer]
		tmpEntity := ruleConfig.GetLastVersionEntity(outerObejct.EntityName)
		if tmpEntity.IsNullObject() == false {
			involveSlice = append(involveSlice, tmpEntity)
		}
	}

	if len(involveSlice) > 1 {
		var involvedTable string
		for involved := range involveSlice {
			involvedTable += involveSlice[involved].Name + ", "
		}
		return ErrCrossDatabase(involvedTable)
	}
	return nil
}

func (ekp *EntityKeyPair) SetKey(entityName, key string, value idgenerator.IdStruct) {
	// 根據table identification 去查找是否存過
	entity := ekp.aliasEntityExist(entityName)

	if entity.IsNullObject() == true {
		return
	}

	tmp := entity
	tmp.Key = KeyAndValue{
		Key:   key,
		Value: value,
	}

	(*ekp)[entity.EntityName] = tmp
}

func ParseWithVarMap(ruleConfig configmanager.Config, sql string, args map[string]interface{}) (result SqlStringAndEnitityData, parserErr error) {
	var (
		targetTable string
		firstTable  string
		err         error
		stmt        sqlparser.Statement
	)
	_, prepared, preparedArgs, err := prepare(sql, args)
	if err != nil {
		return result, err
	}
	var entityKeyPair EntityKeyPair = make(map[string]AliasEntityKeyPair)

	tokens := sqlparser.NewTokenizer(strings.NewReader(prepared))

	for {
		if stmt, err = sqlparser.ParseNext(tokens); err == io.EOF {
			break
		}

		switch stmt := stmt.(type) {
		case *sqlparser.Select, *sqlparser.Union, *sqlparser.Update, *sqlparser.Insert, *sqlparser.Delete:
			if err = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
				return getTableNameToEntityWithArgs(node, args, entityKeyPair, &targetTable, &firstTable, ruleConfig)
			}, stmt); err != nil {
				return result, err
			}

			/* 			if err != nil {
				return result, err
			} */
			if targetTable == "" {
				targetTable = firstTable
			}

			sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
				switch node := node.(type) {
				case *sqlparser.ColName:
					colNameByWalk(entityKeyPair, node, ruleConfig, targetTable)
				case *sqlparser.StarExpr:
					starExprWalk(entityKeyPair, node, ruleConfig, targetTable)
				case *sqlparser.AliasedTableExpr:
					aliasedTableExprWalk(entityKeyPair, node, ruleConfig, targetTable)
				case *sqlparser.ComparisonExpr:
					if err = comparisonTableExprWalk(entityKeyPair, node, ruleConfig, targetTable); err != nil {
						return false, err
					}
				case *sqlparser.Insert:
					tableNameWalk(entityKeyPair, node, ruleConfig, targetTable)
				}

				return true, nil
			}, stmt)
		case nil:
			break
		default:
			parserErr = ErrWrongDML
			break
		}
		if err != nil {
			return result, err
		}
		if err = entityKeyPair.CheckAllEntitiesInSameDatabase(ruleConfig); err != nil {
			return result, err
		}
		result = SqlStringAndEnitityData{
			SqlString:   regexSQLPlaceholder.ReplaceAllString(sqlparser.String(stmt), "?"),
			PrepareArgs: preparedArgs,
			EntityData:  entityKeyPair,
		}
	} // 結束多個語句的歷尋
	if parserErr != nil {
		return result, parserErr
	}

	return result, nil
}

// replace qualifier
func colNameByWalk(entityKeyPair EntityKeyPair, node *sqlparser.ColName, ruleConfig configmanager.Config, tableName string) {
	if node.Qualifier.Name.String() == "" {
		return
	}

	if v, ok := entityKeyPair[node.Qualifier.Name.String()]; ok == true {
		if v.Key.Value.IsNullObject() == true {
			return
		}
		table := ruleConfig.GetSpecificVersion(v.Key.Value.AlgVer).GetTable(tableName)
		tag := table.GetTag(v.Key.Value.BusinessNumber)
		psysicalTableName, _ := MapToPhysicalTable(v.Key.Value, table, tag, tableName)
		node.Qualifier.Name = sqlparser.NewTableIdent(psysicalTableName)
	}
}

// repalce the table name of expression
func starExprWalk(entityKeyPair EntityKeyPair, node *sqlparser.StarExpr, ruleConfig configmanager.Config, tableName string) {
	if node.TableName.Name.String() == "" {
		return
	}

	if v, ok := entityKeyPair[node.TableName.Name.String()]; ok == true {
		if v.Key.Value.IsNullObject() == true {
			return
		}
		table := ruleConfig.GetSpecificVersion(v.Key.Value.AlgVer).GetTable(tableName)
		tag := table.GetTag(v.Key.Value.BusinessNumber)
		psysicalTableName, _ := MapToPhysicalTable(v.Key.Value, table, tag, tableName)
		node.TableName.Name = sqlparser.NewTableIdent(psysicalTableName)
	}
}

// check comparison expression is in sharding rules
func comparisonTableExprWalk(entityKeyPair EntityKeyPair, node *sqlparser.ComparisonExpr, ruleConfig configmanager.Config, tableName string) error {
	var (
		kvPair AliasEntityKeyPair
		exist  bool
		table  configmanager.Table
	)
	if node.Operator == "in" && reflect.TypeOf(node.Left).String() == "*sqlparser.ColName" && reflect.TypeOf(node.Right).String() == "sqlparser.ValTuple" {
		left := node.Left.(*sqlparser.ColName)

		if exist, kvPair = entityKeyPair.checkEntityByKeyAndAlias(left.Qualifier.Name.String(), left.Name.String()); exist == true {
			return ErrWrongDML
		}
		if table = ruleConfig.GetLastVersionEntity(kvPair.EntityName); table.IsNullObject() == true {
			return nil
		}
		if table.Key == left.Name.String() {
			return ErrWrongDML
		}
	}
	return nil
}

// replace the table name of the aliased table expression
func aliasedTableExprWalk(entityKeyPair EntityKeyPair, node *sqlparser.AliasedTableExpr, ruleConfig configmanager.Config, tableName string) {
	if _, ok := node.Expr.(*sqlparser.Subquery); ok == true {
		return
	}

	if node.Expr.(sqlparser.TableName).Name.String() == "" {
		return
	}

	if v, ok := entityKeyPair[node.Expr.(sqlparser.TableName).Name.String()]; ok == true {
		if v.Key.Value.IsNullObject() == true {
			return
		}
		switch node.Expr.(type) {
		case *sqlparser.Subquery:
			break
		case sqlparser.TableName:
			table := ruleConfig.GetSpecificVersion(v.Key.Value.AlgVer).GetTable(tableName)
			tag := table.GetTag(v.Key.Value.BusinessNumber)
			psysicalTableName, _ := MapToPhysicalTable(v.Key.Value, table, tag, tableName)
			node.Expr = &sqlparser.TableName{
				Name: sqlparser.NewTableIdent(psysicalTableName),
			}
		}
	}
}

// replace the table name of the aliased table expression
func tableNameWalk(entityKeyPair EntityKeyPair, node *sqlparser.Insert, ruleConfig configmanager.Config, tableName string) {
	if node.Table.Name.String() == "" {
		return
	}

	if v, ok := entityKeyPair[node.Table.Name.String()]; ok == true {
		/* if v.Key.Value.IsNullObject() == true {
			return
		} */
		table := ruleConfig.GetSpecificVersion(v.Key.Value.AlgVer).GetTable(tableName)
		tag := table.GetTag(v.Key.Value.BusinessNumber)
		psysicalTableName, _ := MapToPhysicalTable(v.Key.Value, table, tag, tableName)

		node.Table = sqlparser.TableName{
			Name: sqlparser.NewTableIdent(psysicalTableName),
		}
	}
}

func MapToPhysicalTable(id idgenerator.IdStruct, table model.Table, tag model.Tag, logicalTableName string) (string, error) {
	result := ""
	if tag.IsNullObject() == true {
		return logicalTableName, nil
	}
	time, _ := table.GetTime(id.Time)
	switch tag.Date {
	case "day":
		result = logicalTableName + "_" + tag.GetIdString() + "_" + time.ToYearString() + "_" + util.Left(time.ToMonthString(), 2, "0") + "_" + util.Left(time.ToDayString(), 2, "0")
	case "month":
		result = logicalTableName + "_" + tag.GetIdString() + "_" + time.ToYearString() + "_" + util.Left(time.ToMonthString(), 2, "0")
	case "year":
		result = logicalTableName + "_" + tag.GetIdString() + "_" + time.ToYearString()
	case "week":
		result = logicalTableName + "_" + tag.GetIdString() + "_" + time.ToYearString() + "_" + util.Left(time.ToWeekString(), 2, "0")
	case "year_day":
		result = logicalTableName + "_" + tag.GetIdString() + "_" + time.ToYearString() + "_" + util.Left(time.ToYearDayString(), 3, "0")
	default:
		result = logicalTableName + "_" + tag.GetIdString()
	}

	switch tag.Amount {
	case 0:
		break
	default:
		result = result + "_" + util.Int64Left(tag.GetHashValue(id.Sequence), 2, "0")
	}

	return result, nil
}

func interpolateParams(query string, args []driver.Value) (string, error) {
	// Number of ? should be same to len(args)
	if strings.Count(query, "?") != len(args) {
		return "", driver.ErrSkip
	}

	buf := takeCompleteBuffer()
	buf = buf[:0]
	argPos := 0

	for i := 0; i < len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf = append(buf, query[i:]...)
			break
		}
		buf = append(buf, query[i:i+q]...)
		i += q

		arg := args[argPos]
		argPos++

		if arg == nil {
			buf = append(buf, "NULL"...)
			continue
		}

		switch v := arg.(type) {
		case int64:
			buf = strconv.AppendInt(buf, v, 10)
		case uint64:
			// Handle uint64 explicitly because our custom ConvertValue emits unsigned values
			buf = strconv.AppendUint(buf, v, 10)
		case float64:
			buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
		case bool:
			if v {
				buf = append(buf, '1')
			} else {
				buf = append(buf, '0')
			}
		case time.Time:
			if v.IsZero() {
				buf = append(buf, "'0000-00-00'"...)
			} else {
				v := v.In(time.UTC)
				v = v.Add(time.Nanosecond * 500) // To round under microsecond
				year := v.Year()
				year100 := year / 100
				year1 := year % 100
				month := v.Month()
				day := v.Day()
				hour := v.Hour()
				minute := v.Minute()
				second := v.Second()
				micro := v.Nanosecond() / 1000

				buf = append(buf, []byte{
					'\'',
					digits10[year100], digits01[year100],
					digits10[year1], digits01[year1],
					'-',
					digits10[month], digits01[month],
					'-',
					digits10[day], digits01[day],
					' ',
					digits10[hour], digits01[hour],
					':',
					digits10[minute], digits01[minute],
					':',
					digits10[second], digits01[second],
				}...)

				if micro != 0 {
					micro10000 := micro / 10000
					micro100 := micro / 100 % 100
					micro1 := micro % 100
					buf = append(buf, []byte{
						'.',
						digits10[micro10000], digits01[micro10000],
						digits10[micro100], digits01[micro100],
						digits10[micro1], digits01[micro1],
					}...)
				}
				buf = append(buf, '\'')
			}
		case []byte:
			if v == nil {
				buf = append(buf, "NULL"...)
			} else {
				buf = append(buf, "_binary'"...)
				if statusNoBackslashEscapes == 0 {
					buf = escapeBytesBackslash(buf, v)
				} else {
					buf = escapeBytesQuotes(buf, v)
				}
				buf = append(buf, '\'')
			}
		case string:
			buf = append(buf, '\'')
			if statusNoBackslashEscapes == 0 {
				buf = escapeStringBackslash(buf, v)
			} else {
				buf = escapeStringQuotes(buf, v)
			}
			buf = append(buf, '\'')
		default:
			return "", driver.ErrSkip
		}

		if len(buf)+4 > maxAllowedPacket {
			return "", driver.ErrSkip
		}
	}
	if argPos != len(args) {
		return "", driver.ErrSkip
	}
	return string(buf), nil
}

// takeCompleteBuffer returns the complete existing buffer.
// This can be used if the necessary buffer size is unknown.
// cap and len of the returned buffer will be equal.
// Only one buffer (total) can be used at a time.
func takeCompleteBuffer() []byte {
	return make([]byte, defaultBufSize)
}

func convertAndInterpolateParams(sql string, args []interface{}) (string, error) {
	if len(args) == 0 {
		return sql, nil
	}
	argValue := make([]driver.Value, len(args))
	for arg := range args {
		dbValue, err := driver.DefaultParameterConverter.ConvertValue(args[arg])
		if err != nil {
			return "", err
		}
		argValue[arg] = dbValue
	}
	temp, err := interpolateParams(sql, argValue)
	if err != nil {
		return "", err
	}
	return temp, nil
}

func prepare(cmd string, args map[string]interface{}) (string, string, []interface{}, error) {
	// get all placeholder and compare amount with args
	tokens := regexSQLVarToken.FindAllStringIndex(cmd, -1)
	if args == nil {
		args = make(map[string]interface{})
	}
	//var err error

	//if n > 0 {
	s := ""
	params := make([]interface{}, 0)
	i := 0

	for _, v := range tokens {
		s += cmd[i:v[0]] + "?"
		name := cmd[v[0]+1 : v[1]-1]

		val, ok := args[strings.ToLower(name)]
		if ok {
			//args[k] = val
			params = append(params, val)
			//	formattedCMD += cmd[i:v[0]] + fmt.Sprintf("%s", val)
		} else {
			return cmd, "", nil, micro.Throw(context.TODO(), ErrParameterMissing, name)
		}
		i = v[1]
	}

	if i < len(cmd) {
		s += cmd[i:]
		//formattedCMD += cmd[i:]
	}
	//stmt, err = c.db.Prepare(s)
	/* 		if c.tx != nil {
	   			stmt, err = c.tx.Prepare(s)
	   		} else {
	   			stmt, err = c.db.Prepare(s)
	   		} */
	convertedSqlStr, err := convertAndInterpolateParams(s, params)
	return convertedSqlStr, s, params, err
	//}
	return cmd, cmd, nil, nil
}

func getTableNameToEntityWithArgs(node sqlparser.SQLNode, args map[string]interface{}, entityKeyPair EntityKeyPair, targetTable, firstTable *string, ruleConfig configmanager.Config) (bool, error) {
	switch node := node.(type) {
	case *sqlparser.AliasedTableExpr:
		if _, ok := node.Expr.(*sqlparser.Subquery); ok == true {
			return true, nil
		}
		if node.Expr.(sqlparser.TableName).Name.String() == "" {
			break
		}
		tempTableNm := node.Expr.(sqlparser.TableName).Name.String()
		tableConfig := ruleConfig.GetLastVersionEntity(tempTableNm)
		if *firstTable == "" {
			*firstTable = tempTableNm
		}
		if *targetTable == "" && (&tableConfig).IsNullObject() == false {
			*targetTable = tempTableNm
		}
		_, ok := entityKeyPair[tempTableNm]
		if ok == false {
			entityKeyPair[tempTableNm] = AliasEntityKeyPair{
				EntityName: tempTableNm,
				Alias:      node.As.String(),
			}
		}
	case *sqlparser.ComparisonExpr:
		if reflect.TypeOf(node.Left).String() == "*sqlparser.ColName" && reflect.TypeOf(node.Right).String() == "*sqlparser.SQLVal" {
			left := node.Left.(*sqlparser.ColName)
			rightValue := args[left.Name.String()]
			idObject := idgenerator.ParseStringId(fmt.Sprintf("%v", rightValue))
			//如果sqlval能夠成功被解析成generic id
			if idObject.IsNullObject() == false {
				entityKeyPair.SetKey(left.Qualifier.Name.String(), left.Name.String(), idObject)
			}
			// 不能被解析成generic id; 表示該表應該沒參與分表規則.
			//  不做特別操作
		}
	case *sqlparser.Insert:
		tempTableNm := node.Table.Name.String()
		tableConfig := ruleConfig.GetLastVersionEntity(tempTableNm)
		if *firstTable == "" {
			*firstTable = tempTableNm
		}
		if *targetTable == "" && (&tableConfig).IsNullObject() == false {
			*targetTable = tempTableNm
		}
		_, ok := entityKeyPair[tempTableNm]
		if ok == false {
			for col := range node.Columns {
				colName := node.Columns[col].CompliantName()
				if arg, ok := args[colName]; ok {
					idObject := idgenerator.ParseStringId(fmt.Sprintf("%v", arg))
					if idObject.IsNullObject() {
						continue
					}
					//entityKeyPair.SetKey(tempTableNm, node.Columns[col].CompliantName(), idObject)
					entityKeyPair[tempTableNm] = AliasEntityKeyPair{
						EntityName: tempTableNm,
						Alias:      "",
						Key: KeyAndValue{
							Key:   colName,
							Value: idObject,
						},
					}
				}
			}
		}

	default:
		return true, nil
	}
	return true, nil
}
