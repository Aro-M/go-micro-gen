package cli

import (
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Aro-M/go-micro-gen/internal/config"
	"github.com/Aro-M/go-micro-gen/internal/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inject microservice structure into the current directory",
	Long:  "Scaffold all standard go-micro-gen templates directly into your existing project.",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg := &config.ServiceConfig{
			GoVersion: goVersion(),
			OutputDir: ".",
		}

		// 1. Determine service name from current directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		cfg.ServiceName = filepath.Base(cwd)

		// 2. Discover module path from go.mod if it exists
		if modBytes, err := os.ReadFile("go.mod"); err == nil {
			if modPath := modfile.ModulePath(modBytes); modPath != "" {
				cfg.ModulePath = modPath
			}
		}

		if cfg.ModulePath == "" {
			cfg.ModulePath = "github.com/acme/" + cfg.ServiceName
		}

		// 3. Ask standard questions, skipping Name, Module, OutputDir
		if err := askArch(cfg); err != nil {
			return err
		}
		if err := askDB(cfg); err != nil {
			return err
		}
		if err := askBroker(cfg); err != nil {
			return err
		}
		if err := askTransport(cfg); err != nil {
			return err
		}
		if err := askGraphQL(cfg); err != nil {
			return err
		}
		if err := askJWT(cfg); err != nil {
			return err
		}
		if err := askServerless(cfg); err != nil {
			return err
		}
		if err := askSeeding(cfg); err != nil {
			return err
		}
		if err := askRedis(cfg); err != nil {
			return err
		}
		if err := askDocker(cfg); err != nil {
			return err
		}
		if err := askPrometheus(cfg); err != nil {
			return err
		}
		if err := askGrafana(cfg); err != nil {
			return err
		}
		if err := askK8s(cfg); err != nil {
			return err
		}
		if err := askHelm(cfg); err != nil {
			return err
		}
		if err := askCloud(cfg); err != nil {
			return err
		}
		if err := askCI(cfg); err != nil {
			return err
		}

		printSummary(cfg)

		if !flagYes {
			confirm := false
			if err := survey.AskOne(&survey.Confirm{
				Message: "Inject service structure with these settings into current directory?",
				Default: true,
			}, &confirm); err != nil || !confirm {
				color.Yellow("Aborted.")
				return nil
			}
		}

		g := generator.New(cfg)
		color.Cyan("\n🚀 Initializing structured project in '%s' ...\n", cfg.ServiceName)

		if err := g.Generate(); err != nil {
			color.Red("❌ Initialization failed: %v", err)
			return err
		}

		printSuccess(cfg)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&flagDB, "db", "", "Database type: postgres | mongo | none")
	initCmd.Flags().StringVar(&flagBroker, "broker", "", "Message Broker: kafka | rabbitmq | nats | none")
	initCmd.Flags().StringVar(&flagTransport, "transport", "", "Transport protocol: http | grpc | both")
	initCmd.Flags().StringVar(&flagArch, "arch", "", "Architecture: clean | hexagonal | ddd | vertical | standard")
	initCmd.Flags().StringVar(&flagCI, "ci", "", "CI/CD: github | gitlab | none")
	initCmd.Flags().StringVar(&flagCloud, "cloud", "", "Cloud Provider: aws | gcp | none")
	initCmd.Flags().BoolVar(&flagRedis, "redis", false, "Include Redis")
	initCmd.Flags().BoolVar(&flagGraphQL, "graphql", false, "Include GraphQL endpoint")
	initCmd.Flags().BoolVar(&flagJWT, "jwt", false, "Include JWT Auth Middleware")
	initCmd.Flags().BoolVar(&flagServerless, "serverless", false, "Include Serverless deployment wrappers")
	initCmd.Flags().BoolVar(&flagSeeding, "seeding", false, "Include DB mock seeder wrappers")
	initCmd.Flags().BoolVar(&flagDocker, "docker", false, "Include Docker setup")
	initCmd.Flags().BoolVar(&flagPrometheus, "prometheus", false, "Include Prometheus & OpenTelemetry")
	initCmd.Flags().BoolVar(&flagGrafana, "grafana", false, "Include Grafana dashboard")
	initCmd.Flags().BoolVar(&flagK8s, "k8s", false, "Include Kubernetes manifests")
	initCmd.Flags().BoolVar(&flagHelm, "helm", false, "Include Helm charts")
	initCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation prompt")

	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		flagRedisSet = initCmd.Flags().Changed("redis")
		flagGraphQLSet = initCmd.Flags().Changed("graphql")
		flagJWTSet = initCmd.Flags().Changed("jwt")
		flagServerlessSet = initCmd.Flags().Changed("serverless")
		flagSeedingSet = initCmd.Flags().Changed("seeding")
		flagDockerSet = initCmd.Flags().Changed("docker")
		flagPrometheusSet = initCmd.Flags().Changed("prometheus")
		flagGrafanaSet = initCmd.Flags().Changed("grafana")
		flagK8sSet = initCmd.Flags().Changed("k8s")
		flagHelmSet = initCmd.Flags().Changed("helm")
	}
}
