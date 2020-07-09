package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type PlanRepository interface {
	Get(name string) (model.Plan, error)
	List() ([]model.Plan, error)
	Page(num, size int) (int, []model.Plan, error)
	Save(plan *model.Plan, zones []string) error
	Delete(name string) error
	Batch(operation string, items []model.Plan) error
}

func NewPlanRepository() PlanRepository {
	return &planRepository{}
}

type planRepository struct {
}

func (p planRepository) Get(name string) (model.Plan, error) {
	var plan model.Plan
	plan.Name = name
	if err := db.DB.Where(plan).First(&plan).Error; err != nil {
		return plan, err
	}
	return plan, nil
}

func (p planRepository) List() ([]model.Plan, error) {
	var plans []model.Plan
	err := db.DB.Model(model.Zone{}).Find(&plans).Error
	return plans, err
}

func (p planRepository) Page(num, size int) (int, []model.Plan, error) {
	var total int
	var plans []model.Plan
	err := db.DB.Model(model.Plan{}).
		Count(&total).
		Find(&plans).
		Offset((num - 1) * size).
		Limit(size).
		Error

	for i, p := range plans {
		var zoneIds []string
		var planZones []model.PlanZones
		db.DB.Model(model.PlanZones{}).Where("plan_id = ?", p.ID).Find(&planZones)
		for _, p := range planZones {
			zoneIds = append(zoneIds, p.ZoneID)
		}
		var zones []model.Zone
		db.DB.Model(model.Zone{}).Where("id in (?)", zoneIds).Find(&zones)
		plans[i].Zones = zones
		var regionIds []string
		for _, z := range zones {
			regionIds = append(regionIds, z.RegionID)
		}
		var regions []model.Region
		db.DB.Model(model.Region{}).Where("id in (?)", regionIds).Find(&regions)
		plans[i].Regions = regions
	}

	return total, plans, err
}

func (p planRepository) Save(plan *model.Plan, zones []string) error {
	if db.DB.NewRecord(plan) {
		tx := db.DB.Begin()
		err := tx.Create(&plan).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, z := range zones {
			err = tx.Create(&model.PlanZones{
				PlanID: plan.ID,
				ZoneID: z,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		return err
	} else {
		tx := db.DB.Begin()
		err := db.DB.Save(&plan).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("plan_id = ?", plan.ID).Delete(&model.PlanZones{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, z := range zones {
			err = tx.Create(model.PlanZones{
				PlanID: plan.ID,
				ZoneID: z,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		return err
	}
}

func (p planRepository) Delete(name string) error {
	plan, err := p.Get(name)
	if err != nil {
		return err
	}
	tx := db.DB.Begin()
	err = tx.Delete(&plan).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("plan_id = ?", plan.ID).Delete(&model.PlanZones{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (p planRepository) Batch(operation string, items []model.Plan) error {
	switch operation {
	case constant.BatchOperationDelete:
		//TODO 关联校验
		//var clusterIds []string
		//for _, item := range items {
		//	clusterIds = append(clusterIds, item.ClusterID)
		//}
		//var clusters []model.Cluster
		//err := db.DB.Where("id in (?)", clusterIds).Find(&clusters).Error
		//if err != nil {
		//	return err
		//}
		//if len(clusters) > 0 {
		//	return errors.New(DeleteFailedError)
		//}
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}

		tx := db.DB.Begin()
		err := db.DB.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("plan_id in (?)", ids).Delete(&model.PlanZones{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()

	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
