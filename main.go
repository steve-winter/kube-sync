package main

import (
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	//cmd.Execute()

	viper.SetConfigName("config") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")        // config current file path
	viper.AddConfigPath("/config/")         // config directory for kubernetes
	viper.AddConfigPath("$HOME/.kube-sync") // call multiple times to add many search paths

	viper.AutomaticEnv() // read value ENV variable
	viper.SetDefault("InCluster", true)
	viper.SetDefault("DryRun", false)
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	var config KubeConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to umarshall config: %v", err)
		return
	}

	configFile := viper.ConfigFileUsed()
	fmt.Printf("Config file used: %v\n", configFile)

	for _, entry := range config.ConfigEntryList {
		fmt.Printf("Running job: %v\n", entry.Name)
		err := validateItem(entry)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		if entry.GitItem.GitDestination == "" {
			dest := randomString()
			fmt.Printf("Git destination randomly set: %v\n", dest)
			entry.GitItem.GitDestination = dest
		}
		if entry.GitItem.GitAuthor == "" {
			author := "Kube Sync"
			fmt.Printf("Git destination randomly set: %v\n", author)
			entry.GitItem.GitAuthor = author
		}
		err = runner(entry.Name,
			config.InCluster,
			entry.GitItem.GitOrigin,
			entry.GitItem.GitDestination,
			entry.PodItem.PodLabel,
			entry.PodItem.PodNamespace,
			entry.PodItem.PodContainer,
			entry.PodItem.PodContainerPath,
			entry.PodItem.UseFirstPod,
			config.DryRun,
			entry.GitItem.GitUsername,
			entry.GitItem.GitAccessToken,
			entry.GitItem.GitEmail,
			entry.GitItem.GitAuthor)
		if err != nil {
			fmt.Printf("Failure running job: %v, error: %v\n", entry.Name, err)
			return
		}
	}
}

func randomString() string {
	rand.Seed(time.Now().Unix())
	randomInt := rand.Intn(10)
	str := strconv.Itoa(randomInt)
	return str
}

func validateItem(entry ConfigEntry) error {
	if entry.GitItem.GitUsername == "" || entry.GitItem.GitOrigin == "" || entry.GitItem.GitAccessToken == "" {
		return fmt.Errorf("git origin, username and accesstoken must be set")
	}
	if entry.PodItem.PodNamespace == "" || entry.PodItem.PodLabel == "" || entry.PodItem.PodContainer == "" || entry.PodItem.PodContainerPath == "" {
		return fmt.Errorf("pod namespace, Label, Container and Container path must be set")
	}
	return nil
}

func runner(name string, inCluster bool,
	origin string,
	destination string,
	podLabel string,
	namespace string,
	container string,
	containerPath string,
	useFirst bool,
	dryRyn bool,
	username string,
	accessToken string,
	email string,
	author string) error {
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("unable to load cluster config: %v", err)
		}
	} else {
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return fmt.Errorf("unable to load local cluster config: %v", err)
		}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("unable to build kube client: %v", err)
	}
	kubeClient := &KubeClient{
		client:     client,
		restConfig: config,
	}
	gitClient := &GitClient{}

	err = Sync(
		kubeClient,
		gitClient,
		origin,
		destination,
		podLabel,
		namespace,
		container,
		containerPath,
		useFirst,
		dryRyn,
		username,
		accessToken,
		email,
		author)
	if err != nil {
		return err
	}
	fmt.Printf("completed job %v", name)
	return nil
}
