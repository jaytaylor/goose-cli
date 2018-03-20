package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/advancedlogic/GoOse"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	OutputFormat = "html"
	Verbose      bool
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", `Output format, one of "json", "html", "text", "yaml"; (default: text)`)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorExit(err)
	}
}

func errorExit(err interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "extract",
	Short: "Disambiguate an article from HTML content",
	Long:  "Runs GoOse command-line client",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(_ *cobra.Command, _ []string) {
		initLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			g       = goose.New()
			article *goose.Article
			err     error
		)
		if args[0] == "-" {
			bs, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errorExit(err)
			}
			article, err = g.ExtractFromRawHTML("", string(bs))
		} else {
			article, err = g.ExtractFromURL(args[0])
		}
		if err != nil {
			errorExit(err)
		}

		// log.Infof("Title: %s", article.Title)
		// log.Infof("Publish Date: %s", article.PublishDate)
		// h, err := article.TopNode.Html()
		// if err != nil {
		// 	errorExit(err)
		// }
		// log.Infof("HTML: ", h)
		// log.Infof("%s", )
		// log.Infof("%s", )
		// log.Infof("%s", )
		// log.Infof("%s", )

		if err := emit(article); err != nil {
			errorExit(err)
		}
	},
}

func emit(article *goose.Article) error {
	switch OutputFormat {
	case "json":
		bs, err := json.MarshalIndent(*article, "", "    ")
		if err != nil {
			return err
		}
		fmt.Print(string(bs))

	case "html":
		// fmt.Printf("%s", article.RawHTML)
		h, err := article.TopNode.Html()
		if err != nil {
			return err
		}
		fmt.Printf(h)

	case "text":
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Publish Date: %s\n", article.PublishDate)
		fmt.Printf("Canonical Link: %s\n", article.CanonicalLink)
		fmt.Printf("Final URL: %s\n", article.FinalURL)
		fmt.Printf("Meta Lang: %s\n", article.MetaLang)
		fmt.Printf("Content:\n%s\n", article.CleanedText)

	case "yaml":
		bs, err := yaml.Marshal(*article)
		if err != nil {
			return err
		}
		fmt.Print(string(bs))

	default:
		return fmt.Errorf("unrecognized output format %q", OutputFormat)
	}

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information for goose-cli",
	Long:  "All software has versions. This is goose-cli's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("goose-cli HTML Content / Article extractor command-line interface v0.0")
	},
}

func initLogging() {
	if Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
