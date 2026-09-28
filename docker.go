package main

import (
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type DockerTool struct {
	Name string
}

type DockerStats struct {
	Containers int
	Running    int
	Paused     int
	Stopped    int
	Images     int
}

type Containers struct {
	ID         string
	IMAGE      string
	COMMAND    string
	CREATED    string
	STATUS     string
	LOCALPORT  string
	REMOTEPORT string
	NAMES      string
	IP         string
}

func (dt DockerTool) RunDocker(args ...string) string {

	cmd := exec.Command("docker", args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "Error running command: " + err.Error()
	}

	return string(output)
}

func (dt DockerTool) getSystemStats() DockerStats {
	var stats DockerStats
	rawInfo := strings.Split(dt.RunDocker("info"), "\n")

	for _, line := range rawInfo {

		if strings.Contains(line, ":") {

			parts := strings.SplitN(line, ":", 2)

			key := strings.TrimSpace(parts[0])
			value, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				continue
			}

			switch key {
			case "Containers":
				stats.Containers = value
			case "Running":
				stats.Running = value
			case "Paused":
				stats.Paused = value
			case "Stopped":
				stats.Stopped = value
			case "Images":
				stats.Images = value
			}
		}
	}

	return stats
}

func (dt DockerTool) getDockerContainers() []Containers {

	var containerList []Containers

	rawData := dt.RunDocker(
		"ps",
		"-a",
		"--format",
		"{{.ID}}|{{.Image}}|{{.Command}}|{{.CreatedAt}}|{{.Status}}|{{.Ports}}|{{.Names}}",
	)

	lines := strings.Split(rawData, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var dockerContainer Containers
		filteredOutput := strings.Split(line, "|")

		dockerContainer.ID = filteredOutput[0]
		dockerContainer.IMAGE = filteredOutput[1]
		dockerContainer.COMMAND = filteredOutput[2]
		dockerContainer.CREATED = filteredOutput[3]
		dockerContainer.STATUS = filteredOutput[4]
		dockerContainer.NAMES = filteredOutput[6]

		if strings.Contains(filteredOutput[5], "->") {
			port := strings.Split(filteredOutput[5], ",")[0]
			ports := strings.SplitN(port, "->", 2)

			dockerContainer.LOCALPORT = ports[0]
			dockerContainer.REMOTEPORT = ports[1]
		}

		containerList = append(containerList, dockerContainer)
	}
	if len(containerList) == 0 {
		fmt.Println("No containers found.")
		return nil
	}
	return containerList
}

func (dt DockerTool) printContainerList() {

	containers := dt.getDockerContainers()

	for _, c := range containers {
		fmt.Printf("ID: %s | Name: %s | Status: %s | Local: %s| Remote: %s| Image: %s\n",
			c.ID, c.NAMES, c.STATUS, c.LOCALPORT, c.REMOTEPORT, c.IMAGE)
	}

}

func (dt DockerTool) changeExistingContainerState() bool {

	containers := dt.getDockerContainers()

	if len(containers) == 0 {
		return false
	}
	selectionOptions := make([]string, 0, len(containers))

	for i, c := range containers {
		fmt.Printf("%d) %-20s | %-25s | %s\n", i+1, c.NAMES, c.STATUS, c.IMAGE)
		selectionOptions = append(selectionOptions, c.ID)
	}

	for {
		fmt.Print("Selection > ")
		var containerChoice int
		fmt.Scan(&containerChoice)

		targetIndex := containerChoice - 1
		if targetIndex < 0 || targetIndex >= len(selectionOptions) {
			fmt.Println("Invalid container selection.")
			continue
		}
		selectedID := selectionOptions[targetIndex]

		fmt.Print(`
What would you like to do with this container?

1: Start
2: Stop
3: Restart
4: Delete
5: Go back to main menu

> `)

		var actionChoice string
		fmt.Scan(&actionChoice)

		switch actionChoice {
		case "1":
			fmt.Println("Starting container...")
			go dt.RunDocker("start", selectedID)
			return true
		case "2":
			fmt.Println("Stopping container...this can take up to ten seconds. Refresh in a few seconds")
			go dt.RunDocker("stop", selectedID)
			return true
		case "3":
			fmt.Println("Restarting container...")
			dt.RunDocker("restart", selectedID)
			return true
		case "4":
			fmt.Println("Deleting container...")
			out := dt.RunDocker("rm", "-f", selectedID)
			fmt.Println(out)
			return true
		case "5":
			fmt.Println("Going back to main menu...")
			return false
		default:
			fmt.Println("Invalid selection, taking you back to the main menu")
			return false
		}
	}
}

func (dt DockerTool) createNewContainer() bool {

	containers := dt.getDockerContainers()
	portsInUse, containerNames := []string{}, []string{}

	for _, c := range containers {
		port := c.LOCALPORT
		if idx := strings.LastIndex(port, ":"); idx != -1 {
			port = port[idx+1:]
		}
		if port != "" {
			portsInUse = append(portsInUse, port)
		}
		containerNames = append(containerNames, c.NAMES)
	}

	localPort := 10000
	var portStr string

	for {
		portStr = strconv.Itoa(localPort)
		if !slices.Contains(portsInUse, portStr) {
			break
		}
		localPort++
	}

	baseName := "Debian-"
	i := 1
	var hostname string

	for {
		hostname = fmt.Sprintf("%s%d", baseName, i)
		if !slices.Contains(containerNames, hostname) {
			break
		}
		i++
	}

	image := "lab-debian:latest"

	dt.RunDocker(
		"run",
		"-d",
		"--name",
		hostname,
		"-p",
		strconv.Itoa(localPort)+":22",
		image,
	)

	fmt.Printf("Container deployed %s at %d\n", hostname, localPort)
	return true
}
