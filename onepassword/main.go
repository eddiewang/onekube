package onepassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func NewOpClient() *Client {
	return &Client{}
}

func (c *Client) runOp(opCmd string, args []string) ([]byte, error) {
	cmdArgs := []string{opCmd}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, "--format", "json")

	cmd := exec.Command("op", cmdArgs...)
	errBuf := bytes.NewBuffer(nil)
	cmd.Stderr = errBuf

	out, err := cmd.Output()
	if err != nil {
		if errBuf.String() != "" {
			return nil, fmt.Errorf("op returned err: %s", errBuf.String())
		}
		return nil, err
	}

	return out, nil
}

func (c *Client) runOpAndUnmarshal(opCmd string, args []string, unmarshalInto any) error {
	out, err := c.runOp(opCmd, args)
	if err != nil {
		return err
	}

	err = json.Unmarshal(out, unmarshalInto)
	if err != nil {
		return err
	}

	return nil
}
