package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/answer1991/docker-images-migration/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os/exec"
	"strings"
)

func newImagePathError(image string) (err error) {
	return fmt.Errorf("Image %s is Invalid, Path Format Must Be $registryDomain/$namespace/$image:$tag", image)
}

func loadAuth() (info *types.AuthInfo, err error) {
	if "" != authFile {
		content, err := ioutil.ReadFile(authFile)

		if nil != err {
			return nil, err
		}

		info := &types.AuthInfo{}
		if err = json.Unmarshal(content, info); nil != err {
			return nil, err
		}

		return info, nil
	}

	return nil, nil
}

func loadSource() (sourceList []*types.Source, err error) {
	if nil != sourceFiles {
		sourceList = make([]*types.Source, len(sourceFiles))

		for i, sourceFile := range sourceFiles {
			content, err := ioutil.ReadFile(sourceFile)

			if nil != err {
				return nil, err
			}

			source := &types.Source{}
			if err = json.Unmarshal(content, source); nil != err {
				return nil, err
			}

			sourceList[i] = source
		}
	}

	return sourceList, nil
}

func getDockerCmdFlag() (flag string) {
	dockerHostCmd := ""

	if "" != dockerHost {
		dockerHostCmd = fmt.Sprintf("-H %s ", dockerHost)
	}

	tlsCmd := ""

	if tls {
		tlsCmd = fmt.Sprintf("--tls --tlscacert %s --tlscert %s --tlskey %s ", tlscacert, tlscert, tlskey)
	}

	return strings.TrimSpace(fmt.Sprintf("docker %s %s", dockerHostCmd, tlsCmd))
}

const defaultImageTag = "latest"
const defaultNamespace = "_"

func getImageNameAndTag(nameAndTag string) (name, tag string, err error) {
	nameAndTagArr := strings.Split(nameAndTag, ":")

	if 1 == len(nameAndTagArr) {
		name = nameAndTagArr[0]
		tag = defaultImageTag
	} else if 2 == len(nameAndTagArr) {
		name = nameAndTagArr[0]
		tag = nameAndTagArr[1]
	} else {
		return "", "", fmt.Errorf("Image And Tag not Correct: %s", nameAndTag)
	}

	return name, tag, nil
}

func getImageInfo(image string) (registryDomain, namespace, name, tag string, err error) {
	arr := strings.Split(image, "/")

	if 1 == len(arr) {
		registryDomain = ""
		namespace = defaultNamespace
		name, tag, err = getImageNameAndTag(arr[0])
	} else if 2 == len(arr) {
		registryDomain = ""
		namespace = arr[1]
		name, tag, err = getImageNameAndTag(arr[1])
	} else if 3 == len(arr) {
		registryDomain = arr[0]
		namespace = arr[1]
		name, tag, err = getImageNameAndTag(arr[2])
	} else {
		return "", "", "", "", newImagePathError(image)
	}

	return registryDomain, namespace, name, tag, err
}

func getTarRegistryDomains(sourceList []*types.Source) (domains []string, err error) {
	countMap := make(map[string]int)

	for _, source := range sourceList {
		for _, image := range source.Images {
			registryDomain, _, _, _, err := getImageInfo(image)
			if nil != err {
				return nil, err
			}

			if v, ok := countMap[registryDomain]; ok {
				countMap[registryDomain] = v + 1
			} else {
				countMap[registryDomain] = 1
			}
		}
	}

	domains = make([]string, len(countMap))

	i := 0
	for key := range countMap {
		domains[i] = key
		i++
	}

	return domains, nil
}

func getImportRegistryDomains(sourceList []*types.Source) (domains []string, err error) {
	countMap := make(map[string]int)

	for _, source := range sourceList {
		if v, ok := countMap[source.TargetRegistryDomain]; ok {
			countMap[source.TargetRegistryDomain] = v + 1
		} else {
			countMap[source.TargetRegistryDomain] = 1
		}
	}

	domains = make([]string, len(countMap))

	i := 0
	for key := range countMap {
		domains[i] = key
		i++
	}

	return domains, nil
}

func getImages(sourceList []*types.Source) (images []string, err error) {
	countMap := make(map[string]int)

	for _, source := range sourceList {
		for _, image := range source.Images {
			if v, ok := countMap[image]; ok {
				countMap[image] = v + 1
			} else {
				countMap[image] = 1
			}
		}
	}

	images = make([]string, len(countMap))

	i := 0
	for key := range countMap {
		images[i] = key
		i++
	}

	return images, nil
}

func translateTargetImages(cmd *cobra.Command, sourceList []*types.Source) (images []string, err error) {
	docker := getDockerCmdFlag()

	images = make([]string, 0)

	for _, source := range sourceList {
		for _, image := range source.Images {
			domain, namespace, name, tag, err := getImageInfo(image)
			if nil != err {
				return nil, err
			}

			newDomain := source.TargetRegistryDomain

			if "" == newDomain {
				newDomain = domain
			}

			newImage := fmt.Sprintf("%s/%s/%s:%s", source.TargetRegistryDomain, namespace, name, tag)

			if image != newImage {
				command := exec.Command(docker, "tag", image, newImage)
				command.Stdout = cmd.OutOrStdout()

				if err := command.Run(); nil != err {
					stderr(cmd, err)
					return nil, err
				}
			}

			images = append(images, newImage)
		}
	}

	return images, nil
}

func loginRegistry(cmd *cobra.Command, info *types.AuthInfo, registry string) (err error) {
	if nil != info && "" != registry {
		docker := getDockerCmdFlag()
		command := exec.Command(docker, "login", "--username", info.Username, "--password", info.Password, registry)

		command.Stdin = strings.NewReader(info.Password)
		command.Stdout = cmd.OutOrStdout()

		if err := command.Run(); nil != err {
			stderr(cmd, err)

			return err
		} else {
			stdout(cmd, fmt.Sprintf("Login Success to %s", registry))
		}
	}

	return nil
}

func logoutRegistries(cmd *cobra.Command, registryDomains []string) {
	for _, registry := range registryDomains {
		if "" != registry {
			docker := getDockerCmdFlag()
			command := exec.Command(docker, "logout", registry)
			command.Stdout = cmd.OutOrStdout()

			if err := command.Run(); nil != err {
				stderr(cmd, err)
			}
		}
	}
}

func stderr(cmd *cobra.Command, err error) {
	if nil != err {
		cmd.OutOrStderr().Write([]byte(fmt.Sprintln(err.Error())))
	}
}

func stdout(cmd *cobra.Command, msgs ...interface{}) {
	for _, msg := range msgs {
		cmd.OutOrStderr().Write([]byte(fmt.Sprintln(msg)))
	}
}

func pullImages(cmd *cobra.Command, images []string) (err error) {
	docker := getDockerCmdFlag()
	for _, image := range images {
		command := exec.Command(docker, "pull", image)
		command.Stdout = cmd.OutOrStdout()

		if err := command.Run(); nil != err {
			stderr(cmd, err)
			return err
		}
	}

	stdout(cmd, "Pull Images Success")

	return nil
}

func pushImages(cmd *cobra.Command, images []string) (err error) {
	docker := getDockerCmdFlag()
	for _, image := range images {
		command := exec.Command(docker, "push", image)
		command.Stdout = cmd.OutOrStdout()

		if err := command.Run(); nil != err {
			stderr(cmd, err)
			return err
		}
	}

	stdout(cmd, "Push Images Success")

	return nil
}

func tarImages(cmd *cobra.Command, images []string) (err error) {
	docker := getDockerCmdFlag()

	args := []string{"save", "-o", target}
	for _, image := range images {
		args = append(args, image)
	}

	command := exec.Command(docker, args...)
	command.Stdout = cmd.OutOrStdout()

	if err := command.Run(); nil != err {
		stderr(cmd, err)
		return err
	}

	return nil
}

func loadTar(cmd *cobra.Command) (err error) {
	docker := getDockerCmdFlag()
	command := exec.Command(docker, "load", "-i", target)
	command.Stdout = cmd.OutOrStdout()

	if err := command.Run(); nil != err {
		stderr(cmd, err)
		return err
	}

	return nil
}
