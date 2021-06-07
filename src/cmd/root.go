package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"io/ioutil"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qbc",
	Short: "Query Builder Compiler",
	Long: `Query Builder Compiler - compiles query builder files into source code.`,
	Run: func(cmd *cobra.Command, args []string) {
    generateFile(cmd, args)
    },
}

func generateFile(cmd *cobra.Command, args []string) error {
    var inputPath string = args[0];
    var outputPath string = args[1];
    fmt.Println("Reading file at path [" + inputPath +"]")

    data, err := ioutil.ReadFile(inputPath)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print(string(data))

    var outputContent string = `
package io.harness.beans;

import io.harness.beans.DelegateTask.DelegateTaskKeys;
import io.harness.persistence.HPersistence;
import io.harness.query.PersistentQuery;

import org.mongodb.morphia.query.Query;

public class DelegateTasksQuery implements PersistentQuery {
  private Query<DelegateTask> query;

  public static DelegateTasksQuery create(HPersistence persistence) {
    return new DelegateTasksQuery(persistence.createQuery(DelegateTask.class)
                                      .project(DelegateTaskKeys.uuid, true)
                                      .project(DelegateTaskKeys.data_timeout, true));
  }

  private DelegateTasksQuery(Query<DelegateTask> query) {
    this.query = query;
  }

  public static class Final {
    DelegateTasksQuery self;
    Final(DelegateTasksQuery self) {
      this.self = self;
    }

    public Query<DelegateTask> query() {
      return self.query;
    }
  }

  public static class FilterStatus {
    DelegateTasksQuery self;
    FilterStatus(DelegateTasksQuery self) {
      this.self = self;
    }

    public Final status(DelegateTask.Status status) {
      self.query.filter(DelegateTaskKeys.status, status);
      return new Final(self);
    }
  }

  public static class FilterDelegateId {
    DelegateTasksQuery self;
    FilterDelegateId(DelegateTasksQuery self) {
      this.self = self;
    }

    public FilterStatus delegateId(String delegateId) {
      self.query.filter(DelegateTaskKeys.delegateId, delegateId);
      return new FilterStatus(self);
    }
  }

  public static class FilterUuids {
    DelegateTasksQuery self;
    FilterUuids(DelegateTasksQuery self) {
      this.self = self;
    }

    public FilterDelegateId uuids(Iterable<String> uuids) {
      self.query.field(DelegateTaskKeys.uuid).in(uuids);
      return new FilterDelegateId(self);
    }
  }

  public FilterUuids accountId(String accountId) {
    query.filter(DelegateTaskKeys.accountId, accountId);
    return new FilterUuids(this);
  }
}
`
    // Write file
    f, err := os.Create(outputPath)
    if err != nil {
       fmt.Println(err)
    }

    _, err = f.WriteString(outputContent)
    if err != nil {
      fmt.Println(err)
    }

    return nil
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

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
		viper.SetConfigName(".qbc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
