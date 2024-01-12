package MongoDB

import "strings"

type Model struct {
  Table string
}

func (M Model) New() *Orm {
      dbtag := "mongodb"
      table := M.Table
      index := strings.Index(table, ".")
      if index>0 {
            dbtag = table[0:index]
            table = table[index+1:]
      }

      return NewOrm(table, dbtag)
}

func (M Model) Insert(data map[string]any) string {
  return M.New().Insert(data)
}

func (M Model) Delete() int64 {
    return M.New().Delete()
}

func (M Model) Update(data map[string]any) int64 {
      return M.New().Update(data)
}

func (M Model) Field(fields string) *Orm {
      return M.New().Field(fields)
}

func (M Model) Where(conds ...any) *Orm {
      return M.New().Where(conds...)
}

func (M Model) Group(fields ...string) *Orm {
      return M.New().Group(fields...)
}

func (M Model) Having(field string, opr string, criteria int) *Orm {
      return M.New().Having(field, opr, criteria)
}

func (M Model) Order(field string, sort string) *Orm {
      return M.New().Order(field, sort)
}

func (M Model) Page(page int64) *Orm {
      return M.New().Page(page)
}

func (M Model) Limit(limit int64) *Orm {
      return M.New().Limit(limit)
}

func (M Model) Select() []map[string]any {
      return M.New().Select()
}

func (M Model) Find() map[string]any {
      return M.New().Find()
}

func (M Model) Value(field string) any {
      return M.New().Value(field)
}

func (M Model) Values(field string) []any {
      return M.New().Values(field)
}

func (M Model) Columns(fields ...string) map[string]any {
      return M.New().Columns(fields...)
}

func (M Model) Sum(field string) int32 {
      return M.New().Sum(field)
}

func (M Model) Count() int32 {
      return M.New().Count()
}

func (M Model) Exist(id string) bool {
      return M.New().Exist(id)
}
