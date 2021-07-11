package casbinquery

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
	// "github.com/pecolynx/casbin-query/gateway"
)

type petEntity struct {
	ID        uint
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (e *petEntity) TableName() string {
	return "pet"
}

func findPets(db *gorm.DB, name string) ([]string, error) {
	objectColumnName := "name"
	subQuery, err := QueryObject(db, objectColumnName, "user_"+name, "read")
	if err != nil {
		return nil, err
	}

	petEntities := []petEntity{}
	if result := db.Model(&petEntity{}).
		Joins("inner join (?) AS t3 ON `pet`.`name`= t3."+objectColumnName, subQuery).
		Order("`pet`.`name`").
		Scan(&petEntities); result.Error != nil {
		return nil, result.Error
	}
	petNames := make([]string, len(petEntities))
	for i, e := range petEntities {
		petNames[i] = e.Name
	}

	return petNames, nil
}
func TestA(t*testing.T){
	
}

func TestQueryObject(t *testing.T) {
	db := openMySQLForTest()

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "bob",
			args:    args{name: "bob"},
			want:    []string{"ewok", "fluffy"},
			wantErr: false,
		},
		{
			name:    "charlie",
			args:    args{name: "charlie"},
			want:    []string{"gordo"},
			wantErr: false,
		},
		{
			name:    "david",
			args:    args{name: "david"},
			want:    []string{"ewok"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findPets(db, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
