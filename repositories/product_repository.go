package repositories

import (
	"database/sql"
	"shopping_system/common"
	datamodels "shopping_system/dataModels"
	"strconv"
)

//先开发对应的接口
//实现定义的接口

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	//判断sql
	if err = p.Conn(); err != nil {
		return
	}

	//准备sql
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}
	defer stmt.Close()

	//传入参数
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
	if err != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	//gorm
	sql := "Update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductManager) SelectByKey(productId int64) (*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID = " + strconv.FormatInt(productId, 10)
	row, err := p.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Product{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult := datamodels.Product{}
	common.DataToStructByTagSql(result, &productResult)
	return &productResult, nil
}

func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "Select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}
