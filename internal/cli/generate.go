package cli

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Aro-M/go-micro-gen/internal/config"
	"github.com/Aro-M/go-micro-gen/internal/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagName          string
	flagModule        string
	flagDB            string
	flagBroker        string
	flagTransport     string
	flagArch          string
	flagCI            string
	flagRedis         bool
	flagRedisSet      bool
	flagGraphQL       bool
	flagGraphQLSet    bool
	flagJWT           bool
	flagJWTSet        bool
	flagServerless    bool
	flagServerlessSet bool
	flagSeeding       bool
	flagSeedingSet    bool
	flagDocker        bool
	flagDockerSet     bool
	flagGrafana       bool
	flagGrafanaSet    bool
	flagK8s           bool
	flagK8sSet        bool
	flagHelm          bool
	flagHelmSet       bool
	flagCloud         string
	flagOutput        string
	flagYes           bool // skip confirmation prompt (for CI / scripted usage)
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new production-ready Go microservice",
	Long:  "Interactively scaffold a fully-wired microservice with Observability, DB, Docker, CI/CD and more.",
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().StringVar(&flagName, "name", "", "Service name (e.g. order-service)")
	generateCmd.Flags().StringVar(&flagModule, "module", "", "Go module path (e.g. github.com/acme/order-service)")
	generateCmd.Flags().StringVar(&flagDB, "db", "", "Database type: postgres | mongo | none")
	generateCmd.Flags().StringVar(&flagBroker, "broker", "", "Message Broker: kafka | rabbitmq | nats | none")
	generateCmd.Flags().StringVar(&flagTransport, "transport", "", "Transport protocol: http | grpc | both")
	generateCmd.Flags().StringVar(&flagArch, "arch", "", "Architecture: clean | hexagonal | ddd | vertical | standard")
	generateCmd.Flags().StringVar(&flagCI, "ci", "", "CI/CD: github | gitlab | none")
	generateCmd.Flags().StringVar(&flagCloud, "cloud", "", "Cloud Provider: aws | gcp | none")
	generateCmd.Flags().BoolVar(&flagRedis, "redis", false, "Include Redis")
	generateCmd.Flags().BoolVar(&flagGraphQL, "graphql", false, "Include GraphQL endpoint")
	generateCmd.Flags().BoolVar(&flagJWT, "jwt", false, "Include JWT Auth Middleware")
	generateCmd.Flags().BoolVar(&flagServerless, "serverless", false, "Include Serverless deployment wrappers")
	generateCmd.Flags().BoolVar(&flagSeeding, "seeding", false, "Include DB mock seeder wrappers")
	generateCmd.Flags().BoolVar(&flagDocker, "docker", false, "Include Docker setup")
	generateCmd.Flags().BoolVar(&flagGrafana, "grafana", false, "Include Grafana dashboard")
	generateCmd.Flags().BoolVar(&flagK8s, "k8s", false, "Include Kubernetes manifests")
	generateCmd.Flags().BoolVar(&flagHelm, "helm", false, "Include Helm charts")
	generateCmd.Flags().StringVar(&flagOutput, "output", "", "Output directory")
	generateCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation prompt")

	generateCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Check if flags were explicitly set
		flagRedisSet = generateCmd.Flags().Changed("redis")
		flagGraphQLSet = generateCmd.Flags().Changed("graphql")
		flagJWTSet = generateCmd.Flags().Changed("jwt")
		flagServerlessSet = generateCmd.Flags().Changed("serverless")
		flagSeedingSet = generateCmd.Flags().Changed("seeding")
		flagDockerSet = generateCmd.Flags().Changed("docker")
		flagGrafanaSet = generateCmd.Flags().Changed("grafana")
		flagK8sSet = generateCmd.Flags().Changed("k8s")
		flagHelmSet = generateCmd.Flags().Changed("helm")
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	cfg := &config.ServiceConfig{
		GoVersion: goVersion(),
	}

	// --- Collect inputs (interactive if flag not provided) ---
	if err := askName(cfg); err != nil {
		return err
	}
	if err := askModule(cfg); err != nil {
		return err
	}
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
	if err := askRedis(cfg); err != nil {
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
	if err := askDocker(cfg); err != nil {
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
	if err := askOutput(cfg); err != nil {
		return err
	}

	// --- Summary + Confirmation ---
	printSummary(cfg)

	if !flagYes {
		confirm := false
		if err := survey.AskOne(&survey.Confirm{
			Message: "Generate service with these settings?",
			Default: true,
		}, &confirm); err != nil || !confirm {
			color.Yellow("Aborted.")
			return nil
		}
	}

	// --- Generate ---
	g := generator.New(cfg)
	color.Cyan("\n🚀 Generating %s ...\n", cfg.ServiceName)
	if err := g.Generate(); err != nil {
		color.Red("❌ Generation failed: %v", err)
		return err
	}

	printSuccess(cfg)
	return nil
}

func askName(cfg *config.ServiceConfig) error {
	if flagName != "" {
		cfg.ServiceName = flagName
		return nil
	}
	return survey.AskOne(&survey.Input{
		Message: "Service name:",
		Help:    "Lowercase, hyphenated (e.g. order-service)",
	}, &cfg.ServiceName, survey.WithValidator(survey.Required))
}

func askModule(cfg *config.ServiceConfig) error {
	if flagModule != "" {
		cfg.ModulePath = flagModule
		return nil
	}
	defaultModule := fmt.Sprintf("github.com/acme/%s", cfg.ServiceName)
	err := survey.AskOne(&survey.Input{
		Message: "Go module path:",
		Default: defaultModule,
	}, &cfg.ModulePath)
	if cfg.ModulePath == "" {
		cfg.ModulePath = defaultModule
	}
	return err
}

func askArch(cfg *config.ServiceConfig) error {
	if flagArch != "" {
		cfg.Architecture = config.ArchType(flagArch)
		return nil
	}
	var answer string
	err := survey.AskOne(&survey.Select{
		Message: "Architecture pattern:",
		Options: []string{"clean", "hexagonal", "ddd", "vertical", "standard"},
		Default: "clean",
	}, &answer)
	cfg.Architecture = config.ArchType(answer)
	return err
}

func askDB(cfg *config.ServiceConfig) error {
	if flagDB != "" {
		cfg.Database = config.DBType(flagDB)
		return nil
	}
	var answer string
	err := survey.AskOne(&survey.Select{
		Message: "Database:",
		Options: []string{"postgres", "mongo", "none"},
		Default: "postgres",
	}, &answer)
	cfg.Database = config.DBType(answer)
	return err
}

func askBroker(cfg *config.ServiceConfig) error {
	if flagBroker != "" {
		cfg.Broker = config.BrokerType(flagBroker)
		return nil
	}
	var answer string
	err := survey.AskOne(&survey.Select{
		Message: "Message Broker:",
		Options: []string{"kafka", "rabbitmq", "nats", "none"},
		Default: "none",
	}, &answer)
	cfg.Broker = config.BrokerType(answer)
	return err
}

func askTransport(cfg *config.ServiceConfig) error {
	if flagTransport != "" {
		cfg.Transport = config.TransportType(flagTransport)
		return nil
	}
	var answer string
	err := survey.AskOne(&survey.Select{
		Message: "Transport Protocol:",
		Options: []string{"http", "grpc", "both", "none"},
		Default: "http",
	}, &answer)
	cfg.Transport = config.TransportType(answer)
	return err
}

func askRedis(cfg *config.ServiceConfig) error {
	if flagRedisSet {
		cfg.IncludeRedis = flagRedis
		return nil
	}
	var res bool
	prompt := &survey.Confirm{
		Message: "Include Redis?",
		Default: false,
	}
	if err := survey.AskOne(prompt, &res); err != nil {
		return err
	}
	cfg.IncludeRedis = res
	return nil
}

func askGraphQL(cfg *config.ServiceConfig) error {
	if flagGraphQLSet {
		cfg.IncludeGraphQL = flagGraphQL
	} else {
		var res bool
		prompt := &survey.Confirm{
			Message: "Include GraphQL endpoint?",
			Default: false,
		}
		if err := survey.AskOne(prompt, &res); err != nil {
			return err
		}
		cfg.IncludeGraphQL = res
	}

	// Validate: GraphQL requires HTTP
	if cfg.IncludeGraphQL && cfg.Transport == config.TransportGRPC {
		color.Yellow("\n⚠️  GraphQL requires HTTP transport. Automatically switching Transport to 'both' (HTTP + gRPC).")
		cfg.Transport = config.TransportBoth
	}

	return nil
}

func askJWT(cfg *config.ServiceConfig) error {
	if flagJWTSet {
		cfg.IncludeJWT = flagJWT
		return nil
	}
	var res bool
	prompt := &survey.Confirm{
		Message: "Include JWT Auth Middleware?",
		Default: false,
	}
	if err := survey.AskOne(prompt, &res); err != nil {
		return err
	}
	cfg.IncludeJWT = res
	return nil
}

func askServerless(cfg *config.ServiceConfig) error {
	if flagServerlessSet {
		cfg.IncludeServerless = flagServerless
		return nil
	}
	var res bool
	prompt := &survey.Confirm{
		Message: "Support Serverless deployment? (AWS Lambda / GCP Functions)",
		Default: false,
	}
	if err := survey.AskOne(prompt, &res); err != nil {
		return err
	}
	cfg.IncludeServerless = res
	return nil
}

func askSeeding(cfg *config.ServiceConfig) error {
	if cfg.Database == config.DBNone {
		cfg.IncludeSeeding = false
		return nil
	}
	if flagSeedingSet {
		cfg.IncludeSeeding = flagSeeding
		return nil
	}
	var res bool
	prompt := &survey.Confirm{
		Message: "Add database seeding script (mock data generation)?",
		Default: false,
	}
	if err := survey.AskOne(prompt, &res); err != nil {
		return err
	}
	cfg.IncludeSeeding = res
	return nil
}

func askDocker(cfg *config.ServiceConfig) error {
	if flagDockerSet {
		cfg.IncludeDocker = flagDocker
		return nil
	}
	return survey.AskOne(&survey.Confirm{
		Message: "Include Docker & Docker Compose setup?",
		Default: true,
	}, &cfg.IncludeDocker)
}

func askGrafana(cfg *config.ServiceConfig) error {
	if !cfg.IncludeDocker {
		cfg.IncludeGrafana = false
		return nil
	}
	if flagGrafanaSet {
		cfg.IncludeGrafana = flagGrafana
		return nil
	}
	return survey.AskOne(&survey.Confirm{
		Message: "Include Grafana dashboard?",
		Default: true,
	}, &cfg.IncludeGrafana)
}

func askK8s(cfg *config.ServiceConfig) error {
	if flagK8sSet {
		cfg.IncludeK8s = flagK8s
		return nil
	}
	return survey.AskOne(&survey.Confirm{
		Message: "Include Kubernetes manifests (Deployment, Service, etc.)?",
		Default: false,
	}, &cfg.IncludeK8s)
}

func askHelm(cfg *config.ServiceConfig) error {
	if flagHelmSet {
		cfg.IncludeHelm = flagHelm
		return nil
	}
	// Need survey check
	return survey.AskOne(&survey.Confirm{
		Message: "Include Helm charts?",
		Default: false,
	}, &cfg.IncludeHelm)
}

func askCloud(cfg *config.ServiceConfig) error {
	if flagCloud != "" {
		cfg.Cloud = config.CloudProvider(flagCloud)
		return nil
	}
	var answer string
	err := survey.AskOne(&survey.Select{
		Message: "Cloud provider:",
		Options: []string{"aws", "gcp", "none"},
		Default: "none",
	}, &answer)
	cfg.Cloud = config.CloudProvider(answer)
	return err
}

func askCI(cfg *config.ServiceConfig) error {
	if flagCI != "" {
		cfg.CI = config.CIType(flagCI)
		return nil
	}
	var answer string
	err := survey.AskOne(&survey.Select{
		Message: "CI/CD provider:",
		Options: []string{"github", "gitlab", "none"},
		Default: "github",
	}, &answer)
	cfg.CI = config.CIType(answer)
	return err
}

func askOutput(cfg *config.ServiceConfig) error {
	if flagOutput != "" {
		cfg.OutputDir = flagOutput
		return nil
	}
	defaultOut := "./" + cfg.ServiceName
	return survey.AskOne(&survey.Input{
		Message: "Output directory:",
		Default: defaultOut,
	}, &cfg.OutputDir)
}

func printSummary(cfg *config.ServiceConfig) {
	bold := color.New(color.Bold).SprintFunc()
	fmt.Println()
	fmt.Printf("  %s  %s\n", bold("Service:"), cfg.ServiceName)
	fmt.Printf("  %s  %s\n", bold("Module: "), cfg.ModulePath)
	fmt.Printf("  %s  %s\n", bold("Arch:   "), cfg.Architecture)
	fmt.Printf("  %s  %s\n", bold("DB:     "), cfg.Database)
	fmt.Printf("  %s  %s\n", bold("Broker: "), cfg.Broker)
	fmt.Printf("  Transp:   %s\n", color.CyanString(string(cfg.Transport)))
	fmt.Printf("  GraphQL:  %s\n", color.CyanString(fmt.Sprintf("%t", cfg.IncludeGraphQL)))
	fmt.Printf("  JWT:      %s\n", color.CyanString(fmt.Sprintf("%t", cfg.IncludeJWT)))
	fmt.Printf("  Redis:    %s\n", color.CyanString(fmt.Sprintf("%t", cfg.IncludeRedis)))
	fmt.Printf("  %s  %v\n", bold("Docker: "), cfg.IncludeDocker)
	fmt.Printf("  %s  %v\n", bold("Grafana:"), cfg.IncludeGrafana)
	fmt.Printf("  %s  %v\n", bold("K8s:    "), cfg.IncludeK8s)
	fmt.Printf("  %s  %v\n", bold("Helm:   "), cfg.IncludeHelm)
	fmt.Printf("  %s  %s\n", bold("Cloud:  "), cfg.Cloud)
	if cfg.IncludeServerless {
		fmt.Printf("  %s  %s\n", bold("S-less: "), color.CyanString("true"))
	}
	if cfg.IncludeSeeding {
		fmt.Printf("  %s  %s\n", bold("Seed:   "), color.CyanString("true"))
	}
	fmt.Printf("  %s  %s\n", bold("CI:     "), cfg.CI)
	fmt.Printf("  %s  %s\n", bold("Output: "), cfg.OutputDir)
	fmt.Println()
}

func printSuccess(cfg *config.ServiceConfig) {
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Println()
	fmt.Println(green("✅ Service generated successfully!"))
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Printf("    %s\n", cyan(fmt.Sprintf("cd %s", cfg.OutputDir)))
	fmt.Printf("    %s\n", cyan("make up          # Start all services with Docker Compose"))
	fmt.Printf("    %s\n", cyan("make run         # Run the service locally"))
	fmt.Printf("    %s\n", cyan("make test        # Run tests"))
	fmt.Printf("    %s\n", cyan("make lint        # Run golangci-lint"))
	fmt.Println()
	if cfg.IncludeGrafana {
		fmt.Printf("  Grafana:    %s\n", cyan("http://localhost:3000  (admin/admin)"))
	}
	fmt.Printf("  Prometheus: %s\n", cyan("http://localhost:9090"))
	fmt.Printf("  Service:    %s\n", cyan("http://localhost:8080"))
	fmt.Println()
}

func goVersion() string {
	v := runtime.Version() // e.g. "go1.22.0"
	v = strings.TrimPrefix(v, "go")
	parts := strings.Split(v, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return "1.22"
}
