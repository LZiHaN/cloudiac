// Copyright (c) 2015-2022 CloudJ Technology Co., Ltd.

package runner

import "time"

/*
provider plugin 的查找逻辑:
1. 检查 cache 目录是否存在目标 provider，存在则直接使用，否则
2. 检查本地 plugins 目录(包含多个目录，具体参考下方文档)下是否存在，存在则拷贝一份(或创建软链接)到 cache 目录，否则
3. 从网络下载文件，并保存到 cache 目录

最后将 cache 目录的文件链接到当前目录的 .terraform/providers

参考文档:
- https://www.terraform.io/docs/cli/config/config-file.html#implied-local-mirror-directories
- https://www.terraform.io/docs/cli/config/config-file.html#provider-plugin-cache
*/

/////
// 以下定义的是 runner 启动任务后容器内部的路径，直接以常量配置即可
const (
	ContainerWorkspace = "/cloudiac/workspace"
	// CodeDir 必须为 ContainerWorkspace 的子目录
	ContainerCodeDir = "/cloudiac/workspace/code"

	ContainerCertificateDir  = "/cloudiac/cert"                    // 挂载consul证书资源
	ContainerAssetsDir       = "/cloudiac/assets"                  // 挂载依赖资源，如 terraform.py 等(己打包到 worker 镜像)
	ContainerPluginPath      = "/cloudiac/terraform/plugins"       // 预置 providers 目录(己打包到镜像)
	ContainerPluginCachePath = "/cloudiac/terraform/plugins-cache" // terraform plugins 缓存目录
)

const (
	TaskScriptName = "run.sh"
	TaskLogName    = "output.log"

	TaskStepInfoFileName      = "step-info.json"
	TaskContainerInfoFileName = "container.json"
	TaskControlFileName       = "control.json"

	TerraformrcFileName = "terraformrc"
	EnvironmentFile     = "environment"

	CloudIacTfFile   = "_cloudiac.tf"
	CloudIacPlayVars = "_cloudiac_play_vars.yml"
	CloudIacTfvarsJson = "_cloudiac.tfvars.json"

	CloudIacAnsibleRequirements = "requirements.yml"

	TFStateJsonFile  = "tfstate.json"
	TFPlanJsonFile   = "tfplan.json"
	TFProviderSchema = "tfproviderschema.json"

	AnsibleStateAnalysisName = "terraform.py"

	FollowLogDelay = time.Second // follow 文件时读到 EOF 后进行下次读取的等待时长

	PoliciesDir      = "policies"
	ScanInputMapFile = "tfmap.json"
	ScanInputFile    = "tfscan.json"
	ScanResultFile   = "scan_result.json"
	ScanLogFile      = "scan.log"
	RegoResultFile   = "scan_raw.json"

	PopulateSourceLineCount = 3
)
