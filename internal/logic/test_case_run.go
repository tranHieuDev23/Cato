package logic

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	testCaseRunProgramFilePathPlaceHolder = "$PROGRAM"
)

type RunOutput struct {
	ReturnCode int64
	StdOut     string
	StdErr     string
}

type TestCaseRun interface {
	Run(
		ctx context.Context,
		programFilePath string,
		input string,
		timeLimitInMillisecond uint64,
		memoryLimitInByte uint64,
	) (RunOutput, error)
}

type testCaseRun struct {
	dockerClient      *client.Client
	logger            *zap.Logger
	language          string
	testCaseRunConfig configs.TestCaseRun
}

func NewTestCaseRun(
	dockerClient *client.Client,
	logger *zap.Logger,
	language string,
	testCaseRunConfig configs.TestCaseRun,
) (TestCaseRun, error) {
	t := &testCaseRun{
		dockerClient: dockerClient,
		logger: logger.
			With(zap.String("language", language)).
			With(zap.Any("test_case_run_config", testCaseRunConfig)),
		language:          language,
		testCaseRunConfig: testCaseRunConfig,
	}

	if _, err := dockerClient.ImagePull(context.Background(), testCaseRunConfig.Image, types.ImagePullOptions{}); err != nil {
		t.logger.With(zap.Error(err)).Error("failed to load test case run image")
		return nil, err
	}

	return t, nil
}

func (t testCaseRun) getContainerCommand(commandTemplate []string, containerProgramFilePath string) []string {
	command := make([]string, len(commandTemplate))
	for i := range command {
		switch commandTemplate[i] {
		case testCaseRunProgramFilePathPlaceHolder:
			command[i] = containerProgramFilePath
		default:
			command[i] = commandTemplate[i]
		}
	}

	return command
}

func (t testCaseRun) Run(
	ctx context.Context,
	programFilePath string,
	input string,
	timeLimitInMillisecond uint64,
	memoryLimitInByte uint64,
) (RunOutput, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	timeoutDuration := time.Millisecond * time.Duration(timeLimitInMillisecond)
	dockerContainerCtx, dockerContainerCancelFunc := context.WithTimeout(ctx, timeoutDuration)
	defer dockerContainerCancelFunc()

	programFileName := filepath.Base(programFilePath)
	containerProgramFilePath := filepath.Join(t.testCaseRunConfig.WorkingDir, programFileName)
	containerCreateResponse, err := t.dockerClient.ContainerCreate(dockerContainerCtx, &container.Config{
		Image:        t.testCaseRunConfig.Image,
		Cmd:          t.getContainerCommand(t.testCaseRunConfig.CommandTemplate, containerProgramFilePath),
		WorkingDir:   t.testCaseRunConfig.WorkingDir,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", programFilePath, containerProgramFilePath),
		},
		Resources: container.Resources{
			CPUShares: t.testCaseRunConfig.CPUShares,
			Memory:    int64(memoryLimitInByte),
		},
		NetworkMode: "none",
	}, nil, nil, "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create run test case container")
		return RunOutput{}, err
	}

	defer func() {
		if err := t.dockerClient.ContainerRemove(ctx, containerCreateResponse.ID, types.ContainerRemoveOptions{}); err != nil {
			logger.With(zap.Error(err)).Error("failed to remove run test case container")
		}
	}()

	containerID := containerCreateResponse.ID
	containerAttachResponse, err := t.dockerClient.ContainerAttach(dockerContainerCtx, containerID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failed to attached to run test case container")
		return RunOutput{}, err
	}

	defer containerAttachResponse.Close()

	containerAttachResponse.Conn.Write([]byte(input))
	containerAttachResponse.Conn.Write([]byte{'\n'})

	if err := t.dockerClient.ContainerStart(dockerContainerCtx, containerID, types.ContainerStartOptions{}); err != nil {
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failed to attached to run test case container")
		return RunOutput{}, err
	}

	dataChan, errChan := t.dockerClient.ContainerWait(dockerContainerCtx, containerID, container.WaitConditionNotRunning)
	select {
	case data := <-dataChan:
		if data.StatusCode != 0 {
			logger.
				With(zap.Int64("status_code", data.StatusCode)).
				Info("running submission failed: program exited with non-zero code")
		}

		stdoutBuffer := new(bytes.Buffer)
		stderrBuffer := new(bytes.Buffer)
		stdcopy.StdCopy(stdoutBuffer, stderrBuffer, containerAttachResponse.Reader)
		return RunOutput{
			ReturnCode: data.StatusCode,
			StdOut:     stdoutBuffer.String(),
			StdErr:     stderrBuffer.String(),
		}, nil

	case err := <-errChan:
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failure happened while waiting for container")
		return RunOutput{}, err
	}
}
