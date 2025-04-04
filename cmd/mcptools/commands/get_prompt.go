package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// GetPromptCmd creates the get-prompt command.
func GetPromptCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "get-prompt prompt [command args...]",
		Short:              "Get a prompt on the MCP server",
		DisableFlagParsing: true,
		SilenceUsage:       true,
		Run: func(thisCmd *cobra.Command, args []string) {
			if len(args) == 1 && (args[0] == FlagHelp || args[0] == FlagHelpShort) {
				_ = thisCmd.Help()
				return
			}

			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Error: prompt name is required")
				fmt.Fprintln(
					os.Stderr,
					"Example: mcp get-prompt read_file npx -y @modelcontextprotocol/server-filesystem ~",
				)
				os.Exit(1)
			}

			cmdArgs := args
			parsedArgs := []string{}
			promptName := ""

			i := 0
			promptExtracted := false

			for i < len(cmdArgs) {
				switch {
				case (cmdArgs[i] == FlagFormat || cmdArgs[i] == FlagFormatShort) && i+1 < len(cmdArgs):
					FormatOption = cmdArgs[i+1]
					i += 2
				case (cmdArgs[i] == FlagParams || cmdArgs[i] == FlagParamsShort) && i+1 < len(cmdArgs):
					ParamsString = cmdArgs[i+1]
					i += 2
				case !promptExtracted:
					promptName = cmdArgs[i]
					promptExtracted = true
					i++
				default:
					parsedArgs = append(parsedArgs, cmdArgs[i])
					i++
				}
			}

			if promptName == "" {
				fmt.Fprintln(os.Stderr, "Error: prompt name is required")
				fmt.Fprintln(
					os.Stderr,
					"Example: mcp get-prompt read_file npx -y @modelcontextprotocol/server-filesystem ~",
				)
				os.Exit(1)
			}

			var params map[string]any
			if ParamsString != "" {
				if jsonErr := json.Unmarshal([]byte(ParamsString), &params); jsonErr != nil {
					fmt.Fprintf(os.Stderr, "Error: invalid JSON for params: %v\n", jsonErr)
					os.Exit(1)
				}
			}

			mcpClient, clientErr := CreateClientFunc(parsedArgs)
			if clientErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", clientErr)
				os.Exit(1)
			}

			resp, execErr := mcpClient.GetPrompt(promptName)
			if formatErr := FormatAndPrintResponse(thisCmd, resp, execErr); formatErr != nil {
				fmt.Fprintf(os.Stderr, "%v\n", formatErr)
				os.Exit(1)
			}
		},
	}
}
