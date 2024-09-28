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
	cfg, err := updateRenovateConfInternal(renovateJson)
	if err != nil {
		return err
	}
	return os.WriteFile(renovateJson, cfg, 0o600)
}

func updateRenovateConfInternal(renovateCfgFile string) ([]byte, error) {
	if _, err := os.Stat(renovateCfgFile); errors.Is(err, os.ErrNotExist) {
		// no renovate config found, abort
		return nil, nil
	}

	renovateCfg, err := readRenovateConfig(renovateCfgFile)
	if err != nil {
		return nil, err
	}

	cms := renovate.CustomManagers{}
	if cm, ok := renovateCfg["customManagers"]; ok {
		if err := covert(&cm, &cms); err != nil {
			return nil, err
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
		return nil, err
	}

	renovateCfg["customManagers"] = merged
	return prettyPrint(renovateCfg)
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
	if err := enc.Encode(renovateConfig); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
