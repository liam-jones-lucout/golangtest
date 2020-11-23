package db

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liam-jones-lucout/golangtest/internal/pkg/spaceshipmodels"
	"go.uber.org/zap"
)

type Db struct {
	logger *zap.Logger
	db     orm.Ormer
}

func NewDb(logger *zap.Logger) *Db {
	return &Db{
		logger: logger,
	}
}

func (s *Db) Initiate() error {

	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		s.logger.Error("Error registering sqlDriver", zap.Error(err))
		return err
	}
	if err := orm.RegisterDataBase("default", "mysql", "root:bob@tcp(127.0.0.1:3306)/starwars"); err != nil {
		s.logger.Error("Error Registering database with ORM", zap.Error(err))
	}

	orm.RegisterModel(new(spaceshipmodels.Armament), new(spaceshipmodels.Spaceship))

	if err := orm.RunSyncdb("default", true, false); err != nil {
		s.logger.Error("Error Syncing database with ORM", zap.Error(err))
	}

	db := orm.NewOrm()

	baseData := GetBaseData()
	if _, err := db.InsertMulti(len(baseData), baseData); err != nil {
		s.logger.Error("Error populating database with basedata", zap.Error(err))
	}

	s.db = db
	return nil
}

func (s *Db) Search(name, class, status []string) (spaceshipmodels.Spaceships, error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return nil, err
	}
	qb = qb.Select("id", "name", "status").From("spaceship")

	if len(name) > 0 {
		qb = qb.Where("name IN " + inCreator(name))
	}

	if len(class) > 0 {
		qb = qb.Where("class IN " + inCreator(class))
	}

	if len(status) > 0 {
		qb = qb.Where("status IN " + inCreator(status))
	}

	var ships spaceshipmodels.Spaceships

	if _, err := s.db.Raw(qb.String()).QueryRows(&ships); err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return nil, err
	}

	return ships, nil
}

func (s *Db) Get(id string) (spaceshipmodels.Spaceship, error) {
	ship := spaceshipmodels.Spaceship{}
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return ship, err
	}
	qb = qb.Select("*").From("spaceship").Where("id = " + id)

	if err := s.db.Raw(qb.String()).QueryRow(&ship); err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return ship, err
	}

	if _, err := s.db.LoadRelated(&ship, "Armament"); err != nil {
		s.logger.Error("Failed to load related fields", zap.Error(err))
		return ship, err
	}

	return ship, nil
}

func (s *Db) Delete(id string) error {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return err
	}
	qb = qb.Delete("*").From("spaceship").Where("id = " + id)

	if _, err := s.db.Raw(qb.String()).Exec(); err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return err
	}

	return nil
}

func (s *Db) Update(id string, ship spaceshipmodels.Spaceship) error {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return err
	}
	qb = qb.Update("spaceship").Set(
		"name = "+ship.Name,
		"class = "+ship.Class,
		"crew = "+strconv.Itoa(ship.Crew),
		"image = "+ship.Image,
		"value = "+strconv.Itoa(ship.Value),
		"status = "+ship.Status,
	).Where("id = " + id)

	if _, err := s.db.Raw(qb.String()).Exec(); err != nil {
		s.logger.Error("Failed to query databse", zap.Error(err))
		return err
	}

	return nil
}

func inCreator(stringValues []string) string {
	return "('" + strings.Join(stringValues, "','") + "')"
}
