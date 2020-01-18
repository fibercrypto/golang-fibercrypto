package local

import (
	"encoding/json"
	"strconv"

	"github.com/fibercrypto/fibercryptowallet/src/errors"
	"github.com/fibercrypto/fibercryptowallet/src/params"
	qtcore "github.com/therecipe/qt/core"
)

const (
	LocalWallet = iota
	RemoteWallet
)

var (
	confManager         *ConfigManager
	OptionNotFoundError = errors.ErrInvalidOptions
)

func init() {
	qs := qtcore.NewQSettings(qtcore.QCoreApplication_OrganizationName(), qtcore.QCoreApplication_ApplicationName(), nil)
	confManager = &ConfigManager{
		setting:  qs,
		sections: make(map[string]*SectionManager),
	}

	valueLifeTime := strconv.FormatUint(params.DataRefreshTimeout, 10)

	cache := map[string]string{"lifeTime": valueLifeTime}

	cacheBytes, err := json.Marshal(cache)
	if err != nil {
		return
	}
	cacheOpt := NewOption("cache", []string{}, false, string(cacheBytes))

	_ = confManager.RegisterSection("global", []*Option{cacheOpt})
}

type ConfigManager struct {
	setting  *qtcore.QSettings
	sections map[string]*SectionManager
}

func (cm *ConfigManager) GetSections() []string {
	return cm.setting.ChildGroups()
}

func (cm *ConfigManager) GetSectionManager(section string) *SectionManager {
	sectionM, ok := cm.sections[section]
	if !ok {
		return nil
	}
	return sectionM
}

func (cm *ConfigManager) RegisterSection(name string, options []*Option) *SectionManager {

	cm.sections[name] = &SectionManager{
		name:     name,
		settings: cm.setting,
		options:  options,
	}

	cm.setting.BeginGroup(name)
	defer cm.setting.EndGroup()
	defer cm.setting.Sync()

	for _, opt := range options {
		depthLevel := 0
		for _, sect := range opt.sectionPath {
			cm.setting.BeginGroup(sect)
			depthLevel++
		}
		if !opt.optional && !cm.setting.Contains(opt.name) {
			cm.setting.SetValue(opt.name, qtcore.NewQVariant1(opt._default))

		}
		for i := 0; i < depthLevel; i++ {
			cm.setting.EndGroup()
		}
	}

	return cm.sections[name]
}

type SectionManager struct {
	name     string
	settings *qtcore.QSettings
	options  []*Option
}

func (sm *SectionManager) GetValue(name string, sectionPath []string) (string, error) {
	sm.settings.BeginGroup(sm.name)
	defer sm.settings.EndGroup()
	for _, sect := range sectionPath {
		groups := sm.settings.ChildGroups()
		finded := false
		for _, grp := range groups {
			if grp == sect {
				finded = true
				break
			}
		}
		if !finded {
			return "", OptionNotFoundError
		}
		sm.settings.BeginGroup(sect)
		defer sm.settings.EndGroup()
	}
	val := sm.settings.Value(name, qtcore.NewQVariant())
	if val.IsNull() {
		return "", OptionNotFoundError
	}
	return val.ToString(), nil
}

func (sm *SectionManager) GetDefaultValue(option string, sectionPath []string, name string) (string, error) {
	for _, opt := range sm.options {
		if compareStringSlices(sectionPath, opt.sectionPath) && option == opt.name {
			store := make(map[string]string)
			err := json.Unmarshal([]byte(opt._default), &store)
			if err != nil {
				return "", err
			}
			val, ok := store[name]
			if !ok {
				return "", errors.ErrInvalidOptions
			}
			return val, nil
		}
	}
	return "", OptionNotFoundError
}

func (sm *SectionManager) Save(name string, sectionPath []string, value string) error {
	sm.settings.BeginGroup(sm.name)
	defer sm.settings.EndGroup()
	for _, sect := range sectionPath {
		groups := sm.settings.ChildGroups()
		finded := false
		for _, grp := range groups {
			if grp == sect {
				finded = true
				break
			}
		}
		if !finded {
			return OptionNotFoundError
		}
		sm.settings.BeginGroup(sect)
		defer sm.settings.EndGroup()
	}
	if !sm.settings.Contains(name) {
		return OptionNotFoundError
	}
	sm.settings.SetValue(name, qtcore.NewQVariant1(value))
	return nil
}

func (sm *SectionManager) GetValues(sectionPath []string) ([]string, error) {
	sm.settings.BeginGroup(sm.name)
	defer sm.settings.EndGroup()

	for _, sect := range sectionPath {

		groups := sm.settings.ChildGroups()
		finded := false
		for _, grp := range groups {
			if grp == sect {
				finded = true
				break
			}
		}
		if !finded {
			return nil, OptionNotFoundError
		}
		sm.settings.BeginGroup(sect)
		defer sm.settings.EndGroup()
	}
	keys := sm.settings.ChildKeys()
	values := make([]string, 0)
	for _, key := range keys {
		values = append(values, sm.settings.Value(key, qtcore.NewQVariant()).ToString())
	}
	return values, nil
}

func (sm *SectionManager) GetPaths() [][]string {
	sm.settings.BeginGroup(sm.name)
	defer sm.settings.EndGroup()
	return sm.getPaths([]string{})
}

func (sm *SectionManager) getPaths(prefix []string) [][]string {
	if len(prefix) > 0 {
		lastGrp := prefix[len(prefix)-1]
		sm.settings.BeginGroup(lastGrp)
		defer sm.settings.EndGroup()
	}

	grps := sm.settings.ChildGroups()
	if len(grps) == 0 {
		return [][]string{prefix}
	}
	values := make([][]string, 0)
	for _, grp := range grps {
		values = append(values, sm.getPaths(append(prefix, grp))...)
	}
	return values

}

type Option struct {
	name        string
	sectionPath []string
	optional    bool
	_default    string
}

func NewOption(name string, sectionPath []string, optional bool, _default string) *Option {
	if !optional && _default == "" {
		return nil
	}
	return &Option{
		name:        name,
		sectionPath: sectionPath,
		optional:    optional,
		_default:    _default,
	}
}

func GetConfigManager() *ConfigManager {
	return confManager
}

func compareStringSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, val := range a {
		if val != b[i] {
			return false
		}
	}
	return true
}
