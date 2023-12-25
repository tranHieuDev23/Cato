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
		if timeoutDuration, err := compileConfig.GetTimeoutDuration(); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to get timeout duration")
			return nil, err
		} else {
			c.timeoutDuration = timeoutDuration
		}

		if memoryInBytes, err := compileConfig.GetMemoryInBytes(); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to get memory in bytes")
			return nil, err
		} else {
			c.memoryInBytes = int64(memoryInBytes)
		}

		if _, err := dockerClient.ImagePull(context.Background(), compileConfig.Image, types.ImagePullOptions{}); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to load compile image")
			return nil, err
		}
	}

	return c, nil
}

func (c compile) getSourceFileName(sourceFile *os.File) string {
	if c.compileConfig.SourceFileName != "" {
		return c.compileConfig.SourceFileName
	}

	return filepath.Base(sourceFile.Name())
}

func (c compile) getProgramFileName() string {
	if c.compileConfig.ProgramFileName != "" {
		return c.compileConfig.ProgramFileName
	}

	return uuid.NewString()
}

func (c compile) getCompileCommand(commandTemplate []string, containerSourceFilePath, containerProgramFilePath string) []string {
	command := make([]string, len(commandTemplate))
	for i := range command {
		switch commandTemplate[i] {
		case compileSourceFilePathPlaceHolder:
			command[i] = containerSourceFilePath
		case containerProgramFilePath:
			command[i] = containerProgramFilePath
		default:
			command[i] = commandTemplate[i]
		}
	}

	return command
}

func (c compile) compileSourceFile(ctx context.Context, sourceFile *os.File) (CompileOutput, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	sourceFileName := c.getSourceFileName(sourceFile)
	sourceFilePath := sourceFile.Name()
	containerSourceFilePath := filepath.Join(c.compileConfig.WorkingDir, sourceFileName)

	programFileName := c.getProgramFileName()
	programFileDirectory, err := os.MkdirTemp("", "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create program file directory")
		return CompileOutput{}, err
	}

	programFilePath := filepath.Join(programFileDirectory, programFileName)
	containerProgramFilePath := filepath.Join(c.compileConfig.WorkingDir, programFileName)

	dockerContainerCtx, dockerContainerCancelFunc := context.WithTimeout(ctx, c.timeoutDuration)
	defer dockerContainerCancelFunc()

	containerCreateResponse, err := c.dockerClient.ContainerCreate(dockerContainerCtx, &container.Config{
		Image:        c.compileConfig.Image,
		WorkingDir:   c.compileConfig.WorkingDir,
		Cmd:          c.getCompileCommand(c.compileConfig.CommandTemplate, containerSourceFilePath, containerProgramFilePath),
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", sourceFilePath, containerSourceFilePath),
			fmt.Sprintf("%s:%s", programFilePath, containerProgramFilePath),
		},
		Resources: container.Resources{
			CPUShares: c.compileConfig.CPUShares,
			Memory:    c.memoryInBytes,
		},
		NetworkMode: "none",
	}, nil, nil, "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create compile container")
		return CompileOutput{}, err
	}

	defer func() {
		if err := c.dockerClient.ContainerRemove(ctx, containerCreateResponse.ID, types.ContainerRemoveOptions{}); err != nil {
			logger.With(zap.Error(err)).Error("failed to remove compile container")
		}
	}()

	containerID := containerCreateResponse.ID
	containerAttachResponse, err := c.dockerClient.ContainerAttach(dockerContainerCtx, containerID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failed to attached to compile container")
		return CompileOutput{}, err
	}

	defer containerAttachResponse.Close()

	if err := c.dockerClient.ContainerStart(dockerContainerCtx, containerID, types.ContainerStartOptions{}); err != nil {
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
			stdcopy.StdCopy(stdoutBuffer, stderrBuffer, containerAttachResponse.Reader)
			return CompileOutput{
				ReturnCode: data.StatusCode,
				StdOut:     stdoutBuffer.String(),
				StdErr:     stderrBuffer.String(),
			}, nil
		}

		return CompileOutput{
			ProgramFilePath: programFilePath,
		}, nil

	case err := <-errChan:
		logger.
			With(zap.String("container_id", containerID)).
			With(zap.Error(err)).
			Error("failure happened while waiting for container")
		return CompileOutput{}, err
	}
}

func (c compile) Compile(ctx context.Context, content string) (CompileOutput, error) {
	logger := utils.LoggerWithContext(ctx, c.logger)

	sourceFile, err := os.CreateTemp("", "")
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create source file")
		return CompileOutput{}, err
	}

	if _, err := sourceFile.Write([]byte(content)); err != nil {
		logger.With(zap.Error(err)).Error("failed to write source file")
		return CompileOutput{}, err
	}

	if c.compileConfig == nil {
		// Interpreted language, the source file is the program itself
		return CompileOutput{
			ProgramFilePath: sourceFile.Name(),
		}, nil
	}

	defer func() {
		if err := os.Remove(sourceFile.Name()); err != nil {
			logger.With(zap.Error(err)).Error("failed to remove source file")
		}
	}()

	// Compiled language, need to compile inside a container first
	return c.compileSourceFile(ctx, sourceFile)
}
