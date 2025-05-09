package makefile

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/bakito/toolbox/pkg/types/renovate"
)

const renovateJSON = "renovate.json"

func updateRenovateConf() error {
	withRenovate, cfg, err := updateRenovateConfInternal(renovateJSON)
	if err != nil || !withRenovate {
		return err
	}
	return os.WriteFile(renovateJSON, cfg, 0o600)
}

func updateRenovateConfInternal(renovateCfgFile string) (bool, []byte, error) {
	if _, err := os.Stat(renovateCfgFile); errors.Is(err, os.ErrNotExist) {
		// no renovate config found, abort
		return false, nil, nil
	}

	renovateCfg, err := readRenovateConfig(renovateCfgFile)
	if err != nil {
		return false, nil, err
	}

	cms := renovate.CustomManagers{}
	if cm, ok := renovateCfg["customManagers"]; ok {
		if err := covert(&cm, &cms); err != nil {
			return false, nil, err
		}

		found := false
		for i, manager := range cms {
			if manager.IsToolbox() {
				found = true
				manager.UpdateParams()
				cms[i] = manager
			} else if len(manager.FileMatch) > 0 && len(manager.ManagerFilePatterns) == 0 {
				manager.ManagerFilePatterns = manager.FileMatch
				manager.FileMatch = nil
				cms[i] = manager
			}
		}
		if !found {
			// add toolbox config
			cms = append(cms, renovate.Config())
		}
	} else {
		// add default config
		cms = []renovate.CustomManager{renovate.Config()}
	}

	var merged []map[string]any
	if err := covert(&cms, &merged); err != nil {
		return false, nil, err
	}

	renovateCfg["customManagers"] = merged
	pp, err := prettyPrint(renovateCfg)
	return true, pp, err
}

func covert(from, to any) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}

func readRenovateConfig(renovateCfgFile string) (map[string]any, error) {
	b, err := os.ReadFile(renovateCfgFile)
	if err != nil {
		return nil, err
	}
	renovateConfig := make(map[string]any)
	if err := json.Unmarshal(b, &renovateConfig); err != nil {
		return nil, err
	}
	return renovateConfig, nil
}

func prettyPrint(renovateConfig map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(renovateConfig); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
