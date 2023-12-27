package logic

import (
	"bytes"
	"context"
	"errors"
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
	testCaseRunProgramFilePathPlaceHolder    = "$PROGRAM"
	testCaseRunTimeLimitInSecondsPlaceHolder = "$TIME_LIMIT"

	millisecondPerSecond   = 1000
	timeoutErrorReturnCode = 124
	oomErrorReturnCode     = 137
)

type RunOutput struct {
	ReturnCode          int64
	TimeLimitExceeded   bool
	MemoryLimitExceeded bool
	StdOut              string
	StdErr              string
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
	appArguments      utils.Arguments
}

func (t testCaseRun) pullImage() error {
	t.logger.Info("pulling test case run image")
	_, err := t.dockerClient.ImagePull(context.Background(), t.testCaseRunConfig.Image, types.ImagePullOptions{})
	if err != nil {
		t.logger.With(zap.Error(err)).Error("failed to pull test case run image")
		return err
	}

	t.logger.Info("pulled test case run image successfully")
	return nil
}

func NewTestCaseRun(
	dockerClient *client.Client,
	logger *zap.Logger,
	language string,
	testCaseRunConfig configs.TestCaseRun,
	appArguments utils.Arguments,
) (TestCaseRun, error) {
	t := &testCaseRun{
		dockerClient: dockerClient,
		logger: logger.
			With(zap.String("language", language)).
			With(zap.Any("test_case_run_config", testCaseRunConfig)),
		language:          language,
		testCaseRunConfig: testCaseRunConfig,
		appArguments:      appArguments,
	}

	if appArguments.PullImageAtStartUp {
		if err := t.pullImage(); err != nil {
			return nil, err
		}
	} else {
		go func() {
			_ = t.pullImage()
		}()
	}

	return t, nil
}

func (t testCaseRun) getWorkingDir() string {
	if t.testCaseRunConfig.WorkingDir != "" {
		return t.testCaseRunConfig.WorkingDir
	}

	return "/work"
}

func (t testCaseRun) getContainerCommand(
	commandTemplate []string,
	containerProgramFilePath string,
	timeLimitInMillisecond uint64,
) []string {
	timeLimitInSecondString := fmt.Sprintf("%.3fs", float64(timeLimitInMillisecond)/millisecondPerSecond)

	command := make([]string, len(commandTemplate))
	for i := range command {
		switch commandTemplate[i] {
		case testCaseRunProgramFilePathPlaceHolder:
			command[i] = containerProgramFilePath
		case testCaseRunTimeLimitInSecondsPlaceHolder:
			command[i] = timeLimitInSecondString
		default:
			command[i] = commandTemplate[i]
		}
	}

	return command
}

func (t testCaseRun) onContainerWaitData(
	ctx context.Context,
	data container.WaitResponse,
	containerAttachResponse types.HijackedResponse,
) (RunOutput, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	if _, err := stdcopy.StdCopy(stdoutBuffer, stderrBuffer, containerAttachResponse.Reader); err != nil {
		logger.With(zap.Error(err)).Error("failed to read from stdout and stderr of container")
		return RunOutput{}, err
	}

	stdOut := utils.TrimSpaceRight(stdoutBuffer.String())
	stdErr := utils.TrimSpaceRight(stderrBuffer.String())

	switch data.StatusCode {
	case 0:
		logger.Info("running submission successfully")
		return RunOutput{
			StdOut: stdOut,
			StdErr: stdErr,
		}, nil
	case timeoutErrorReturnCode:
		logger.Info("running submission failed: time limit exceeded")
		return RunOutput{
			ReturnCode:        data.StatusCode,
			TimeLimitExceeded: true,
			StdOut:            stdOut,
			StdErr:            stdErr,
		}, nil
	case oomErrorReturnCode:
		logger.Info("running submission failed: memory limit exceeded")
		return RunOutput{
			ReturnCode:          data.StatusCode,
			MemoryLimitExceeded: true,
			StdOut:              stdOut,
			StdErr:              stdErr,
		}, nil
	default:
		logger.With(zap.Int64("status_code", data.StatusCode)).
			Info("running submission failed: program exited with non-zero code")
		return RunOutput{
			ReturnCode: data.StatusCode,
			StdOut:     stdOut,
			StdErr:     stdErr,
		}, nil
	}
}

func (t testCaseRun) onContainerWaitError(ctx context.Context, containerID string, err error) (RunOutput, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	if errors.Is(err, context.DeadlineExceeded) {
		logger.Info("running submission failed: time limit exceeded")
		return RunOutput{TimeLimitExceeded: true}, nil
	}

	logger.
		With(zap.String("container_id", containerID)).
		With(zap.Error(err)).
		Error("failure happened while waiting for container")
	return RunOutput{}, err
}

func (t testCaseRun) Run(
	ctx context.Context,
	programFilePath string,
	input string,
	timeLimitInMillisecond uint64,
	memoryLimitInByte uint64,
) (RunOutput, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	workingDir := t.getWorkingDir()
	programFileDirectory, programFileName := filepath.Split(programFilePath)
	containerProgramFilePath := filepath.Join(workingDir, programFileName)

	containerCreateResponse, err := t.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: t.testCaseRunConfig.Image,
		Cmd: t.getContainerCommand(
			t.testCaseRunConfig.CommandTemplate, containerProgramFilePath, timeLimitInMillisecond),
		WorkingDir:   workingDir,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:%s", programFileDirectory, workingDir)},
		Resources: container.Resources{
			CPUQuota: t.testCaseRunConfig.CPUQuota,
			Memory:   int64(memoryLimitInByte),
		},
		NetworkMode: "none",
	}, nil, nil, "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create run test case container")
		return RunOutput{}, err
	}

	defer func() {
		err = t.dockerClient.ContainerRemove(ctx, containerCreateResponse.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to remove run test case container")
		}
	}()

	containerID := containerCreateResponse.ID
	containerAttachResponse, err := t.dockerClient.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{
		Stream: true, Stdin: true, Stdout: true, Stderr: true,
	})
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).
			Error("failed to attached to stdin of run test case container")
		return RunOutput{}, err
	}

	defer containerAttachResponse.Close()

	_, err = containerAttachResponse.Conn.Write(append([]byte(input), '\n'))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to write to stdin of container")
		return RunOutput{}, err
	}

	err = t.dockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).
			Error("failed to start run test case container")
		return RunOutput{}, err
	}

	containerWaitCtx, containerWaitCancelFunc := context.WithTimeout(ctx, time.Minute)
	defer containerWaitCancelFunc()

	dataChan, errChan := t.dockerClient.ContainerWait(containerWaitCtx, containerID, container.WaitConditionNotRunning)
	select {
	case data := <-dataChan:
		return t.onContainerWaitData(ctx, data, containerAttachResponse)

	case err = <-errChan:
		return t.onContainerWaitError(ctx, containerID, err)

	case <-containerWaitCtx.Done():
		return t.onContainerWaitError(ctx, containerID, containerWaitCtx.Err())
	}
}
