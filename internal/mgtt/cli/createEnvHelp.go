package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

// CreateEnvHelpFileCommand is the command to create the help for env-parameter
type CreateEnvHelpFileCommand struct {
	HelpCreateEnv CreateEnvHelpFileFlag `kong:"cmd,name='help-create-env',help='Will create an ascidoc-file with an table of variables'"`
}

// CreateEnvHelpFileFlag ar parameter for the help-create-env-command
type CreateEnvHelpFileFlag struct {
	OutputFile string `kong:"help='The output where we write the help',default='docs/usage/_env.adoc'"`
}

// Run will start the server
func (s *CreateEnvHelpFileFlag) Run() error {

	programm := os.Args[0]
	os.Args = []string{programm, "--help"}

	file, err := os.OpenFile(s.OutputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	//file, err := os.OpenFile("_envvars.adoc", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed creating file")
	}
	datawriter := bufio.NewWriter(file)

	datawriter.WriteString(".Environment Parameter\n")
	datawriter.WriteString("[width=\"100%\",options=\"header\"]\n")
	datawriter.WriteString("|======================\n")

	kongInstance, _ := kong.New(&cliData)
	kongInstance.Parse([]string{})

	// flags
	for _, flag := range kongInstance.Model.Flags {
		kongFlagWriteToFile(flag, datawriter)
	}

	for _, child := range kongInstance.Model.Children {
		kongCommandWriteToFile(child, datawriter)
	}

	datawriter.WriteString("|======================\n")
	datawriter.Flush()
	file.Close()

	return nil
}

func kongFlagWriteToFile(flag *kong.Flag, writer *bufio.Writer) {
	if flag == nil {
		return
	}
	if flag.Env == "" {
		return
	}
	fileLine := fmt.Sprintf("|%s|%s|%s\n", flag.Env, flag.Help, flag.Default)
	writer.WriteString(fileLine)
	writer.Flush()
}

func kongCommandWriteToFile(node *kong.Node, writer *bufio.Writer) {
	if node == nil {
		return
	}

	// handle children
	if node.Children != nil {
		for _, child := range node.Children {
			kongCommandWriteToFile(child, writer)
		}
	}

	// handle flags
	if node.Flags != nil {
		for _, flag := range node.Flags {
			kongFlagWriteToFile(flag, writer)
		}
	}

}
