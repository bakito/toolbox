package makefile

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/bakito/toolbox/pkg/types/renovate"
)

const renovateJson = "renovate.json"

func updateRenovateConf() error {
	withRenovate, cfg, err := updateRenovateConfInternal(renovateJson)
	if err != nil || !withRenovate {
		return err
	}
	return os.WriteFile(renovateJson, cfg, 0o600)
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

	var merged []map[string]interface{}
	if err := covert(&cms, &merged); err != nil {
		return false, nil, err
	}

	renovateCfg["customManagers"] = merged
	pp, err := prettyPrint(renovateCfg)
	return true, pp, err
}

func covert(from interface{}, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, to); err != nil {
		return err
	}
	return nil
}

func readRenovateConfig(renovateCfgFile string) (map[string]interface{}, error) {
	b, err := os.ReadFile(renovateCfgFile)
	if err != nil {
		return nil, err
	}
	renovateConfig := make(map[string]interface{})
	if err := json.Unmarshal(b, &renovateConfig); err != nil {
		return nil, err
	}
	return renovateConfig, nil
}

func prettyPrint(renovateConfig map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(renovateConfig); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
