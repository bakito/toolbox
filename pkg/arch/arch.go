package arch

import (
	"debug/elf"
	"debug/pe"
	"fmt"
	"os"
	"runtime"
)

func DoesBinaryMatchCurrentOSArch(fileName string) (bool, error) {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	file, err := os.Open(fileName)
	if err != nil {
		return false, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	switch currentOS {
	case "windows":
		return checkPEFileArch(file, currentArch)
	case "linux":
		return checkELFFileArch(file, currentArch)
	default:
		return false, fmt.Errorf("unsupported OS: %s", currentOS)
	}
}

func checkPEFileArch(file *os.File, currentArch string) (bool, error) {
	peFile, err := pe.NewFile(file)
	if err != nil {
		return false, fmt.Errorf("error opening PE file: %w", err)
	}

	var arch string
	switch peFile.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		arch = "386"
	case pe.IMAGE_FILE_MACHINE_AMD64:
		arch = "amd64"
	default:
		return false, fmt.Errorf("unsupported architecture in PE file: %v", peFile.Machine)
	}

	return arch == currentArch, nil
}

func checkELFFileArch(file *os.File, currentArch string) (bool, error) {
	elfFile, err := elf.NewFile(file)
	if err != nil {
		return false, fmt.Errorf("error opening ELF file: %w", err)
	}

	var arch string
	switch elfFile.Machine {
	case elf.EM_386:
		arch = "386"
	case elf.EM_X86_64:
		arch = "amd64"
	default:
		return false, fmt.Errorf("unsupported architecture in ELF file: %v", elfFile.Machine)
	}

	return arch == currentArch, nil
}
