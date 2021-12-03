package repositories

import (
	"database/sql"
	"shopping_system/common"
	datamodels "shopping_system/dataModels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table: table, mysqlConn: sql}
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "`order`"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT " + o.table + " set userID=?,productID=?,orderStatus=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil {
		return productID, errStmt
	}
	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderID int64) (isOK bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sql := "delete from " + o.table + " where ID =?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil {
		return
	}
	_, err := stmt.Exec(orderID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}
	sql := "Update " + o.table + " set userID=?,productID=?,orderStatus=? where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil {
		return errStmt
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return
}

func (o *OrderManagerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}
	sql := "Select * from " + o.table + "where ID=" + strconv.FormatInt(orderID, 10)
	row, errRow := o.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select * from " + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	defer rows.Close()
	if errRows != nil {
		return nil, errRows
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select o.ID,p.productName,o.OrderStatus From testP.order as o left join product as p on o.productID=p.ID"
	rows, errRows := o.mysqlConn.Query(sql)
	defer rows.Close()
	if errRows != nil {
		return nil, errRows
	}
	return common.GetResultRows(rows), err
}
