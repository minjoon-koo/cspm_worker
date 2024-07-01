package CLI

import (
	"os/exec"
)

func SteampipeQuery(sql string) []byte {
	cmd := exec.Command("steampipe", "query", sql, "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		println("Failed to execute steampipe query")

	}
	return output
}
