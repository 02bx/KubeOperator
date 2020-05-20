package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	adm2 "github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"log"
	"time"
)

func Page(num, size int) (clusters []clusterModel.Cluster, total int, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Find(&clusters).
		Offset((num - 1) * size).
		Limit(size).
		Count(&total).Error
	return
}

func List() (clusters []clusterModel.Cluster, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).Find(&clusters).Error
	return
}

func Get(name string) (*clusterModel.Cluster, error) {
	var result clusterModel.Cluster
	err := db.DB.Model(clusterModel.Cluster{}).
		Where(&result).
		First(&result).
		Error
	return &result, err
}

func Save(item *clusterModel.Cluster) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	var c clusterModel.Cluster
	c.Name = name
	return db.DB.Delete(&c).Error
}

func Batch(operation string, items []clusterModel.Cluster) ([]clusterModel.Cluster, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(clusterModel.Cluster{}).Delete(&item).Error
			if err != nil {
				tx.Rollback()
			}
		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
}

func InitCluster(c clusterModel.Cluster) {
	adm, err := adm2.NewClusterAdm()
	if err != nil {
		log.Fatal(err)
	}
	for {
		//start := time.Now()
		resp, err := adm.OnInitialize(c)
		if err != nil {
		}
		condition := resp.Status.Conditions[len(resp.Status.Conditions)-1]
		switch condition.Status {
		case constant.ConditionFalse:
		case constant.ConditionUnknown:
		case constant.ConditionTrue:
		default:
		}
		time.Sleep(5 * time.Second)
	}
}
