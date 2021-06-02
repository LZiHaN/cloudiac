package services

import (
	"bufio"
	"cloudiac/consts"
	"cloudiac/consts/e"
	"cloudiac/libs/db"
	"cloudiac/models"
	"cloudiac/models/forms"
	"cloudiac/utils"
	"cloudiac/utils/logs"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func CreateTask(tx *db.Session, task models.Task) (*models.Task, e.Error) {
	if err := models.Create(tx, &task); err != nil {
		if e.IsDuplicate(err) {
			return nil, e.New(e.TaskAlreadyExists, err)
		}
		return nil, e.New(e.DBError, err)
	}

	return &task, nil
}

func UpdateTask(tx *db.Session, id uint, attrs models.Attrs) (org *models.Task, re e.Error) {
	org = &models.Task{}
	if _, err := models.UpdateAttr(tx.Where("id = ?", id), &models.Task{}, attrs); err != nil {
		return nil, e.New(e.DBError, fmt.Errorf("update task error: %v", err))
	}
	if err := tx.Where("id = ?", id).First(org); err != nil {
		return nil, e.New(e.DBError, fmt.Errorf("query task error: %v", err))
	}
	return
}

func GetTaskById(tx *db.Session, id uint) (*models.Task, e.Error) {
	o := models.Task{}
	if err := tx.Where("id = ?", id).First(&o); err != nil {
		if e.IsRecordNotFound(err) {
			return nil, e.New(e.TaskNotExists, err)
		}
		return nil, e.New(e.DBError, err)
	}
	return &o, nil
}

func GetTaskByGuid(tx *db.Session, guid string) (*models.Task, e.Error) {
	o := models.Task{}
	if err := tx.Where("guid = ?", guid).First(&o); err != nil {
		if e.IsRecordNotFound(err) {
			return nil, e.New(e.TaskNotExists, err)
		}
		return nil, e.New(e.DBError, err)
	}
	return &o, nil
}

func QueryTask(tx *db.Session, status, q string, tplId uint) *db.Session {
	query := tx.Table(fmt.Sprintf("%s as task", models.Task{}.TableName())).
		Where("template_id = ?", tplId).
		Joins(fmt.Sprintf("left join %s as tpl on tpl.id = task.template_id", models.Template{}.TableName())).
		LazySelectAppend("task.*, tpl.repo_branch")
	if status != "" {
		query = query.Where("task.status = ?", status)
	}
	if q != "" {
		qs := "%" + q + "%"
		query = query.Where("task.name LIKE ? OR task.description LIKE ?", qs, qs)
	}

	return query.Order("task.created_at DESC")
}

func TaskDetail(tx *db.Session, taskId uint) *db.Session {
	return tx.Table(models.Task{}.TableName()).Select(fmt.Sprintf("%s.*, tpl.*", models.Task{}.TableName())).
		Joins(fmt.Sprintf("left join %s as tpl on tpl.id = %s.template_id", models.Template{}.TableName(), models.Task{}.TableName())).
		Where(fmt.Sprintf("%s.id = %d", models.Task{}.TableName(), taskId))
}

func LastTask(tx *db.Session, tplId uint) *db.Session {
	return tx.Table(models.Task{}.TableName()).Where("template_id = ?", tplId)
}

type LastTaskInfo struct {
	Status    string    `json:"status"`
	Guid      string    `json:"taskGuid"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetTaskByTplId(tx *db.Session, tplId uint) (*LastTaskInfo, e.Error) {
	lastTaskInfo := LastTaskInfo{}
	err := tx.Table(models.Task{}.TableName()).
		Select("status, guid, updated_at").
		Where("template_id = ?", tplId).
		//Where("status in (?)",statusList).
		Find(&lastTaskInfo)
	if err != nil {
		return nil, e.New(e.DBError, err)
	}
	return &lastTaskInfo, nil
}

//var (
//	taskTicker *time.Ticker = time.NewTicker(time.Duration(configs.Get().Task.TimeTicker) * time.Second)
//	runnerAddr string       = configs.Get().Task.Addr
//	//runnerAddr string       = ""
//)

func runningTaskEnvParam(tpl *models.Template, runnerId string, task *models.Task) interface{} {
	tplVars := make([]forms.Var, 0)
	taskVars := make([]forms.VarOpen, 0)
	param := make(map[string]interface{})

	tplVarsByte, _ := tpl.Vars.MarshalJSON()
	taskVarsByte, _ := task.SourceVars.MarshalJSON()

	if !tpl.Vars.IsNull() {
		_ = json.Unmarshal(tplVarsByte, &tplVars)
	}

	if !task.SourceVars.IsNull() {
		_ = json.Unmarshal(taskVarsByte, &taskVars)
	}

	tplVars = append(tplVars, resourceEnvParam(runnerId, tpl.OrgId)...)
	for _, v := range varsDuplicateRemoval(taskVars, tplVars) {
		if v.Key == "" {
			continue
		}
		if v.Type == consts.Terraform && !strings.HasPrefix(v.Key, consts.TerraformVar) {
			v.Key = fmt.Sprintf("%s%s", consts.TerraformVar, v.Key)
		}
		if v.IsSecret != nil && *v.IsSecret {
			param[v.Key] = utils.AesDecrypt(v.Value)
		} else {
			param[v.Key] = v.Value
		}
	}
	return param
}

func varsDuplicateRemoval(taskVars []forms.VarOpen, tplVars []forms.Var) []forms.Var {
	if taskVars == nil || len(taskVars) == 0 {
		return tplVars
	}
	vars := make([]forms.Var, 0)
	//taskV := make(map[string]forms.VarOpen, 0)
	tplV := make(map[string]forms.Var, 0)
	for _, tplv := range tplVars {
		tplV[tplv.Key] = tplv
	}
	isSecret := false
	for _, taskv := range taskVars {
		if taskv.Name == "" {
			continue
		}
		if taskv.Value == "" {
			vars = append(vars, tplV[taskv.Name])
		} else {
			vars = append(vars, forms.Var{
				Key:      taskv.Name,
				Value:    taskv.Value,
				IsSecret: &isSecret,
			})
		}
	}
	return vars
}

func resourceEnvParam(runnerId string, orgId uint) []forms.Var {
	vars := make([]forms.Var, 0)
	ra := []models.ResourceAccount{}
	//org,_:=getorg
	if err := db.Get().Debug().Joins(fmt.Sprintf("left join %s as crm on %s.id = crm.resource_account_id",
		models.CtResourceMap{}.TableName(), models.ResourceAccount{}.TableName())).
		Where("crm.ct_service_id = ?", runnerId).
		Where(fmt.Sprintf("%s.status = '%s'", models.ResourceAccount{}.TableName(), consts.ResourceAccountEnable)).
		Where(fmt.Sprintf("%s.org_id = %d", models.ResourceAccount{}.TableName(), orgId)).
		Find(&ra); err != nil {
		logs.Get().Errorf("ResourceAccount db err %v: ", err)
		return nil
	}

	for _, raInfo := range ra {
		varsByte, _ := raInfo.Params.MarshalJSON()
		if !raInfo.Params.IsNull() {
			v := make([]forms.Var, 0)
			_ = json.Unmarshal(varsByte, &v)
			vars = append(vars, v...)
		}
	}

	return vars
}

func getBackendInfo(backendInfo models.JSON, containerId string) []byte {
	attr := models.Attrs{}
	_ = json.Unmarshal(backendInfo, &attr)
	attr["container_id"] = containerId
	b, _ := json.Marshal(attr)
	return b
}

func GetTFLog(logPath string) map[string]interface{} {
	loggers := logs.Get()
	path := fmt.Sprintf("%s/%s", logPath, consts.TaskLogName)
	f, err := os.Open(path)
	if err != nil {
		loggers.Error(err)
		return nil
	}
	defer f.Close()
	result := map[string]interface{}{}
	rd := bufio.NewReader(f)
	for {
		str, _, err := rd.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				loggers.Error("Read Error:", err.Error())
				break
			}
		}
		LogStr := string(str)
		if strings.Contains(LogStr, "No changes. Infrastructure is up-to-date.") {
			result["add"] = "0"
			result["change"] = "0"
			result["destroy"] = "0"
			result["allowApply"] = false
			break
		} else if strings.Contains(LogStr, `Plan:`) {
			r, _ := regexp.Compile(`([\d]+) to add, ([\d]+) to change, ([\d]+) to destroy`)
			params := r.FindStringSubmatch(LogStr)
			if len(params) == 4 {
				result["add"] = params[1]
				result["change"] = params[2]
				result["destroy"] = params[3]
				result["allowApply"] = true
			}
			break
		} else if strings.Contains(LogStr, `Apply complete!`) {
			r, _ := regexp.Compile(`Apply complete! Resources: ([\d]+) added, ([\d]+) changed, ([\d]+) destroyed.`)
			params := r.FindStringSubmatch(LogStr)
			if len(params) == 4 {
				result["add"] = params[1]
				result["change"] = params[2]
				result["destroy"] = params[3]
				result["allowApply"] = false
			}
			break
		}
	}
	return result
}
