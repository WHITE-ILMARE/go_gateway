package dao

import (
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dto"
	"gorm.io/gorm"
	"time"
)

type ServiceInfo struct {
	ID          int       `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"加载类型"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete    int       `json:"is_delete" grom:"column:is_delete" description:"是否已经删除"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) PageList(tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize
	query := tx
	query = query.Table(t.TableName()).Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where("service_name like ? or service_desc like ?", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (t *ServiceInfo) ServiceDetail(tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	httpRule := &HttpRule{ServiceID: int64(search.ID)}
	httpRule, err := httpRule.Find(tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tcpRule := &TcpRule{ServiceID: int64(search.ID)}
	tcpRule, err = tcpRule.Find(tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	grpcRule := &GrpcRule{ServiceID: int64(search.ID)}
	grpcRule, err = grpcRule.Find(tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	accessControl := &AccessControl{ServiceID: int64(search.ID)}
	accessControl, err = accessControl.Find(tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	loadBalance := &LoadBalance{ServiceID: int64(search.ID)}
	loadBalance, err = loadBalance.Find(tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

func (t *ServiceInfo) Find(tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	//fmt.Printf("search condition: %+v\n", search)
	out := &ServiceInfo{}
	err := tx.Where(search).First(out).Error
	//fmt.Printf("found: %+v\n", out)
	//fmt.Printf("err: %v\n", err)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceInfo) Save(tx *gorm.DB) error {
	return tx.Save(t).Error
}
