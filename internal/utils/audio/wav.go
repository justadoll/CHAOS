package audio

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/justadoll/CHAOS/internal/utils/constants"
)

// WriteWav write a given content to a wav file
func WriteWav(content []byte) (string, error) {
	filename := fmt.Sprint(uuid.New().String(), ".wav")
	file, err := os.OpenFile(
		fmt.Sprint(constants.TempDirectory, filename),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		os.ModePerm,
	)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return "", err
	}
	return filename, nil
}
