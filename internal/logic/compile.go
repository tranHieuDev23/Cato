package logic

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	compileSourceFilePathPlaceHolder  = "$SOURCE"
	compileProgramFilePathPlaceHolder = "$PROGRAM"
)

type CompileOutput struct {
	ProgramFilePath string
	ReturnCode      int64
	StdOut          string
	StdErr          string
}

type Compile interface {
	Compile(ctx context.Context, content string) (CompileOutput, error)
}

type compile struct {
	dockerClient  *client.Client
	logger        *zap.Logger
	language      string
	compileConfig *configs.Compile

	timeoutDuration time.Duration
	memoryInBytes   int64
}

func NewCompile(
	dockerClient *client.Client,
	logger *zap.Logger,
	language string,
	compileConfig *configs.Compile,
) (Compile, error) {
	c := &compile{
		dockerClient: dockerClient,
		logger: logger.
			With(zap.String("language", language)).
			With(zap.Any("compile_config", compileConfig)),
		language:      language,
		compileConfig: compileConfig,
	}

	if compileConfig != nil {
		timeoutDuration, err := compileConfig.GetTimeoutDuration()
		if err != nil {
			c.logger.With(zap.Error(err)).Error("failed to get timeout duration")
			return nil, err
		}

		c.timeoutDuration = timeoutDuration

		memoryInBytes, err := compileConfig.GetMemoryInBytes()
		if err != nil {
			c.logger.With(zap.Error(err)).Error("failed to get memory in bytes")
			return nil, err
		}

		c.memoryInBytes = int64(memoryInBytes)

		c.logger.Info("pulling compile image")
		_, err = dockerClient.ImagePull(context.Background(), compileConfig.Image, types.ImagePullOptions{})
		if err != nil {
			c.logger.With(zap.Error(err)).Error("failed to load compile image")
			return nil, err
		}
	}

	return c, nil
}

func (c compile) getProgramFileName() string {
	if c.compileConfig.ProgramFileName != "" {
		return c.compileConfig.ProgramFileName
	}

	return uuid.NewString()
}

func (c compile) getWorkingDir() string {
	if c.compileConfig.WorkingDir != "" {
		return c.compileConfig.WorkingDir
	}

	return "/work"
}

func (c compile) getCompileCommand(
	commandTemplate []string,
	containerSourceFilePath,
	containerProgramFilePath string,
) []string {
	command := make([]string, len(commandTemplate))
	for i := range command {
		switch commandTemplate[i] {
		case compileSourceFilePathPlaceHolder:
			command[i] = containerSourceFilePath
		case compileProgramFilePathPlaceHolder:
			command[i] = containerProgramFilePath
		default:
			command[i] = commandTemplate[i]
		}
	}

	return command
}

func (c compile) compileSourceFile(
	ctx context.Context,
	hostWorkingDir string,
) (CompileOutput, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	workingDir := c.getWorkingDir()
	containerSourceFilePath := filepath.Join(workingDir, c.compileConfig.SourceFileName)
	programFileName := c.getProgramFileName()
	programFilePath := filepath.Join(hostWorkingDir, programFileName)
	containerProgramFilePath := filepath.Join(workingDir, programFileName)

	dockerContainerCtx, dockerContainerCancelFunc := context.WithTimeout(ctx, c.timeoutDuration)
	defer dockerContainerCancelFunc()

	containerCreateResponse, err := c.dockerClient.ContainerCreate(dockerContainerCtx, &container.Config{
		Image:        c.compileConfig.Image,
		WorkingDir:   workingDir,
		Cmd:          c.getCompileCommand(c.compileConfig.CommandTemplate, containerSourceFilePath, containerProgramFilePath),
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:%s", hostWorkingDir, workingDir)},
		Resources: container.Resources{
			CPUQuota: c.compileConfig.CPUQuota,
			Memory:   c.memoryInBytes,
		},
		NetworkMode: "none",
	}, nil, nil, "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create compile container")
		return CompileOutput{}, err
	}

	defer func() {
		err = c.dockerClient.ContainerRemove(ctx, containerCreateResponse.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to remove compile container")
		}
	}()

	containerID := containerCreateResponse.ID
	containerAttachResponse, err := c.dockerClient.
		ContainerAttach(dockerContainerCtx, containerID, types.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		})
	if err != nil {
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failed to start compile container")
		return CompileOutput{}, err
	}

	defer containerAttachResponse.Close()

	err = c.dockerClient.ContainerStart(dockerContainerCtx, containerID, types.ContainerStartOptions{})
	if err != nil {
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failed to attached to compile container")
		return CompileOutput{}, err
	}

	dataChan, errChan := c.dockerClient.ContainerWait(dockerContainerCtx, containerID, container.WaitConditionNotRunning)
	select {
	case data := <-dataChan:
		if data.StatusCode != 0 {
			logger.
				With(zap.Int64("status_code", data.StatusCode)).
				Info("compiling submission failed: compiler exited with non-zero code")

			stdoutBuffer := new(bytes.Buffer)
			stderrBuffer := new(bytes.Buffer)
			_, err = stdcopy.StdCopy(stdoutBuffer, stderrBuffer, containerAttachResponse.Reader)
			if err != nil {
				logger.With(zap.Error(err)).Error("failed to read from stdout and stderr of container")
				return CompileOutput{}, err
			}

			return CompileOutput{
				ReturnCode: data.StatusCode,
				StdOut:     stdoutBuffer.String(),
				StdErr:     stderrBuffer.String(),
			}, nil
		}

		return CompileOutput{
			ProgramFilePath: programFilePath,
		}, nil

	case err = <-errChan:
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failure happened while waiting for container")
		return CompileOutput{}, err
	}
}

func (c compile) createSourceFile(
	ctx context.Context,
	hostWorkingDir,
	sourceFileName,
	content string,
) (*os.File, error) {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("host_woking_dir", hostWorkingDir)).
		With(zap.String("source_file_name", sourceFileName))

	sourceFilePath := filepath.Join(hostWorkingDir, sourceFileName)
	sourceFile, err := os.Create(sourceFilePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create source file")
		return nil, err
	}

	defer sourceFile.Close()

	_, err = sourceFile.WriteString(content)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to write source file")
		return nil, err
	}

	return sourceFile, nil
}

func (c compile) Compile(ctx context.Context, content string) (CompileOutput, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	hostWorkingDir, err := os.MkdirTemp("", "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create working directory for compiling")
		return CompileOutput{}, err
	}

	var sourceFile *os.File
	if c.compileConfig == nil {
		sourceFile, err = c.createSourceFile(ctx, hostWorkingDir, uuid.NewString(), content)
		if err != nil {
			return CompileOutput{}, err
		}

		return CompileOutput{
			ProgramFilePath: sourceFile.Name(),
		}, nil
	}

	sourceFile, err = c.createSourceFile(ctx, hostWorkingDir, c.compileConfig.SourceFileName, content)
	if err != nil {
		return CompileOutput{}, err
	}

	defer func() {
		if err = os.Remove(sourceFile.Name()); err != nil {
			logger.With(zap.Error(err)).Error("failed to remove source file")
		}
	}()

	// Compiled language, need to compile inside a container first
	return c.compileSourceFile(ctx, hostWorkingDir)
}
