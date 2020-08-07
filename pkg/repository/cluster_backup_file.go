package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterBackupFileRepository interface {
	Page(num, size int, clusterId string) (int, []model.ClusterBackupFile, error)
	Save(file *model.ClusterBackupFile) error
	Batch(operation string, items []model.ClusterBackupFile) error
}

type clusterBackupFileRepository struct {
}

func NewClusterBackupFileRepository() ClusterBackupFileRepository {
	return &clusterBackupFileRepository{}
}

func (c clusterBackupFileRepository) Page(num, size int, clusterId string) (int, []model.ClusterBackupFile, error) {
	var total int
	var files []model.ClusterBackupFile
	err := db.DB.Model(model.ClusterBackupFile{}).Where(model.ClusterBackupFile{ClusterID: clusterId}).Count(&total).
		Find(&files).Offset((num - 1) * size).Limit(size).Error
	return total, files, err
}

func (c clusterBackupFileRepository) Save(file *model.ClusterBackupFile) error {
	if db.DB.NewRecord(file) {
		return db.DB.Create(&file).Error
	} else {
		return db.DB.Updates(&file).Error
	}
}

func (c clusterBackupFileRepository) Batch(operation string, items []model.ClusterBackupFile) error {

	tx := db.DB.Begin()
	switch operation {
	case constant.BatchOperationDelete:
		for i := range items {

			var file model.ClusterBackupFile
			if err := db.DB.Where(model.ClusterBackupFile{Name: items[i].Name}).First(&file).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := db.DB.Delete(&file).Error; err != nil {
				tx.Rollback()
				return err
			}

		}
	default:
		return constant.NotSupportedBatchOperation
	}
	tx.Commit()
	return nil
}
