package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/query-builder-generator/src/compiler"
	"github.com/query-builder-generator/src/dom/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qbc",
	Short: "Query Builder Compiler",
	Long: `Query Builder Compiler - compiles query builder files into source code.`,
	Run: func(cmd *cobra.Command, args []string) {
	    fmt.Println("Root command. Use subcommands")
        //generateFile(cmd, args)
    },
}

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate Java classes",
    Long: "Generate Java classes - to be used in portal code",
    Run: func (cmd *cobra.Command, args []string) {
        err := generateFile(cmd)
        if err != nil {
        	panic(err.Error())
		}
    },
}


func generateFile(cmd *cobra.Command) error {
    fmt.Println("Reading file at path [" + inputFilePath +"]")
    data, err := ioutil.ReadFile(inputFilePath)
    if err != nil {
        return err
    }

    var document = parser.Parse(string(data))
    document.Package = strings.Replace(path.Dir(classPath), "/", ".", -1)

    var compiler = compiler.Compiler{}
    var outputContent = compiler.Generate(&document)

    var absoluteFilePath = path.Join(rootPath, classPath)

    fmt.Println("Writing file at path [" + absoluteFilePath +"]")
    err = ioutil.WriteFile(absoluteFilePath, []byte(outputContent), 0777)
    if err != nil {
		return err
	}

    return nil
}

func addCommands() {
    rootCmd.AddCommand(generateCmd)
}

func Execute() {
    addCommands()
	cobra.CheckErr(rootCmd.Execute())
}


// adds flags
var inputFilePath, rootPath, classPath string
func init() {
	cobra.OnInitialize(initConfig)

	generateCmd.Flags().StringVar(&inputFilePath, "input", "default", "--input=<Input File path>")
    generateCmd.Flags().StringVar(&rootPath, "root", "default", "--root=<Root directory path>")
	generateCmd.Flags().StringVar(&classPath, "class", "default", "--class=<Class file path>")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qbc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".qbc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".qbg")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
