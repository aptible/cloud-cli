package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/cmd/asset"
	"github.com/aptible/cloud-cli/config"
)

var (
	cfgFile    string
	token      string
	authDomain string
	apiDomain  string
	org        string
	env        string
	debug      bool
)

var logo = `
      ..'...''..             .','.                          ... .',.         .,'.
    .;oxo;.'cddc'.          .:0NXo.                  .od,  .:kc.'kK:        .lXx.
  .;oxo:'....,lxxl,.        ,OXk0Kc.   ......'''..  .:KNo....,. 'OXc.'''..  .oNk.   ..''...
 'oxo;.,lc.,l:',cxx:.      .xNx'cX0,   'x0xxxxkOOd,.:ONW0d,.l0l.'0W0xxxkOx,..oNk. .cxkxxkOo'.
 ,l;.,lxOl.;xkd:.'cc.     .lXO, .dNx.  ,0W0:....lKK:.:KNo...dWx.'0WO;..'oX0;.oNk..xXx,..'oKO,
 ..,lxkkOl.;xOkxd:...    .:KWOlclxXNo. ,0No.    .xWd.,0Xc. .dWx.'0Nl.   .kNo.oNk.:KNOddddx00c.
 .lxxc;oOl.;xkc;lxd;.    ,ON0xxxxxOXXc.,0Wx.    'ONo.,0Xc. .dWx.'OWd.   ,0Nl.oNk.,0Xo'''';lc'
 ,oc'..lOl.;xk:..,ll.   .xNO'     .dN0,,0WXxlccoOKx' 'kNOl'.dNx.'OWKxlco0Xd..oNk..:OOdcclk0o.
 ...  .,:'..;;.   ...   .cl,.      .cl,;0Xdcodxdl,.  .'col'.,l,..:l::odol,. .,l;. ..;lddoc,.
                                       ,OK:.
                                       .,;.
`

var desc = `aptible is a command line interface to the Aptible.com platform.

It allows users to manage authentication, application launch,
deployment, logging, and more with just the one command.

* Provision an app with the app create command
* Provision a datastore with the datastore create command
* View a deployed web application with the open command
* View detailed information about an app or datastore with the info command

To read more, use the docs command to view Aptible's help on the web.`

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "aptible",
		Short: "aptible is a command line interface to the Aptible.com platform.",
		Long:  fmt.Sprintf("%s\n%s", logo, desc),
	}

	cobra.OnInitialize(initConfig())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "common", "", "common file (default is $HOME/.aptible.yaml)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "jwt token")
	rootCmd.PersistentFlags().StringVar(&authDomain, "auth-domain", "auth-api-master.aptible-staging.com", "auth domain")
	rootCmd.PersistentFlags().StringVar(&apiDomain, "api-domain", "cloud-api.aptible.com", "api domain")
	rootCmd.PersistentFlags().StringVar(&org, "org", "", "organization id")
	rootCmd.PersistentFlags().StringVar(&env, "env", "", "environment id")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug logging")

	errs := []error{
		viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")),
		viper.BindPFlag("auth-domain", rootCmd.PersistentFlags().Lookup("auth-domain")),
		viper.BindPFlag("api-domain", rootCmd.PersistentFlags().Lookup("api-domain")),
		viper.BindPFlag("org", rootCmd.PersistentFlags().Lookup("org")),
		viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env")),
		viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")),
	}

	viperErrOnInit := false
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			viperErrOnInit = true
		}
	}
	if viperErrOnInit {
		log.Println("Unable to initialize viper config")
		os.Exit(1)
	}

	assetCmd := asset.NewAssetCmd()
	envCmd := NewEnvCmd()
	dsCmd := asset.NewDatastoreCmd()
	orgCmd := NewOrgCmd()
	configCmd := config.NewConfigCmd()
	vpcCmd := asset.NewVPCCmd()

	rootCmd.AddCommand(
		assetCmd,
		configCmd,
		dsCmd,
		envCmd,
		orgCmd,
		vpcCmd,
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(root *cobra.Command) {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() func() {
	return func() {
		vconfig := viper.GetViper()
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		if cfgFile != "" {
			// Use common file from the flag.
			vconfig.SetConfigFile(cfgFile)
		} else {
			vconfig.AddConfigPath(home)
			vconfig.SetConfigName(".aptible")
			vconfig.SetConfigType("yaml")
		}

		vconfig.AutomaticEnv()

		if token == "" {
			token, err = config.FindToken(home, fmt.Sprintf("https://%s", authDomain))
			if err != nil {
				fmt.Println("Unable to load token")
				os.Exit(1)
			}
		}
		vconfig.Set("token", token)

		if err := vconfig.ReadInConfig(); err == nil {
			fmt.Println("Using common file:", vconfig.ConfigFileUsed())
		}
	}
}
