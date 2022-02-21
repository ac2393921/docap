package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/jesseduffield/yaml"
)

type UserConfig struct {
	Gui              GuiConfig
	ConfirmOnQuit    bool
	CommandTemplates CommandTemplatesConfig
	CustomCommands   CustomCommands
	BulkCommands     CustomCommands
	OS               OSConfig
	Update           UpdateConfig
	Stats            StatsConfig
}

type CommandTemplatesConfig struct {
	RestartService           string
	DockerCompose            string
	StopService              string
	ServiceLogs              string
	ViewServiceLogs          string
	RebuildService           string
	RecreateService          string
	ViewContainerLogs        string
	ContainerLogs            string
	AllLogs                  string
	ViewAllLogs              string
	DockerComposeConfig      string
	CheckDockerComposeConfig string
	ServiceTop               string
}

type ThemeConfig struct {
	ActiveBorderColor   []string
	InactiveBorderColor []string
	OptionsTextColor    []string
}

type GuiConfig struct {
	ScrollHeight         int
	Language             string
	ScrollPastBottom     bool
	IgnoreMouseEvents    bool
	Theme                ThemeConfig
	ShowAllContainers    bool
	ReturnImmediately    bool
	WrapMainPanel        bool
	LegacySortContainers bool
}

type AppConfig struct {
	Debug       bool
	Version     string
	Commit      string
	BuildDate   string
	Name        string
	BuildSource string
	UserConfig  *UserConfig
	ConfigDir   string
	ProjectDir  string
}

type OSConfig struct {
	OpenCommand     string
	OpenLickCommand string
}

type UpdateConfig struct {
	DockerRefreshInterval time.Duration
}

type StatsConfig struct {
	Graphs      []GraphConfig
	MaxDuration time.Duration
}

type GraphConfig struct {
	Min      float64
	Max      float64
	Height   int
	Caption  string
	StatPath string
	Color    string
	MinType  string
	MaxType  string
}

type CustomCommand struct {
	Name             string
	Attach           bool
	Command          string
	ServiceNames     []string
	InternalFunction func() error
}

type CustomCommands struct {
	Containers []CustomCommand
	Services   []CustomCommand
	Images     []CustomCommand
	Volumes    []CustomCommand
}

func GetDefaultConfig() UserConfig {
	duration, err := time.ParseDuration("3m")
	if err != nil {
		panic(err)
	}

	return UserConfig{
		Gui: GuiConfig{
			ScrollHeight:      2,
			Language:          "auto",
			ScrollPastBottom:  false,
			IgnoreMouseEvents: false,
			Theme: ThemeConfig{
				ActiveBorderColor:   []string{"green", "bold"},
				InactiveBorderColor: []string{"default"},
				OptionsTextColor:    []string{"blue"},
			},
			ShowAllContainers:    false,
			ReturnImmediately:    false,
			WrapMainPanel:        false,
			LegacySortContainers: false,
		},
		ConfirmOnQuit: false,
		CommandTemplates: CommandTemplatesConfig{
			DockerCompose:            "docker-compose",
			RestartService:           "{{ .DockerCompose }} restart {{ .Service.Name }}",
			RebuildService:           "{{ .DockerCompose }} up -d --build {{ .Service.Name }}",
			RecreateService:          "{{ .DockerCompose }} up -d --force-recreate {{ .Service.Name }}",
			StopService:              "{{ .DockerCompose }} stop {{ .Service.Name }}",
			ServiceLogs:              "{{ .DockerCompose }} logs --since=60m --follow {{ .Service.Name }}",
			ViewServiceLogs:          "{{ .DockerCompose }} logs --follow {{ .Service.Name }}",
			AllLogs:                  "{{ .DockerCompose }} logs --tail=300 --follow",
			ViewAllLogs:              "{{ .DockerCompose }} logs",
			DockerComposeConfig:      "{{ .DockerCompose }} config",
			CheckDockerComposeConfig: "{{ .DockerCompose }} config --quiet",
			ContainerLogs:            "docker logs --timestamps --follow --since=60m {{ .Container.ID }}",
			ViewContainerLogs:        "docker logs --timestamps --follow --since=60m {{ .Container.ID }}",
			ServiceTop:               "{{ .DockerCompose }} top {{ .Service.Name }}",
		},
		CustomCommands: CustomCommands{
			Containers: []CustomCommand{
				{
					Name:    "bash",
					Command: "docker exec -it {{ .Container.ID }} /bin/sh -c 'eval $(grep ^$(id -un): /etc/passwd | cut -d : -f 7-)'",
					Attach:  true,
				},
			},
			Services: []CustomCommand{},
			Images:   []CustomCommand{},
			Volumes:  []CustomCommand{},
		},
		BulkCommands: CustomCommands{
			Services: []CustomCommand{
				{
					Name:    "up",
					Command: "{{ .DockerCompose }} up -d",
				},
				{
					Name:    "up (attached)",
					Command: "{{ .DockerCompose }} up",
					Attach:  true,
				},
				{
					Name:    "stop",
					Command: "{{ .DockerCompose }} stop",
				},
				{
					Name:    "pull",
					Command: "{{ .DockerCompose }} pull",
					Attach:  true,
				},
				{
					Name:    "build",
					Command: "{{ .DockerCompose }} build --parallel --force-rm",
					Attach:  true,
				},
				{
					Name:    "down",
					Command: "{{ .DockerCompose }} down",
				},
				{
					Name:    "down with volumes",
					Command: "{{ .DockerCompose }} down --volumes",
				},
				{
					Name:    "down with images",
					Command: "{{ .DockerCompose }} down --rmi all",
				},
				{
					Name:    "down with volumes and images",
					Command: "{{ .DockerCompose }} down --volumes --rmi all",
				},
			},
			Containers: []CustomCommand{},
			Images:     []CustomCommand{},
			Volumes:    []CustomCommand{},
		},
		OS: GetPlatformDefaultConfig(),
		Update: UpdateConfig{
			DockerRefreshInterval: time.Millisecond * 100,
		},
		Stats: StatsConfig{
			MaxDuration: duration,
			Graphs: []GraphConfig{
				{
					Caption:  "CPU (%)",
					StatPath: "DerivedStats.CPUPercentage",
					Color:    "cyan",
				},
				{
					Caption:  "Memory (%)",
					StatPath: "DerivedStats.MemoryPercentage",
					Color:    "green",
				},
			},
		},
	}
}

func configDirForVendor(vendor string, projectName string) string {
	envConfigDir := os.Getenv("CONFIG_DIR")
	if envConfigDir != "" {
		return envConfigDir
	}
	configDirs := xdg.New(vendor, projectName)

	return configDirs.ConfigHome()
}

func configDir(projectName string) string {
	legacyConfigDirectory := configDirForVendor("ac2393921", projectName)
	if _, err := os.Stat(legacyConfigDirectory); !os.IsNotExist(err) {
		return legacyConfigDirectory
	}
	configDirectory := configDirForVendor("", projectName)

	return configDirectory
}

func findOrCreateConfigDir(projectName string) (string, error) {
	folder := configDir(projectName)

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return "", os.ErrClosed
	}

	return folder, nil
}

func (c *AppConfig) ConfigFilename() string {
	return filepath.Join(c.ConfigDir, "config.yml")
}

func (c *AppConfig) WriteToUserConfig(updateConfig func(*UserConfig) error) error {
	userConfig, err := loadUserConfig(c.ConfigDir, &UserConfig{})
	if err != nil {
		return err
	}

	if err := updateConfig(userConfig); err != nil {
		return err
	}

	file, err := os.OpenFile(c.ConfigFilename(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(file).Encode(userConfig)
}

func loadUserConfig(configDir string, base *UserConfig) (*UserConfig, error) {
	fileName := filepath.Join(configDir, "config.yml")

	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(fileName)
			if err != nil {
				return nil, err
			}
			file.Close()
		} else {
			return nil, err
		}
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, base); err != nil {
		return nil, err
	}

	return base, nil
}

func loadUserConfigWithDefaults(configDir string) (*UserConfig, error) {
	config := GetDefaultConfig()

	return loadUserConfig(configDir, &config)
}

func NewAppConfig(name, version, commit, date, buildSource, projectDir string, debuggingFlag bool, composeFiles []string) (*AppConfig, error) {
	configDir, err := findOrCreateConfigDir(name)
	if err != nil {
		return nil, err
	}

	userConfig, err := loadUserConfigWithDefaults(configDir)
	if err != nil {
		return nil, err
	}

	if len(composeFiles) > 0 {
		userConfig.CommandTemplates.DockerCompose += " -f " + strings.Join(composeFiles, " -f ")
	}

	appConfig := &AppConfig{
		Name:        name,
		Version:     version,
		Commit:      commit,
		BuildDate:   date,
		Debug:       debuggingFlag || os.Getenv("DEBUG") == "TRUE",
		BuildSource: buildSource,
		UserConfig:  userConfig,
		ConfigDir:   configDir,
		ProjectDir:  projectDir,
	}

	return appConfig, nil
}
