package main

import (
	"fmt"
	"os"
)

func cleanUp(path string) error {
	return os.RemoveAll(path)
}

func Sync(kubeClient *KubeClient,
	gitClient *GitClient,
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
	author string,
	force bool) error {
	err := syncProcess(kubeClient,
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
		author,
		force)
	if err != nil {
		silentCleanup(destination)
	}
	return err
}

func syncProcess(kubeClient *KubeClient,
	gitClient *GitClient,
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
	author string,
	force bool) error {
	err := gitClient.CloneRepo(origin, destination, username, accessToken, author, email)
	if err != nil {
		return err
	}
	if gitClient.repo == nil {
		return fmt.Errorf("git client not configured")
	}
	if kubeClient.client == nil || kubeClient.restConfig == nil {
		return fmt.Errorf("kube client not configured")
	}
	pod, err := kubeClient.CopyByLabel(podLabel, namespace, container, containerPath, useFirst)
	if err != nil {
		return err
	}
	err = UnTarAll(pod, destination, containerPath)
	if err != nil {
		return fmt.Errorf("error untarring: %v", err)
	}
	clean, err := gitClient.IsClean()
	if err != nil {
		return err
	}
	if clean {
		fmt.Println("no changes to commit")
		return nil
	}
	err = gitClient.AddAll()
	if err != nil {
		return err
	}
	if !dryRyn {
		fmt.Println("committing and pushing changes")
		err := gitClient.CommitAndPush(accessToken, username, force)
		if err != nil {
			return err
		}
	}
	return cleanUp(destination)
}

func silentCleanup(destination string) {
	err := cleanUp(destination)
	if err != nil {
		fmt.Printf("error cleaning: %v", err)
	}
}
